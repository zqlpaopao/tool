package pkg

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"sync/atomic"
	"time"
)

func (p *PingMan) receiveIpv4() {
	defer p.opt.saveRecover()
	var (
		err        error
		n          int
		cm         *ipv4.ControlMessage
		addr       net.Addr
		ttl        = 0
		receivePkg *packet
		tag        string
	)
	for {
		select {

		case _, ok := <-p.close:
			if !ok {
				goto END
			}

		default:
			if err = p.conn4.SetReadDeadline(time.Now().Add(p.opt.connReadTimeout)); nil != err && p.opt.debug {
				if p.opt.debug {
					color.Red("set conn4  SetReadDeadline err-%v, connReadTimeout is-%d\n", err, p.opt.connReadTimeout)
				}
			}
			bytes := make([]byte, 512)
			n, cm, addr, err = p.conn4.IPv4PacketConn().ReadFrom(bytes)
			if cm != nil {
				ttl = cm.TTL
			}
			if err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Timeout() {
						continue
					}
					if p.opt.debug {
						color.Red("conn4  IPv4PacketConn ReadFrom err-%s\n", err)
					}
					continue

				}
			}

			receivePkg = &packet{bytes: bytes, nBytes: n, ttl: ttl, proto: protoIpv4, addr: addr}

			if tag, err = p.processPacket(receivePkg); err != nil && p.opt.debug {
				color.Red("conn4  IPv4PacketConn processPacket err-tag err-%v, receivePkg-%v \n", tag, err, receivePkg)
			}
		}
	}
END:
	p.wg.Done()

}

func (p *PingMan) receiveIpv6() {
	defer p.opt.saveRecover()
	var (
		err        error
		n          int
		cm         *ipv6.ControlMessage
		addr       net.Addr
		ttl        = 0
		bytes      = make([]byte, 512)
		receivePkg *packet
		tag        string
	)
	for {
		select {

		case _, ok := <-p.close:
			if !ok {
				goto END
			}

		default:

			if err = p.conn6.SetReadDeadline(time.Now().Add(p.opt.connReadTimeout)); nil != err {
				if p.opt.debug {
					color.Red("set conn4  SetReadDeadline err-%v, connReadTimeout is-%d\n", err, p.opt.connReadTimeout)
				}
			}

			n, cm, addr, err = p.conn6.IPv6PacketConn().ReadFrom(bytes)
			if cm != nil {
				ttl = cm.HopLimit
			}
			if err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Timeout() {
						// Read timeout
						continue
					}
					if p.opt.debug {
						color.Red("conn4  IPv6PacketConn ReadFrom err-%s\n", err)
					}
					continue
				}
			}

			receivePkg = &packet{bytes: bytes, nBytes: n, ttl: ttl, proto: protoIpv6, addr: addr}

			if tag, err = p.processPacket(receivePkg); err != nil && p.opt.debug {
				color.Red("conn6  IPv6PacketConn processPacket err tag-%s err-%v, receivePkg-%v \n", tag, err, receivePkg)
			}
		}

	}
END:
	p.wg.Done()

}

func (p *PingMan) processPacket(receive *packet) (tag string, err error) {
	var (
		proto      int
		m          *icmp.Message
		receivedAt = time.Now()
		ip         string
	)
	if receive.proto == protoIpv4 {
		proto = protocolICMP
	} else {
		proto = protocolIPv6ICMP
	}

	if m, err = icmp.ParseMessage(proto, receive.bytes); err != nil {
		return "icmp.ParseMessage", err
	}

	if m.Type != ipv4.ICMPTypeEchoReply && m.Type != ipv6.ICMPTypeEchoReply {
		if p.opt.debug {
			color.White("pkg drop-%v \n", m)
		}
		return
	}

	switch pkt := m.Body.(type) {
	case *icmp.Echo:

		// If we are privileged, we can match icmp.ID
		if pkt.ID != p.id {
			if p.opt.debug {
				color.White("drop pkg-%+v id-%v addr-%s \n", pkt, p.id, receive.addr)
			}
			return
		}
		if len(pkt.Data) < timeSliceLength+trackerLength {
			return fmt.Sprintf("data received len %d  less-%d data-%v", len(pkt.Data), timeSliceLength+trackerLength, pkt.Data), errors.New("received data less")
		}

		timestamp := bytesToTime(pkt.Data[:timeSliceLength])

		if p.opt.network == Udp.String() {
			if ip, _, err = net.SplitHostPort(receive.addr.String()); err != nil {
				return fmt.Sprintf("err ip-%v", receive.addr), err
			}
		} else {
			ip = receive.addr.String()
		}

		p.resultQueue <- &ResPing{
			seqID:      pkt.ID,
			pid:        pkt.ID,
			ip:         ip,
			receivedAt: receivedAt.Sub(timestamp),
		}
		atomic.AddInt64(&p.current, 1)

	default:
		// Very bad, not sure how this can happen
		return fmt.Sprintf("invalid ICMP echo reply; type: '%T', '%v'", pkt, pkt), errors.New("parser is err")
	}

	return

}
