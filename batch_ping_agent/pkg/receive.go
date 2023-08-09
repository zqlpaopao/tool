package ping

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"runtime"
	"time"
)

// revIpv4 rev the ipv4 packet
func (p *Pool) revIpv4() {
	defer p.option.Recover()
	// Start by waiting for 50 µs and increase to a possible maximum of ~ 100 ms.
	expBackoffTime := newExpBackoff(50*time.Microsecond, 11)
	delay := expBackoffTime.Get()
	var (
		bytes  []byte
		n, ttl int
		err    error
		addr   net.Addr
	)
	if p.option.protocol != "icmp" && runtime.GOOS == "darwin" {
		bytes = make([]byte, p.getMessageLength()+NoIcmpAndDarwin)
	} else {
		bytes = make([]byte, p.getMessageLength())
	}
	for {
		select {
		case <-p.done:
			goto END
		default:
			if err = p.conn4Rev.SetReadDeadline(time.Now().Add(delay)); err != nil {
				errIn := p.getErrRes()
				errIn.ping, errIn.err = nil, errors.New("conn4Rev.SetReadDeadline--->"+err.Error())
				p.errChan <- errIn
				continue
			}
			if n, ttl, addr, err = p.conn4Rev.ReadFrom(bytes); err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Timeout() {
						// Read timeout
						delay = expBackoffTime.Get()
						continue
					}
				}
				errIn := p.getErrRes()
				errIn.ping, errIn.err = nil, errors.New("conn4Rev.ReadFrom--->"+err.Error())
				p.errChan <- errIn
				continue
			}
			pktR := p.getRevPacketRes()
			pktR.receivedAt,
				pktR.addr,
				pktR.bytes,
				pktR.nBytes,
				pktR.ttl,
				pktR.ipv4 =
				time.Now(),
				addr,
				bytes,
				n,
				ttl,
				true
			if err = p.processPacket(pktR); nil != err {
				p.SetRevPacketRes(pktR)
				errIn := p.getErrRes()
				errIn.ping, errIn.err = nil, errors.New("conn4Rev.processPacket--->"+err.Error())
				p.errChan <- errIn
				continue
			}
			p.SetRevPacketRes(pktR)
		}
	}
END:
	p.readWg.Done()
}

// revIpv6 rev the ipv6 packet
func (p *Pool) revIpv6() {
	defer p.option.Recover()
	// Start by waiting for 50 µs and increase to a possible maximum of ~ 100 ms.
	expBackoffTime := newExpBackoff(50*time.Microsecond, 11)
	delay := expBackoffTime.Get()

	var (
		offset int
		bytes  []byte
		n, ttl int
		err    error
		addr   net.Addr
	)
	if p.option.protocol != "icmp" && runtime.GOOS == "darwin" {
		offset, bytes = NoIcmpAndDarwin, make([]byte, p.getMessageLength()+offset)
	} else {
		bytes = make([]byte, p.getMessageLength())
	}

	for {
		select {
		case <-p.done:
			goto END
		default:
			if err = p.conn6Rev.SetReadDeadline(time.Now().Add(delay)); err != nil {
				errIn := p.getErrRes()
				errIn.ping, errIn.err = nil, errors.New("conn6Rev.processPacket--->"+err.Error())
				p.errChan <- errIn
				continue
			}
			if n, ttl, addr, err = p.conn6Rev.ReadFrom(bytes); err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Timeout() {
						// Read timeout
						delay = expBackoffTime.Get()
						continue
					}
				}
				errIn := p.getErrRes()
				errIn.ping, errIn.err = nil, errors.New("conn6Rev.ReadFrom--->"+err.Error())
				p.errChan <- errIn
				continue
			}
			pktR := p.getRevPacketRes()
			pktR.receivedAt,
				pktR.addr,
				pktR.bytes,
				pktR.nBytes,
				pktR.ttl =
				time.Now(),
				addr,
				bytes,
				n,
				ttl
			if err = p.processPacket(pktR); nil != err {
				p.SetRevPacketRes(pktR)
				errIn := p.getErrRes()
				errIn.ping, errIn.err = nil, errors.New("conn6Rev.processPacket--->"+err.Error())
				p.errChan <- errIn
				continue
			}
			p.SetRevPacketRes(pktR)
		}
	}
END:
	p.readWg.Done()
}

// Attempts to match the ID of an ICMP packet.
func (p *Pool) matchID(ID int) bool {
	// On Linux we can only match ID if we are privileged.
	if p.option.protocol == "icmp" {
		return ID == p.option.pid
	}
	return true
}

// Returns the length of an ICMP message.
func (p *Pool) getMessageLength() int {
	return p.option.dataSize + 8
}

// processPacket process the packet
func (p *Pool) processPacket(rev *packet) (err error) {
	var (
		proto   int
		m       *icmp.Message
		ip      string
		pktUUID uuid.UUID
	)
	if rev.ipv4 {
		proto = protocolICMP
		rev.nBytes = stripIPv4Header(rev.nBytes, rev.bytes)
	} else {
		proto = protocolIPv6ICMP
	}

	if m, err = icmp.ParseMessage(proto, rev.bytes); err != nil {
		return fmt.Errorf("error parsing icmp message: %w", err)
	}

	if m.Type != ipv4.ICMPTypeEchoReply && m.Type != ipv6.ICMPTypeEchoReply {
		// Not an echo reply, ignore it
		return nil
	}

	switch pkt := m.Body.(type) {
	case *icmp.Echo:
		if !p.matchID(pkt.ID) {
			return nil
		}

		if p.option.protocol == Udp {
			if ip, _, err = net.SplitHostPort(rev.addr.String()); err != nil {
				return fmt.Errorf("err ip : %v, err %v", rev.addr, err)
			}
		} else {
			ip = rev.addr.String()
		}

		if len(pkt.Data) < timeSliceLength+trackerLength {
			return fmt.Errorf("insufficient data received; got: %d %v",
				len(pkt.Data), pkt.Data)
		}
		if pktUUID, err = p.getPacketUUID(pkt.Data); err != nil || pktUUID.String() == "" {
			return err
		}
		if _, ok := p.get(pktUUID.String()); !ok {
			return errors.New("not have the uuid")
		}
		pkg := p.getPacketRes()
		pkg.StartTime,
			pkg.EndTime,
			pkg.Addr,
			pkg.NBytes,
			pkg.Seq,
			pkg.ID,
			pkg.TTL,
			pkg.uuid =
			bytesToTime(pkt.Data[:timeSliceLength]),
			rev.receivedAt,
			ip,
			rev.nBytes,
			pkt.Seq,
			pkt.ID,
			rev.ttl,
			pktUUID.String()

		p.onRevChan <- pkg
	default:
		// Very bad, not sure how this can happen
		return fmt.Errorf("invalid ICMP echo reply; type: '%T', '%v'", pkt, pkt)
	}

	return
}
