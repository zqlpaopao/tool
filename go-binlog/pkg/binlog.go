package pkg

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	//"github.com/davecgh/go-spew/spew"
	//"github.com/go-mysql-org/go-mysql/replication"
	"github.com/go-mysql-org/go-mysql/replication"
	"io"
	"net"
	"os"
	"strconv"
	"time"
)

type server struct {
	cfg             *Config
	ctx             context.Context
	conn            net.Conn
	io              *PacketIo
	opt             *option
	registerSuccess bool
	err             error
	errMsg chan error
	tryDump chan struct{}
	tryNum  int32
	//dataMsg
}

func NewServer(cfg *Config, f ...Options) *server {
	o := &option{}
	return &server{
		cfg:             cfg,
		ctx:             nil,
		conn:            nil,
		io:              nil,
		registerSuccess: false,
		err:             nil,
		opt:             o.WithOption(f...),
		errMsg: make(chan error,2),
		tryDump: make(chan struct{}),
	}
}

func (s *server) Run() {
	if s.err = s.checkArgs(); s.err != nil {
		os.Exit(3)
		return
	}
	s.tidyArgs()

	defer s.Quit()
	s.dump()
}

//tidy args
func (s *server) tidyArgs() {
	if s.opt.conTime < ConnTime {
		s.opt.conTime = ConnTime
	}
}

//checkArgs config checkout
func (s *server) checkArgs() error {
	if s.cfg.Host == "" {
		return errors.New("host is empty")
	}
	if s.cfg.Port < MinPort {
		return errors.New("port is less " + strconv.Itoa(MinPort))
	}
	if s.cfg.User == "" {
		return errors.New("user is empty")
	}
	if s.cfg.LogFile == "" {
		return errors.New("logFile is empty")
	}
	if s.cfg.Position < 0 {
		return errors.New("position is less 0")
	}
	if s.cfg.ServerId < 0 {
		return errors.New("serverId is less 0")
	}
	return nil
}

func (s *server) dump() {
	if s.err = s.handshake(); nil != s.err {
		panic(s.err)
	}

	s.invalidChecksum()

	if s.err = s.register(); nil != s.err {
		os.Exit(56)
		return
	}

	if s.err = s.writeDumpCommand(); nil != s.err {
		os.Exit(78)

		return
	}
	parser := replication.NewBinlogParser()
	for {
		select {
			case _,ok := <-s.tryDump:
				if ok{
					goto END
				}
			default:
				data, err := s.io.readPacket()
				if err != nil || len(data) == 0 {
					//write chan error
					//s.errMsg<-err
					continue
				}

				//
				//if data[0] == OK_HEADER {
				//	//skip ok
				//	data = data[1:]
				//	var e *replication.BinlogEvent
				//	if e, s.err = parser.Parse(data); s.err == nil {
				//		fmt.Println(12345,string(data))
				//		e = e
				//		//e.Dump(os.Stdout)
				//		fmt.Println(s.err)
				//	} else {
				//		fmt.Println(123,err)
				//	}
				//} else {
				//	fmt.Println(123)
				//	fmt.Println(string(data))
				//	s.io.HandleError(data)
				//}
				//var data1  = data
				//spew.Dump("data1",data1)
				//spew.Dump("data",data)
				s.HandleData(data,parser)
		}
		time.Sleep(1*time.Second)
		//time.Sleep(2 * time.Second)
		//s.query("select 1")



		//if data[0] == OK_HEADER {
		//	//skip ok
		//	data = data[1:]
		//	var e *replication.BinlogEvent
		//	if e, s.err = parser.Parse(data); s.err == nil {
		//		fmt.Println(12345,string(data))
		//		e = e
		//		//e.Dump(os.Stdout)
		//		fmt.Println(s.err)
		//	} else {
		//		fmt.Println(123,err)
		//	}
		//} else {
		//	fmt.Println(123)
		//	fmt.Println(string(data))
		//	s.io.HandleError(data)
		//}
	}

	END:
		fmt.Println("try------")
	s.dump()
}

//invalidChecksum
//When binlog event is sent back, 4 additional bytes will be added
//for verification when the event content is finally obtained
func (s *server) invalidChecksum() {
	sql := `SET @master_binlog_checksum='NONE'`
	if err := s.query(sql); err != nil {
		fmt.Println(err)
	}
	//readPacket must read from tcp connection , either will be blocked
	_, _ = s.io.readPacket()
}

//handshake Establish connection
func (s *server) handshake() (err error) {
	var (
		conn     net.Conn
		data, pk []byte
	)
	if conn, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port), time.Duration(s.opt.conTime)*time.Second); err != nil {
		return err
	}

	tc := conn.(*net.TCPConn)
	if err = tc.SetKeepAlive(true); nil != err {
		return
	}
	if err = tc.SetNoDelay(true); nil != err {
		return
	}

	s.conn = tc
	s.io = &PacketIo{}
	s.io.r = bufio.NewReaderSize(s.conn, 16*1024)
	s.io.w = tc

	if data, err = s.io.readPacket(); err != nil {
		return
	}

	if data[0] == ERR_HEADER {
		return errors.New("error packet")
	}
	if data[0] < MinProtocolVersion {
		return fmt.Errorf("version is too lower, current:%d", data[0])
	}

	pos := 1 + bytes.IndexByte(data[1:], 0x00) + 1
	connId := binary.LittleEndian.Uint32(data[pos : pos+4])
	pos += 4
	salt := data[pos : pos+8]

	pos += 8 + 1
	capability := uint32(binary.LittleEndian.Uint16(data[pos : pos+2]))

	pos += 2

	var status uint16
	var pluginName string
	if len(data) > pos {
		//skip charset
		pos++
		status = binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2
		capability = uint32(binary.LittleEndian.Uint16(data[pos:pos+2]))<<16 | capability
		pos += 2

		pos += 10 + 1
		salt = append(salt, data[pos:pos+12]...)
		pos += 13

		if end := bytes.IndexByte(data[pos:], 0x00); end != -1 {
			pluginName = string(data[pos : pos+end])
		} else {
			pluginName = string(data[pos:])
		}
	}

	fmt.Printf("conn_id:%v, status:%d, plugin:%v\n", connId, status, pluginName)

	//write
	capability = 500357
	length := 4 + 4 + 1 + 23
	length += len(s.cfg.User) + 1

	pass := []byte(s.cfg.Pass)
	auth := calPassword(salt[:20], pass)
	length += 1 + len(auth)
	data = make([]byte, length+4)

	data[4] = byte(capability)
	data[5] = byte(capability >> 8)
	data[6] = byte(capability >> 16)
	data[7] = byte(capability >> 24)

	//utf8
	data[12] = byte(33)
	pos = 13 + 23
	if len(s.cfg.User) > 0 {
		pos += copy(data[pos:], s.cfg.User)
	}

	pos++
	data[pos] = byte(len(auth))
	pos += 1 + copy(data[pos+1:], auth)

	if err = s.io.writePacket(data); err != nil {
		return fmt.Errorf("write auth packet error")
	}

	if pk, err = s.io.readPacket(); err != nil {
		return
	}

	if pk[0] == OK_HEADER {
		return nil
	} else if pk[0] == ERR_HEADER {
		s.io.HandleError(pk)
		return errors.New("handshake error ")
	}

	return nil
}

func (s *server) writeDumpCommand() (err error) {
	s.io.seq = 0
	data := make([]byte, 4+1+4+2+4+len(s.cfg.LogFile))
	pos := 4
	data[pos] = 18 //dump binlog
	pos++
	binary.LittleEndian.PutUint32(data[pos:], uint32(s.cfg.Position))
	pos += 4

	//dump command flag
	binary.LittleEndian.PutUint16(data[pos:], 0)
	pos += 2

	binary.LittleEndian.PutUint32(data[pos:], uint32(s.cfg.ServerId))
	pos += 4

	copy(data[pos:], s.cfg.LogFile)

	if err = s.io.writePacket(data); nil != err {
		return err
	}
	//ok
	res, _ := s.io.readPacket()
	if res[0] == OK_HEADER {
		return
	} else {
		s.io.HandleError(res)
	}
	return nil
}

func (s *server) register() (err error) {
	s.io.seq = 0
	hostname, _ := os.Hostname()
	data := make([]byte, 4+1+4+1+len(hostname)+1+len(s.cfg.User)+1+len(s.cfg.Pass)+2+4+4)
	pos := 4
	data[pos] = 21 //register slave  command
	pos++
	binary.LittleEndian.PutUint32(data[pos:], uint32(s.cfg.ServerId))
	pos += 4

	data[pos] = uint8(len(hostname))
	pos++
	n := copy(data[pos:], hostname)
	pos += n

	data[pos] = uint8(len(s.cfg.User))
	pos++
	n = copy(data[pos:], s.cfg.User)
	pos += n

	data[pos] = uint8(len(s.cfg.Pass))
	pos++
	n = copy(data[pos:], s.cfg.Pass)
	pos += n

	binary.LittleEndian.PutUint16(data[pos:], uint16(s.cfg.Port))
	pos += 2

	binary.LittleEndian.PutUint32(data[pos:], 0)
	pos += 4

	//master id = 0
	binary.LittleEndian.PutUint32(data[pos:], 0)

	if err = s.io.writePacket(data); nil != err {
		return
	}
	var res []byte
	//ok
	if res, err = s.io.readPacket(); nil != err {

	}
	if res[0] == OK_HEADER {
		s.registerSuccess = true
	} else {
		s.io.HandleError(data)
	}
	return nil
}

func (s *server) writeCommand(command byte) {
	s.io.seq = 0
	_ = s.io.writePacket([]byte{
		0x01, //1 byte long
		0x00,
		0x00,
		0x00, //seq
		command,
	})
}

func (s *server) query(q string) error {
	s.io.seq = 0
	length := len(q) + 1
	data := make([]byte, length+4)
	data[4] = 3
	copy(data[5:], q)
	return s.io.writePacket(data)
}

func (s *server) Quit() {
	//quit
	s.writeCommand(byte(1))
	//maybe only close
	if err := s.conn.Close(); nil != err {
		fmt.Printf("error in close :%v\n", err)
	}
}

type PacketIo struct {
	r   *bufio.Reader
	w   io.Writer
	seq uint8
}

func (p *PacketIo) readPacket() ([]byte, error) {
	//to read header
	header := []byte{0, 0, 0, 0}
	if _, err := io.ReadFull(p.r, header); err != nil {
		return nil, err
	}

	length := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)
	if length == 0 {
		p.seq++
		return []byte{}, nil
	}

	if length == 1 {
		return nil, fmt.Errorf("invalid payload")
	}

	seq := uint8(header[3])
	if p.seq != seq {
		return nil, fmt.Errorf("invalid seq %d", seq)
	}

	p.seq++
	data := make([]byte, length)
	if _, err := io.ReadFull(p.r, data); err != nil {
		return nil, err
	} else {
		if length < MaxPayloadLength {
			return data, nil
		}
		var buf []byte
		buf, err = p.readPacket()
		if err != nil {
			return nil, err
		}
		if len(buf) == 0 {
			return data, nil
		} else {
			return append(data, buf...), nil
		}
	}
}

func (p *PacketIo) writePacket(data []byte) error {
	length := len(data) - 4
	if length >= MaxPayloadLength {
		data[0] = 0xff
		data[1] = 0xff
		data[2] = 0xff
		data[3] = p.seq

		if n, err := p.w.Write(data[:4+MaxPayloadLength]); err != nil {
			return fmt.Errorf("write find error")
		} else if n != 4+MaxPayloadLength {
			return fmt.Errorf("not equal max pay load length")
		} else {
			p.seq++
			length -= MaxPayloadLength
			data = data[MaxPayloadLength:]
		}
	}

	data[0] = byte(length)
	data[1] = byte(length >> 8)
	data[2] = byte(length >> 16)
	data[3] = p.seq

	if n, err := p.w.Write(data); err != nil {
		return errors.New("write find error")
	} else if n != len(data) {
		return errors.New("not equal length")
	} else {
		p.seq++
		return nil
	}
}

func calPassword(scramble, password []byte) []byte {
	crypt := sha1.New()
	crypt.Write(password)
	stage1 := crypt.Sum(nil)

	crypt.Reset()
	crypt.Write(stage1)
	hash := crypt.Sum(nil)

	crypt.Reset()
	crypt.Write(scramble)
	crypt.Write(hash)
	scramble = crypt.Sum(nil)

	for i := range scramble {
		scramble[i] ^= stage1[i]
	}

	return scramble
}

func (p *PacketIo) HandleError(data []byte) {
	pos := 1
	code := binary.LittleEndian.Uint16(data[pos:])
	pos += 2
	pos++
	state := string(data[pos : pos+5])
	pos += 5
	msg := string(data[pos:])
	fmt.Printf("code:%d, state:%s, msg:%s\n", code, state, msg)
}

func (s *server) HandleData(data []byte,parser *replication.BinlogParser) {
	var (
		errData *dataInfo
		length int
		err error
	)
	if data[0] == OK_HEADER {
		data = data[1:]
		var e *replication.BinlogEvent
		if e, err = parser.Parse(data); err != nil {
			if s.opt.tryDump && atomic.LoadInt32(&s.tryNum) < s.opt.tryNum && strings.Contains(err.Error(),"UnknownEvent") && strings.Contains(err.Error(),"less event length"){
				//123 invalid data size 101 in event UnknownEvent, less event length 30445
				if err == nil{
					return
				}
				if length,err = s.tryDum(err.Error());nil != err{
					s.errMsg<-err
				}
				atomic.AddInt32(&s.tryNum,1)
				s.startDump(length)
			}
			return
		}
		//var ds []byte
		//err = e.Event.Decode(ds)

		fmt.Println(e.Header.EventType)
		//fmt.Println(e.Event.Decode(data).Error())
		fmt.Println(replication.WRITE_ROWS_EVENTv2)

		fmt.Println(err)
		//fmt.Println(ds)
		//os.Exit(3)
		e.Dump(os.Stdout)
		fmt.Println(123456)
		//fmt.Println(err)
		//e.Dump(os.Stdout)
		e = e


	} else {
		pos := 1
		code := binary.LittleEndian.Uint16(data[pos:])
		pos += 2
		pos++
		state := string(data[pos : pos+5])
		pos += 5
		msg := string(data[pos:])

		if s.opt.print{
			fmt.Printf("code:%d, state:%s, msg:%s\n", code, state, msg)
		}
		errData = &dataInfo{
			code:    FailCode,
			errInfo: &errInfo{
				code:  code,
				state: state,
				msg:   msg,
			},
		}
		errData = errData
		if s.opt.tryDump && atomic.LoadInt32(&s.tryNum) < s.opt.tryNum && strings.Contains(msg,ReadFrom)  {
			if length ,err  = s.lastByte(msg);nil != err{
				fmt.Println(err)
				os.Exit(33)
				//s.errMsg<-err
				return
			}

			s.startDump(length)
		}
		atomic.AddInt32(&s.tryNum,1)
		//spew.Dump(errData)
	}



	//fmt.Printf("code:%d, state:%s, msg:%s\n", code, state, msg)
}

func(s *server)startDump(length int){
	s.cfg.Position = length
	fmt.Println(s.cfg.Position)
	//os.Exit(3)
	go func() {
		s.tryDump<- struct{}{}
	}()
}

func (s *server)tryDum(str string)(length int,err error){
	if str == ""{
		return 0,errors.New("str is empty")
	}
	var (
		info []string
	)
	info = strings.Split(str," ")

	if len(info)<11{
		return
	}
	if length,err = strconv.Atoi(info[3]);nil != err{
		return
	}

	if length< s.cfg.Position{
		return 0,errors.New("Position is error")
	}
	return
}

func (s *server)lastByte(str string)(length int,err error){
	if str == ""{
		return 0,errors.New("str is empty")
	}
	var (
		info []string
	)
	info = strings.Split(str,",")
	if len(info)<2{
		return
	}
	info = info[1:2]
	if len(info)<1{
		return
	}
	info = strings.Split(info[0]," ")
	if length,err = strconv.Atoi(info[len(info)-1]);nil != err{
		return
	}

	if length< s.cfg.Position{
		return 0,errors.New("position is error")
	}
	return
}


func (s *server)Error()error{
	return s.err
}