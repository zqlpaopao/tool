package icmp

import (
	"encoding/binary"
	"github.com/zqlpaopao/tool/data-any-pool/pkg"
	"os"
	"sync"
	"syscall"
)

type ErrInfo struct {
	Err  error
	Ping syscall.Sockaddr
	Tag  string
}

type Pool struct {
	err             error
	resChan         chan *ReceiveMMsg
	option          *Option
	prepareChanV4   chan *Ping
	prepareChanV6   chan *Ping
	errChan         chan *ErrInfo
	done            chan struct{}
	errP            pkg.Pool[*ErrInfo]
	PacketP         pkg.Pool[*Packet]
	receiveMMsgPool pkg.Pool[*ReceiveMMsg]
	pingP           pkg.Pool[*Ping]
	seqP            pkg.Pool[[]byte]
	wgSend          sync.WaitGroup
	wgRec           sync.WaitGroup
	wgErr           sync.WaitGroup
	wgCall          sync.WaitGroup
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
	return p
}

// init -- --------------------------
// --> @Describe init the object
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) init() {
	if p.option.v6 {
		p.prepareChanV6 = make(chan *Ping, p.option.prepareV6ChLen)
	}
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
	var fdV4, fdV6 int
	p.wgRec.Add(1)
	if fdV4,
		err = p.makeRecFDV4(); nil != err {
		return
	}
	go p.revIpv4(fdV4)

	if !p.option.v6 {
		return
	}
	p.wgRec.Add(1)
	if fdV6,
		err = p.makeRecFDV6(); nil != err {
		return
	}
	go p.revIpv6(fdV6)
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
		var conn int
		if conn, err = p.makeSendConn4(); nil != err {
			return
		}
		go p.startSenderV4(conn)
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
		var conn int
		if conn, err = p.makeSendConn6(); nil != err {
			return
		}
		go p.startSenderV6(conn)
	}
	return
}

// startSender -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) startSenderV4(connSend int) {
	defer p.option.Recover()
LOOP:
	for {
		select {
		case v, ok := <-p.prepareChanV4:
			if !ok {
				break LOOP
			}
			p.Ping4(connSend, v)
		}
	}
	_ = syscall.Close(connSend)
	p.wgSend.Done()
}

// startSender -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) startSenderV6(fd int) {
	defer p.option.Recover()
LOOP:
	for {
		select {
		case v, ok := <-p.prepareChanV6:
			if !ok {
				break LOOP
			}
			p.Ping6(fd, v)
		}
	}
	_ = syscall.Close(fd)
	p.wgSend.Done()
}

// Ping4 -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) Ping4(connSend int, ps *Ping) {
	if err := p.sendICMPV4(ps, connSend); err != nil {
		errIn := p.errP.Get()
		errIn.Tag, errIn.Ping, errIn.Err = "sendICMPV4", ps.SocketAddrV4, err
		p.errChan <- errIn
		return
	}
}

// Ping6 -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) Ping6(fd int, ps *Ping) {
	if err := p.sendICMPV6(ps, fd); err != nil {
		errIn := p.errP.Get()
		errIn.Tag, errIn.Ping, errIn.Err = "sendICMPV6", ps.SocketAddrV6, err
		p.errChan <- errIn
		return
	}
}

// Submit -- --------------------------
// --> @Describe submit the icmp  object
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) Submit(ps *Ping) (err error) {
	if p.err != nil {
		return p.err
	}
	if ps.Ipv4 {
		p.prepareChanV4 <- ps
		return
	}
	p.prepareChanV6 <- ps
	return
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
func (p *Pool) callBack(pkt *ReceiveMMsg) {
	var packet *Packet
	if pkt.V4 {
		packet = p.ParserIcmpV4(pkt)
	} else {
		packet = p.ParserIcmpV6(pkt)
	}
	p.option.onRevFunc(packet)
	p.receiveMMsgPool.Put(pkt)
	p.PacketP.Put(packet)
}

// ParserIcmpV4 -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) ParserIcmpV4(pkt *ReceiveMMsg) *Packet {
	packet := p.PacketP.Get()
	packet.TXTime,
		packet.RXTime,
		packet.Ttl,
		packet.Seq =
		bytesToTime(pkt.Data[28:36]),
		pkt.RXTime,
		binary.BigEndian.Uint16(pkt.Data[36:38]),
		binary.BigEndian.Uint16(pkt.Data[38:40])

	if v, ok := pkt.Dest.(*syscall.SockaddrInet4); ok {
		packet.Dest = v.Addr[:]
	}
	return packet

}

// ParserIcmpV6 -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) ParserIcmpV6(pkt *ReceiveMMsg) *Packet {
	packet := p.PacketP.Get()
	packet.TXTime,
		packet.RXTime,
		packet.Ttl,
		packet.Seq =
		bytesToTime(pkt.Data[8:16]),
		pkt.RXTime,
		binary.BigEndian.Uint16(pkt.Data[16:18]),
		binary.BigEndian.Uint16(pkt.Data[18:20])

	if v, ok := pkt.Dest.(*syscall.SockaddrInet6); ok {
		packet.Dest = v.Addr[:]
	}
	return packet

}

// sendICMP -- --------------------------
// --> @Describe send the packet
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) sendICMPV4(ps *Ping, fd int) (err error) {

	if p.option.everyTTL {
		if err = setIPv4HopLimit(fd, p.option.ttl); nil != err {
			return
		}

	}
	for {
		p.fillIcmpData(ps)
		if err = syscall.Sendto(fd, ps.Data, 0, ps.SocketAddrV4); err != nil {
			if netErr, ok := err.(syscall.Errno); ok {
				if netErr.Error() == "no buffer space available" {
					continue
				}
				if netErr.Error() == "message too long" {
					continue
				}
			}
			return
		}
		break
	}

	return
}

// fillIcmpData-- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) fillIcmpData(ps *Ping) {
	ps.ReplaceTimeToBytes()

	b := p.seqP.Get()
	defer p.seqP.Put(b)
	if ps.IsOrderSeq {
		binary.BigEndian.PutUint16(b, ps.Sequence)
		ps.Data[19],
			ps.Data[20] =
			b[0],
			b[1]
	}
	// Place checksum back in header; using ^= avoids the
	// assumption the checksum bytes are zero.
	ps.Data[2], ps.Data[3] = 0, 0
	s := checksum(ps.Data)
	ps.Data[2] ^= byte(s)
	ps.Data[3] ^= byte(s >> 8)

}

// sendICMP -- --------------------------
// --> @Describe send the packet
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) sendICMPV6(ps *Ping, fd int) (err error) {
	if p.option.everyTTL {
		if err = setIPv6HopLimit(fd, p.option.ttl); nil != err {
			return
		}
	}
	for {
		p.fillIcmpData(ps)
		//	Port:   33434,
		if err = syscall.Sendto(fd, ps.Data, 0, ps.SocketAddrV6); err != nil {
			if netErr, ok := err.(syscall.Errno); ok {
				if netErr.Error() == "no buffer space available" {
					continue
				}
				if netErr.Error() == "message too long" {
					continue
				}
			}
			return
		}
		break
	}
	ps.ReplaceTimeToBytes()
	return
}

// -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func setIPv4HopLimit(fd int, v int) error {
	err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_TTL, v)
	if err != nil {
		return os.NewSyscallError("setIPv4HopLimit:setSockOpt", err)
	}
	return nil
}

// -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func setIPv6HopLimit(fd int, v int) error {
	err := syscall.SetsockoptInt(fd, syscall.IPPROTO_IPV6, syscall.IPV6_UNICAST_HOPS, v)
	if err != nil {
		return os.NewSyscallError("setIPv6HopLimit:setSockOpt", err)
	}
	return nil
}

// GetObjectPingPool -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) GetObjectPingPool() *Ping {
	return p.pingP.Get()
}

// GetPid -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) GetPid() int {
	return p.option.pid
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
	close(p.prepareChanV4)
	close(p.prepareChanV6)
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
