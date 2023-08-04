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

func (p *Pool) revIpv4() {
	defer p.option.Recover()
	// Start by waiting for 50 µs and increase to a possible maximum of ~ 100 ms.
	expBackoffTime := newExpBackoff(50*time.Microsecond, 11)
	delay := expBackoffTime.Get()

	// Workaround for https://github.com/golang/go/issues/47369
	offset := 0
	if p.option.protocol != "icmp" && runtime.GOOS == "darwin" {
		offset = 20
	}
	for {
		select {
		case <-p.done:
			goto END
		default:
			bytes := make([]byte, p.getMessageLength()+offset)
			if err := p.conn4Rev.SetReadDeadline(time.Now().Add(delay)); err != nil {
				p.errChan <- &ErrInfo{
					ping: nil,
					err:  err,
				}
				continue
			}
			var (
				n, ttl int
				err    error
				addr   net.Addr
			)

			n, ttl, addr, err = p.conn4Rev.ReadFrom(bytes)
			if err != nil {
				p.errChan <- &ErrInfo{
					ping: nil,
					err:  err,
				}
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Timeout() {
						// Read timeout
						delay = expBackoffTime.Get()
						continue
					}
				}
				continue
			}

			if err := p.processPacket(&packet{addr: addr, bytes: bytes, nBytes: n, ttl: ttl, ipv4: true}); nil != err {
				p.errChan <- &ErrInfo{ping: nil, err: err}
				continue
			}
		}
	}
END:
	p.readWg.Done()
}

func (p *Pool) revIpv6() {
	defer p.option.Recover()
	// Start by waiting for 50 µs and increase to a possible maximum of ~ 100 ms.
	expBackoffTime := newExpBackoff(50*time.Microsecond, 11)
	delay := expBackoffTime.Get()

	// Workaround for https://github.com/golang/go/issues/47369
	offset := 0
	if p.option.protocol != "icmp" && runtime.GOOS == "darwin" {
		offset = 20
	}

	for {
		select {
		case <-p.done:
			goto END
		default:
			bytes := make([]byte, p.getMessageLength()+offset)
			if err := p.conn6Rev.SetReadDeadline(time.Now().Add(delay)); err != nil {
				p.errChan <- &ErrInfo{
					ping: nil,
					err:  err,
				}
				continue
			}
			var (
				n, ttl int
				err    error
				addr   net.Addr
			)
			n, ttl, addr, err = p.conn6Rev.ReadFrom(bytes)
			if err != nil {
				p.errChan <- &ErrInfo{
					ping: nil,
					err:  err,
				}
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Timeout() {
						// Read timeout
						delay = expBackoffTime.Get()
						continue
					}
				}
				continue
			}

			if err := p.processPacket(&packet{addr: addr, bytes: bytes, nBytes: n, ttl: ttl}); nil != err {
				p.errChan <- &ErrInfo{ping: nil, err: err}
				continue
			}
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
func (p *Pool) processPacket(rev *packet) (err error) {
	receivedAt := time.Now()
	var (
		proto int
		m     *icmp.Message
	)
	if rev.ipv4 {
		proto = protocolICMP
		// Workaround for https://github.com/golang/go/issues/47369
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
		var (
			ip      string
			pktUUID uuid.UUID
		)
		if p.option.protocol == "udp" {
			if ip, _, err = net.SplitHostPort(rev.addr.String()); err != nil {
				return fmt.Errorf("err ip : %v, err %v", rev.addr, err)
			}
		} else {
			ip = rev.addr.String()
		}
		if !p.matchID(pkt.ID) {
			return nil
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
		p.onRevChan <- &Packet{
			StartTime: bytesToTime(pkt.Data[:timeSliceLength]),
			EndTime:   receivedAt,
			//Rtt:       time.Duration(rev.ttl),
			//IPAddr: rev.addr,
			Addr:   ip,
			NBytes: rev.nBytes,
			Seq:    pkt.Seq,
			ID:     pkt.ID,
			TTL:    rev.ttl,
			uuid:   pktUUID.String(),
		}
	default:
		// Very bad, not sure how this can happen
		return fmt.Errorf("invalid ICMP echo reply; type: '%T', '%v'", pkt, pkt)
	}

	return nil
}
