package pkg

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"time"
)

func (p *PingPool) revIpv4() {
	var ttl int
	for {
		select {
		case _, ok := <-p.done:
			if !ok {
				goto END
			}
		default:
			bytes := make([]byte, 512)
			if err := p.conn4.SetReadDeadline(time.Now().Add(p.option.readTimeout)); nil != err {
				p.errChan <- &ErrInfo{ping: nil, err: err}
			}
			n, cm, addr, err := p.conn4.IPv4PacketConn().ReadFrom(bytes)
			if cm != nil {
				ttl = cm.TTL
			}

			if err != nil {
				if netRrr, ok := err.(*net.OpError); ok {
					if netRrr.Timeout() {
						// Read timeout
						continue
					}
					p.errChan <- &ErrInfo{ping: nil, err: err}
					//return
				}
			}
			revPkg := &packet{bytes: bytes, nBytes: n, ttl: ttl, proto: protoIpv4, addr: addr}

			if err =
				p.processPacket(revPkg); err != nil {
				p.errChan <- &ErrInfo{ping: nil, err: err}
			}
		}
	}
END:
	p.readWg.Done()
}

func (p *PingPool) revIpv6() {
	var ttl int
	for {
		select {
		case _, ok := <-p.done:
			if !ok {
				goto END
			}
		default:
			bytes := make([]byte, 512)
			if err := p.conn6.SetReadDeadline(time.Now().Add(p.option.readTimeout)); nil != err {
				p.errChan <- &ErrInfo{ping: nil, err: err}
			}
			n, cm, addr, err := p.conn6.IPv6PacketConn().ReadFrom(bytes)
			if cm != nil {
				ttl = cm.HopLimit
			}
			if err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Timeout() {
						// Read timeout
						continue
					}
					p.errChan <- &ErrInfo{ping: nil, err: err}
				}
			}

			revPkg := &packet{bytes: bytes, nBytes: n, ttl: ttl, proto: protoIpv6, addr: addr}

			if err =
				p.processPacket(revPkg); err != nil {
				p.errChan <- &ErrInfo{ping: nil, err: err}
			}
		}

	}
END:
	p.readWg.Done()
}

func (p *PingPool) processPacket(rev *packet) (err error) {
	receivedAt := time.Now()
	var (
		proto int
		m     *icmp.Message
	)
	if rev.proto == protoIpv4 {
		proto = protocolICMP
	} else {
		proto = protocolIPv6ICMP
	}

	if m, err = icmp.ParseMessage(proto, rev.bytes); err != nil {
		return fmt.Errorf("error parsing icmp message: %s", err.Error())
	}

	if m.Type != ipv4.ICMPTypeEchoReply && m.Type != ipv6.ICMPTypeEchoReply {
		// Not an echo reply, ignore it
		//if bp.option.debug {
		//	log.Printf("pkg drop %v \n", m)
		//}
		return nil
	}

	switch pkt := m.Body.(type) {
	case *icmp.Echo:
		// If we are privileged, we can match icmp.ID
		if pkt.ID != p.option.pid {
			return nil
		}
		if len(pkt.Data) < timeSliceLength+trackerLength {
			return fmt.Errorf("insufficient data received; got: %d %v",
				len(pkt.Data), pkt.Data)
		}

		var ip string
		if p.option.network == "udp" {
			if ip, _, err = net.SplitHostPort(rev.addr.String()); err != nil {
				return fmt.Errorf("err ip : %v, err %v", rev.addr, err)
			}
		} else {
			ip = rev.addr.String()
		}

		p.onRevChan <- &Packet{
			StartTime: bytesToTime(pkt.Data[:timeSliceLength]),
			EndTime:   receivedAt,
			//Rtt:       time.Duration(rev.ttl),
			IPAddr: ip,
			Addr:   rev.addr.String(),
			NBytes: rev.nBytes,
			Seq:    pkt.Seq,
			Ttl:    rev.ttl,
			ID:     pkt.ID,
			Bytes:  rev.bytes,
		}
	default:
		// Very bad, not sure how this can happen
		return fmt.Errorf("invalid ICMP echo reply; type: '%T', '%v'", pkt, pkt)
	}

	return nil

}
