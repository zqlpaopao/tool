package ping

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/orcaman/concurrent-map/v2"
	"github.com/zqlpaopao/tool/data-any-pool/pkg"
	"github.com/zqlpaopao/tool/string-byte/src"
	"golang.org/x/net/icmp"
	"net"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

type ErrInfo struct {
	ping *Ping
	err  error
}

type Pool struct {
	option      *Option
	conn4       chan packetConn
	conn4Rev    packetConn
	conn6Rev    packetConn
	conn6       chan packetConn
	pingMap     *cmap.ConcurrentMap[string, *Ping]
	err         error
	errChan     chan *ErrInfo
	done        chan struct{}
	onRevChan   chan *Packet
	wg          *sync.WaitGroup
	errWg       *sync.WaitGroup
	readyChan   chan *Ping
	readWg      *sync.WaitGroup
	revWg       *sync.WaitGroup
	errInfoPool chan *ErrInfo
	recPacket   chan *packet
	packet      chan *Packet
	ping        chan *Ping
	tData       *pkg.Pool[[]byte]
	body        *pkg.Pool[icmp.Echo]
	msg         *pkg.Pool[icmp.Message]
	msgBytes    *pkg.Pool[[]byte]
	uuidEncoded *pkg.Pool[[]byte]
}

func (p *Pool) Run() *Pool {
	defer p.option.Recover()
	p.init()
	p.listen()
	p.loopErr()
	p.loopRev()
	p.rev()
	p.startWorker()
	go p.loopTidyMap()
	return p
}

func (p *Pool) init() {
	cmap.SHARD_COUNT = p.option.currentMapSize
	mp := cmap.New[*Ping]()
	p.pingMap = &mp
	InitUUIDPool(p.option.uuidPool.arraySize, p.option.uuidPool.uuidSize)
}

func (p *Pool) listen() {
	if p.conn4Rev, p.err = p.makeConn4(); p.err != nil {
		return
	}

	if p.conn6Rev, p.err = p.makeConn6(); p.err != nil {
		return
	}

	for i := 0; i < p.option.conn4Size; i++ {
		if conn4, err := p.makeConn4(); err == nil {
			p.conn4 <- conn4
		}
	}

	for i := 0; i < p.option.conn6Size; i++ {
		if conn6, err := p.makeConn6(); err == nil {
			p.conn6 <- conn6
		}
	}
}

func (p *Pool) loopTidyMap() {
	defer p.option.Recover()
	var t = time.NewTicker(p.option.mapTidyInterval)
	for {
		select {
		case <-t.C:
			p.tidyMap()
			t.Reset(p.option.mapTidyInterval)
		case <-p.done:
			goto END
		}

	}
END:
	t.Stop()
}

func (p *Pool) tidyMap() {
	its := p.pingMap.Items()
	newMap := cmap.New[*Ping]()

	for k, v := range its {
		newMap.Set(k, v)
	}

	old := unsafe.Pointer(&(*p.pingMap))
	newMapPointer := unsafe.Pointer(&newMap)
	atomic.SwapPointer(&old, newMapPointer)
}

func (p *Pool) startWorker() {
	if p.err != nil {
		return
	}
	p.wg.Add(p.option.customerSize)
	for i := 0; i < p.option.customerSize; i++ {
		go p.loop()
	}
}

func (p *Pool) loop() {
	defer p.option.Recover()
LOOP:
	for {
		select {
		case v, ok := <-p.readyChan:
			if !ok {
				break LOOP
			}
			p.Ping(v)
		}
	}
	p.wg.Done()
}

func (p *Pool) Ping(ps *Ping) {
	var (
		conn packetConn
		err  error
	)
	if conn, err = p.getConn(ps.ipv4); conn == nil || err != nil {
		errIn := p.getErrRes()
		errIn.ping, errIn.err = ps, errors.New("getConn--->"+err.Error())
		p.errChan <- errIn
		return
	}
	defer p.setConn(ps.ipv4, conn)
	if err = p.sendICMP(ps, conn); err != nil {
		errIn := p.getErrRes()
		errIn.ping, errIn.err = ps, errors.New("sendICMP->"+err.Error())
		p.errChan <- errIn
		return
	}
}

func (p *Pool) getConn(ipv4 bool) (packetConn, error) {
	if ipv4 {
		return p.getConn4()
	}
	return p.getConn6()
}

func (p *Pool) setConn(ipv4 bool, conn packetConn) {
	if ipv4 {
		p.setConn4(conn)
		return
	}
	p.setConn6(conn)

}

func (p *Pool) Submit(ps *Ping) (err error) {
	if p.err != nil {
		return p.err
	}
	if ps.Size < timeSliceLength+trackerLength {
		return fmt.Errorf("size %d is less than minimum required size %d", ps.Size, timeSliceLength+trackerLength)
	}
	ps.id,
		ps.network,
		ps.protocol, ps.Size =
		p.option.pid,
		p.option.network,
		p.option.protocol,
		p.option.dataSize

	p.readyChan <- ps
	return
}

func (p *Pool) set(uuid string, ps *Ping) {
	p.pingMap.Set(uuid, ps)
}

func (p *Pool) get(uuid string) (ping *Ping, ok bool) {
	return p.pingMap.Get(uuid)
}

func (p *Pool) loopErr() {
	p.errWg.Add(1)
	for i := 0; i < 1; i++ {
		go p.callBackErrFunc()
	}

}
func (p *Pool) callBackErrFunc() {
	defer p.option.Recover()
	for {
		select {
		case v, ok := <-p.errChan:
			if !ok {
				goto END
			}
			p.option.errCallBack(v.ping, v.err)
			p.SetErrRes(v)
		}
	}
END:
	p.errWg.Done()
}

func (p *Pool) loopRev() {
	p.revWg.Add(p.option.onRevWorkerNum)
	for i := 0; i < p.option.onRevWorkerNum; i++ {
		go p.revCallBack()
	}
}

func (p *Pool) revCallBack() {
	defer p.option.Recover()
	for {
		select {
		case v, ok := <-p.onRevChan:
			if !ok {
				goto END
			}
			p.callBack(v)
		}
	}
END:
	p.revWg.Done()
}

func (p *Pool) callBack(pkt *Packet) {
	var (
		ping *Ping
		ok   bool
	)
	if ping, ok = p.get(pkt.uuid); ping == nil || !ok {
		return
	}
	p.option.OnRevFunc(ping, pkt)

	SetUUID(*ping.Uuid())
	arr := (*[36]byte)(src.String2Bytes(pkt.Uuid())[0:36])
	SetByte(*arr)

	p.pingMap.Remove(pkt.uuid)
	p.SetPacketRes(pkt)

	//end  race
	p.SetPing(ping)

}

func (p *Pool) rev() {
	if p.err != nil {
		return
	}
	p.readWg.Add(2)
	go p.revIpv4()
	go p.revIpv6()
}

func (p *Pool) Error() error {
	return p.err
}

func (p *Pool) Close() *Pool {
	close(p.readyChan)
	p.wg.Wait()
	return p
}

func (p *Pool) CloseRev() {
	close(p.done)
	p.readWg.Wait()
	close(p.onRevChan)
	p.revWg.Wait()
	close(p.errChan)
	p.errWg.Wait()
	close(p.conn6)
	close(p.conn4)
	p.CloseConn()
}

func (p *Pool) CloseConn() {
	for conn := range p.conn4 {
		_ = conn.Close()
	}
	for conn := range p.conn6 {
		_ = conn.Close()
	}

}

// sendICMP send the packet
func (p *Pool) sendICMP(ps *Ping, conn packetConn) (err error) {
	var (
		dst net.Addr = ps.ipaddr
		t            = *p.tData.Get()
	)

	defer p.tData.Put(&t)

	if ps.mark != 0 {
		if err = conn.SetMark(ps.mark); err != nil {
			return fmt.Errorf("error setting mark: %v", err)
		}
	}

	if ps.df {
		if err = conn.SetDoNotFragment(); err != nil {
			return fmt.Errorf("error setting do-not-fragment: %v", err)
		}
	}

	conn.SetTTL(ps.TTL)

	if ps.protocol == Udp {
		dst = &net.UDPAddr{IP: ps.ipaddr.IP, Zone: ps.ipaddr.Zone}
	}

	uuidEncoded := *p.uuidEncoded.Get()
	defer p.uuidEncoded.Put(&uuidEncoded)
	if uuidEncoded, err = ps.uuid.MarshalBinary(); err != nil {
		return fmt.Errorf("unable to marshal UUID binary: %w", err)
	}

	t = append(timeToBytes(time.Now()), uuidEncoded...)
	if remainSize := ps.Size - timeSliceLength - trackerLength; remainSize > 0 {
		t = append(t, bytes.Repeat([]byte{1}, remainSize)...)
	}

	body := p.body.Get()
	defer p.body.Put(body)
	body.ID,
		body.Seq,
		body.Data =
		ps.id,
		ps.sequence,
		t

	msg := p.msg.Get()
	defer p.msg.Put(msg)
	msg.Type,
		msg.Code,
		msg.Body =
		conn.ICMPRequestType(),
		0,
		body

	p.set(ps.uuid.String(), ps)

	msgBytes := *p.msgBytes.Get()
	defer p.msgBytes.Put(&msgBytes)

	if msgBytes, err = msg.Marshal(nil); err != nil {
		return
	}

	for {
		if _, err = conn.WriteTo(msgBytes, dst); err != nil {
			if netErr, ok := err.(*net.OpError); ok {
				if netErr.Err == syscall.ENOBUFS {
					continue
				}
			}
			return
		}
		break
	}
	return
}
