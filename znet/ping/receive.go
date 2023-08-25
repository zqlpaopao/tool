package ping

import (
	"errors"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"runtime"
	"time"
)

// revIpv4 -- --------------------------
// --> @Describe revIpv4 rev the ipv4 packet
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) revIpv4(conn net.PacketConn, id int) {
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
			if err = conn.SetReadDeadline(time.Now().Add(delay)); err != nil {
				errIn := p.errP.Get()
				errIn.Tag,
					errIn.Ping,
					errIn.Err =
					"revIpv4-SetReadDeadline",
					"",
					errors.New("conn4Rev.SetReadDeadline--->"+err.Error())
				p.errChan <- errIn
				continue
			}
			if n, addr, err = conn.ReadFrom(bytes); err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Timeout() {
						// Read timeout
						delay = expBackoffTime.Get()
						continue
					}
				}
				errIn := p.errP.Get()
				errIn.Tag,
					errIn.Ping,
					errIn.Err =
					"revIpv4-ReadFrom",
					"",
					errors.New("conn4Rev.SetReadDeadline--->"+err.Error())
				p.errChan <- errIn
				continue
			}

			//packet := gopacket.NewPacket(bytes, layers.LayerTypeICMPv4, gopacket.NoCopy)
			//packet = packet
			//fmt.Println("revIp4", id, packet)
			//os.Exit(1)
			//	// Get the UDP layer from this packet
			//	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			//		udp, _ := udpLayer.(*layers.UDP)
			//		if app := packet.ApplicationLayer(); app != nil {
			//			data, err := EncodeUDPPacket(net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.1"), uint16(udp.DstPort), uint16(udp.SrcPort), app.Payload())
			//			if err != nil {
			//				log.Printf("failed to EncodePacket: %v", err)
			//				return
			//			}
			//			if _, err := conn.WriteTo(data, remoteaddr); err != nil {
			//				log.Printf("failed to write packet: %v", err)
			//				conn.Close()
			//				return
			//			}
			//		}
			//	}
			//}

			if err = p.processPacket(time.Now(),
				addr,
				bytes,
				n,
				ttl,
				true); nil != err {
				errIn := p.errP.Get()
				errIn.Tag,
					errIn.Ping,
					errIn.Err =
					"revIpv4-ReadFrom",
					"",
					errors.New("conn4Rev.processPacket--->"+err.Error())
				p.errChan <- errIn

				continue
			}
		}
	}
END:
	_ = conn.Close()
	p.wgRec.Done()
}

// revIpv6 -- --------------------------
// --> @Describe revIpv6 rev the ipv6 packet
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) revIpv6(conn net.PacketConn, id int) {
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
			if err = conn.SetReadDeadline(time.Now().Add(delay)); err != nil {

				errIn := p.errP.Get()
				errIn.Tag,
					errIn.Ping,
					errIn.Err =
					"revIpv4-ReadFrom",
					"",
					errors.New("conn6Rev.SetReadDeadline--->"+err.Error())
				p.errChan <- errIn

				continue
			}
			if n, addr, err = conn.ReadFrom(bytes); err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Timeout() {
						// Read timeout
						delay = expBackoffTime.Get()
						continue
					}
				}

				errIn := p.errP.Get()
				errIn.Tag,
					errIn.Ping,
					errIn.Err =
					"revIpv4-ReadFrom",
					"",
					errors.New("conn6Rev.ReadFrom--->"+err.Error())
				p.errChan <- errIn

				continue
			}
			if err = p.processPacket(time.Now(),
				addr,
				bytes,
				n,
				ttl,
				false); nil != err {

				errIn := p.errP.Get()
				errIn.Tag,
					errIn.Ping,
					errIn.Err =
					"revIpv6-ReadFrom",
					"",
					errors.New("conn6Rev.processPacket--->"+err.Error())
				p.errChan <- errIn

				continue
			}
		}
	}
END:
	_ = conn.Close()
	p.wgRec.Done()
}

// matchID-- --------------------------
// --> @Describe  Attempts to match the ID of an ICMP packet.
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) matchID(ID int, id int) bool {
	// On Linux we can only match ID if we are privileged.
	if p.option.protocol == "icmp" {
		return ID == id
	}
	return true
}

// getMessageLength -- --------------------------
// --> @Describe Returns the length of an ICMP message.
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) getMessageLength() int {
	return p.option.dataSize + 8
}

// processPacket -- --------------------------
// --> @Describe process the packet
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) processPacket(receivedAt time.Time,
	addr net.Addr,
	bytes []byte,
	nBytes int,
	ttl int,
	ipV4 bool) (err error) {
	var (
		proto   int
		m       *icmp.Message
		ip      string
		pktUUID string
	)

	if ipV4 {
		proto = protocolICMP
		nBytes = stripIPv4Header(nBytes, bytes)
	} else {
		proto = protocolIPv6ICMP
	}

	if m, err = icmp.ParseMessage(proto, bytes); err != nil {
		return fmt.Errorf("error parsing icmp message: %w", err)
	}

	if m.Type != ipv4.ICMPTypeEchoReply && m.Type != ipv6.ICMPTypeEchoReply {
		// Not an echo reply, ignore it
		return nil
	}

	switch pkt := m.Body.(type) {
	case *icmp.Echo:
		if !p.matchID(pkt.ID, p.option.pid) {
			return nil
		}

		if p.option.protocol == Udp {
			if ip, _, err = net.SplitHostPort(addr.String()); err != nil {
				return fmt.Errorf("err ip : %v, err %v", addr, err)
			}
		} else {
			ip = addr.String()
		}

		if len(pkt.Data) < timeSliceLength+trackerLength {
			return fmt.Errorf("insufficient data received; got: %d %v",
				len(pkt.Data), pkt.Data)
		}
		if pktUUID = p.getPacketUUID(pkt.Data); err != nil {
			return err
		}

		if _, ok := p.get(pktUUID); !ok {
			return errors.New("not have the uuid")
		}
		pkg := p.PacketP.Get()
		pkg.StartTime,
			pkg.EndTime,
			pkg.Addr,
			pkg.NBytes,
			pkg.Seq,
			pkg.ID,
			pkg.TTL,
			pkg.uuid =
			bytesToTime(pkt.Data[:timeSliceLength]),
			receivedAt,
			ip,
			nBytes,
			pkt.Seq,
			pkt.ID,
			ttl,
			pktUUID

		p.resChan <- pkg

	default:
		// Very bad, not sure how this can happen
		return fmt.Errorf("invalid ICMP echo reply; type: '%T', '%v'", pkt, pkt)
	}

	return
}
