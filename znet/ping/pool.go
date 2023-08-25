package ping

import (
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
	Err  error
	Ping string
	Tag  string
}

type Pool struct {
	pingMap      cmap.ConcurrentMap[string, *Ping]
	err          error
	resChan      chan *Packet
	option       *Option
	prepareChan  chan *Ping
	errChan      chan *ErrInfo
	done         chan struct{}
	errP         pkg.Pool[*ErrInfo]
	PacketP      pkg.Pool[*Packet]
	packetP      pkg.Pool[*packet]
	pingP        pkg.Pool[*Ping]
	body         pkg.Pool[*icmp.Echo]
	msg          pkg.Pool[*icmp.Message]
	msgBytes     pkg.Pool[unsafe.Pointer]
	timeToBytesP pkg.Pool[[]byte]
	wgSend       sync.WaitGroup
	wgRec        sync.WaitGroup
	wgErr        sync.WaitGroup
	wgCall       sync.WaitGroup
}

// Run -- --------------------------
// --> @Describe start the Pool
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) Run() *Pool {
	defer p.option.Recover()
	p.init()

	p.loopRev()

	p.loopErr()

	p.initClient()

	go p.loopTidyMap()

	return p
}

// init -- --------------------------
// --> @Describe init the object
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) init() {
	cmap.SHARD_COUNT = p.option.pingMapLen
	p.pingMap = cmap.New[*Ping]()
}

// loopErr -- --------------------------
// --> @Describe goroutine loop the errCh
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) loopErr() {
	if p.err != nil {
		return
	}
	p.wgErr.Add(1)
	for i := 0; i < 1; i++ {
		go p.callBackErrFunc()
	}

}

// callBackErrFunc -- ------------------------------
// --> @Describe call back the func to send err info
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) callBackErrFunc() {
	defer p.option.Recover()
	for {
		select {
		case v, ok := <-p.errChan:
			if !ok {
				goto END
			}
			p.option.errCallBack(v)
			p.errP.Put(v)
		}
	}
END:
	p.wgErr.Done()
}

// initClient -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) initClient() {
	p.err = p.makeRecClient()
	p.err = p.makeSendClient()
}

// makeRecClient -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeRecClient() (err error) {
	var (
		connRec4 net.PacketConn
		connRec6 net.PacketConn
	)
	p.wgRec.Add(1)
	if connRec4, err = p.makeRecConn4(p.option.pid); nil != err {
		return
	}
	go p.revIpv4(connRec4, p.option.pid)

	if !p.option.v6 {
		return
	}
	p.wgRec.Add(1)
	if connRec6, err = p.makeRecConn6(p.option.pid); nil != err {
		return
	}
	go p.revIpv6(connRec6, p.option.pid)
	return
}

// makeSendClient -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeSendClient() (err error) {
	if err = p.SendClientV4(); nil != err {
		return
	}
	return p.SendClientV6()
}

// SendClientV4 -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) SendClientV4() (err error) {
	p.wgSend.Add(p.option.sendWorker)
	for i := 0; i < p.option.sendWorker; i++ {
		var conn packetConn
		if conn, err = p.makeSendConn4(); nil != err {
			return
		}
		go p.startSender(conn)
	}
	return
}

// SendClientV6 -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) SendClientV6() (err error) {
	if !p.option.v6 {
		return
	}
	p.wgSend.Add(p.option.sendWorker)
	for i := 0; i < p.option.sendWorker; i++ {
		var conn packetConn
		if conn, err = p.makeSendConn6(); nil != err {
			return
		}
		go p.startSender(conn)
	}
	return
}

// startSender -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) startSender(connSend packetConn) {
	defer p.option.Recover()
LOOP:
	for {
		select {
		case v, ok := <-p.prepareChan:
			if !ok {
				break LOOP
			}
			v.Pid = p.option.pid
			p.Ping(connSend, v)
		}
	}
	_ = connSend.Close()
	p.wgSend.Done()
}

// Ping -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) Ping(connSend packetConn, ps *Ping) {
	if err := p.sendICMP(ps, connSend); err != nil {
		errIn := p.errP.Get()
		errIn.Tag,
			errIn.Ping,
			errIn.Err =
			"Ping",
			ps.Addr,
			errors.New("sendICMP->"+err.Error())
		p.errChan <- errIn

		return
	}
}

// loopTidyMap-- --------------------------
// --> @Describe tidy the current map
// --> @params
// --> @return
// -- ------------------------------------
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

// tidyMap -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) tidyMap() {
	its := p.pingMap.Items()
	newMap := cmap.New[*Ping]()

	for k, v := range its {
		newMap.Set(k, v)
	}

	old := unsafe.Pointer(&(p.pingMap))
	newMapPointer := unsafe.Pointer(&newMap)
	atomic.SwapPointer(&old, newMapPointer)
}

// Submit -- --------------------------
// --> @Describe submit the ping  object
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) Submit(ps *Ping) (err error) {
	if p.err != nil {
		return p.err
	}
	if ps.Size < timeSliceLength+trackerLength {
		return fmt.Errorf("size %d is less than minimum required size %d", ps.Size, timeSliceLength+trackerLength)
	}
	//ps.network,
	//	ps.protocol,
	//	ps.Size =
	//	p.option.network,
	//	p.option.protocol,
	//	p.option.dataSize

	p.prepareChan <- ps
	return
}

// set -- --------------------------
// --> @Describe set the ping object
// --> to the current map,and the call
// --> back func check result
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) set(uuid string, ps *Ping) {
	p.pingMap.Set(uuid, ps)
}

// get -- --------------------------
// --> @Describe set the ping object
// --> to the current map,and the call
// --> back func check result
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) get(uuid string) (ping *Ping, ok bool) {
	return p.pingMap.Get(uuid)
}

// loopRev -- --------------------------
// --> @Describe loop the resp chan
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) loopRev() {
	p.wgCall.Add(p.option.callbackWorker)
	for i := 0; i < p.option.callbackWorker; i++ {
		go p.revCallBack()
	}
}

// revCallBack -- --------------------------
// --> @Describe call back
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) revCallBack() {
	defer p.option.Recover()
	for {
		select {
		case v, ok := <-p.resChan:
			if !ok {
				goto END
			}
			p.callBack(v)
		}
	}
END:
	p.wgCall.Done()
}

// callback  -- --------------------------
// --> @Describe call back of get the map
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) callBack(pkt *Packet) {
	var (
		ping *Ping
		ok   bool
	)
	if ping, ok = p.get(pkt.uuid); ping == nil || !ok {
		return
	}
	p.option.onRevFunc(ping, pkt)

	//p.pingMap.Remove(pkt.uuid)

	p.PacketP.Put(pkt)

	//end  race
	//p.pingP.Put(ping)

}

// sendICMP -- --------------------------
// --> @Describe send the packet
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) sendICMP(ps *Ping, conn packetConn) (err error) {
	var (
		dst            net.Addr = ps.Ipaddr
		t, uuidEncoded []byte
	)

	if ps.Mark != 0 {
		if err = conn.SetMark(ps.Mark); err != nil {
			return fmt.Errorf("error setting mark: %v", err)
		}
	}

	if ps.Df {
		if err = conn.SetDoNotFragment(); err != nil {
			return fmt.Errorf("error setting do-not-fragment: %v", err)
		}
	}

	if ps.Protocol == Udp {
		dst = &net.UDPAddr{IP: ps.Ipaddr.IP, Zone: ps.Ipaddr.Zone}
	}

	uuidEncoded = src.String2Bytes(ps.Uuid)[:]

	t = p.timeToBytes(time.Now())
	for i := timeSliceLength; i < timeSliceLength+trackerLength; i++ {
		t[i] = uuidEncoded[i-8]
	}
	if remainSize := ps.Size - timeSliceLength - trackerLength; remainSize > 0 {
		for i := timeSliceLength + trackerLength; i < ps.Size; i++ {
			t[i] = 1
		}
		//t = append(t, bytes.Repeat([]byte{1}, remainSize)...)
	}
	defer p.timeToBytesP.Put(t)

	body := p.body.Get()
	defer p.body.Put(body)
	body.ID,
		body.Seq,
		body.Data =
		ps.Pid,
		ps.Sequence,
		t

	msg := p.msg.Get()
	defer p.msg.Put(msg)
	msg.Type,
		msg.Code,
		msg.Body =
		conn.ICMPRequestType(),
		0,
		body

	p.set(ps.Uuid, ps)

	msgBytes := *(*[]byte)(p.msgBytes.Get())
	defer p.msgBytes.Put(unsafe.Pointer(&msgBytes))
	if msgBytes, err = msg.Marshal(nil); err != nil {
		return
	}
	if p.option.everyTTL {
		conn.SetTTL(p.option.ttl)
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

// GetObjectPingPool -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) GetObjectPingPool() *Ping {
	return p.pingP.Get()
}

// RemoveCurrentMap -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) RemoveCurrentMap(uuid string) {
	p.pingMap.Remove(uuid)
}

// PutObjectPingPool -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) PutObjectPingPool(ping *Ping) {
	p.pingP.Put(ping)
}

// Error -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) Error() error {
	return p.err
}

// Close -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) Close() *Pool {
	close(p.prepareChan)
	p.wgSend.Wait()
	return p
}

// CloseRev -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) CloseRev() {
	close(p.done)
	p.wgRec.Wait()
	close(p.resChan)
	p.wgCall.Wait()
	close(p.errChan)
	p.wgErr.Wait()
}
