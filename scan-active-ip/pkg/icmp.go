package pkg

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

// Scanner represents a ICMP scanner. It contains a pcap handle and
// other information that is needed to scan the network.
type Scanner struct {
	opt *OptionA

	// iFace is the network interface on which to scan.
	iFace *net.Interface

	// gw is the gateway address.
	gw net.IP
	// gwHardwareAddr is the gateway hardware address.
	gwHardwareAddr *net.HardwareAddr

	// src is the source IP address.
	src net.IP
	// handle is the pcap handle.
	handle *pcap.Handle

	// opts and buf allow us to easily serialize packets in the send method.
	opts   gopacket.SerializeOptions
	buf    gopacket.SerializeBuffer
	ips    chan string
	err    error
	errs   chan error
	wgSend *sync.WaitGroup
	wgRe   *sync.WaitGroup
	errSL  []error
}

// NewScanner creates a new Scanner.
func NewScanner(opt *OptionA) *Scanner {
	return &Scanner{
		iFace:          &net.Interface{},
		gw:             net.IP{},
		gwHardwareAddr: &net.HardwareAddr{},
		src:            net.IP{},
		handle:         &pcap.Handle{},
		err:            nil,
		opt:            opt,
		opts:           gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true},
		buf:            gopacket.NewSerializeBuffer(),
		wgSend:         &sync.WaitGroup{},
		wgRe:           &sync.WaitGroup{},
		ips:            make(chan string, opt.chanBuf),
		errs:           make(chan error, 20),
		errSL:          make([]error, 0, 20),
	}
}

func (s *Scanner) Do() {
	s.initParams()
	s.Scan()

}

// initParams init params
func (s *Scanner) initParams() {

	var (
		newRoute *Router
		iFace    net.Interface
	)

	//获取本机路由信息
	if newRoute, s.err = GetRouteInfo(); s.err != nil {
		return
	}

	// figure out the route by using the IP.
	if s.iFace, s.gw, s.src, s.err = newRoute.Route(net.ParseIP(s.opt.netMark)); nil != s.err {
		return
	}

	// open the handle for reading/writing.
	if s.handle, s.err = pcap.OpenLive(iFace.Name, 100, true, pcap.BlockForever); s.err != nil {
		return
	}

	if *s.gwHardwareAddr = s.getHwAddr(); s.err != nil {
		return
	}
	return
}

// getHwAddr gets the hardware address of the gateway by sending an ARP request.
func (s *Scanner) getHwAddr() net.HardwareAddr {
	if s.err != nil {
		return nil
	}
	start := time.Now()

	// prepare the layers to send for an ARP request.
	// send a single ARP request packet we never retry a sender
	if s.err = s.sendPackets(&layers.Ethernet{
		SrcMAC:       s.iFace.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}, &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   s.iFace.HardwareAddr,
		SourceProtAddress: s.src,
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    s.gw,
	}); s.err != nil {
		return nil
	}
	// wait 3 seconds for an ARP reply.
	for {
		if time.Since(start) > time.Second*3 {
			s.err = errors.New("timeout getting ARP reply")
			return nil
		}
		var data []byte
		if data, _, s.err = s.handle.ReadPacketData(); s.err == pcap.NextErrorTimeoutExpired {
			continue
		} else if s.err != nil {
			return nil
		}

		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
		if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
			arp := arpLayer.(*layers.ARP)
			if net.IP(arp.SourceProtAddress).Equal(s.gw) {
				return arp.SourceHwAddress
			}
		}
	}
}

func (s *Scanner) loopSend() {
	for v := range s.ips {
		s.send(v)
	}
	s.wgSend.Done()
}
func (s *Scanner) loopRec() {
	s.recv()
}

func (s *Scanner) Error() {
	for v := range s.errs {
		s.errSL = append(s.errSL, v)
	}
}

// sendPackets sends a packet with the given layers.
func (s *Scanner) sendPackets(l ...gopacket.SerializableLayer) error {
	if err := gopacket.SerializeLayers(s.buf, s.opts, l...); err != nil {
		return err
	}

	return s.handle.WritePacketData(s.buf.Bytes())
}

// Scan scans the network and returns a channel that contains the
// IP addresses of the hosts that respond to ICMP echo requests.
func (s *Scanner) Scan() {
	s.wgSend.Add(s.opt.worker)
	for i := 0; i < s.opt.worker; i++ {
		go s.loopSend()
	}

	s.wgRe.Add(s.opt.worker)
	for i := 0; i < s.opt.worker; i++ {
		go s.loopRec()
	}
	go s.Error()
}

// send sends a single ICMP echo request packet for each ip in the input channel.
func (s *Scanner) send(ip string) {
	id := uint16(os.Getpid())

	seq := uint16(0)
	dstIP := net.ParseIP(ip)
	if dstIP == nil {
		s.errs <- errors.New(ip + " parser is error")
		return
	}
	dstIP = dstIP.To4()
	if dstIP == nil {
		s.errs <- errors.New(ip + " ipv4 is error")
		return
	}

	// construct all the network layers we need.
	eth := layers.Ethernet{
		SrcMAC:       s.iFace.HardwareAddr,
		DstMAC:       *s.gwHardwareAddr,
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip4 := layers.IPv4{
		SrcIP:    s.src,
		DstIP:    dstIP.To4(),
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolICMPv4,
	}
	icmpLayer := layers.ICMPv4{
		TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
		Id:       id,
		Seq:      seq,
	}
	seq++

	err := s.sendPackets(&eth, &ip4, &icmpLayer)
	if err != nil {
		s.errs <- errors.New(ip + " parser is error")
		return
	}

}

// recv receives ICMP echo reply packets and sends the IP addresses
func (s *Scanner) recv() {
	defer s.wgRe.Done()

	// set the filter to only receive ICMP echo reply packets.
	s.handle.SetBPFFilter("dst host " + s.src.To4().String() + " and icmp")

	for {
		// read in the next packet.
		data, _, err := s.handle.ReadPacketData()
		if err == pcap.NextErrorTimeoutExpired {
			s.errs <- err
		} else if errors.Is(err, io.EOF) {
			s.errs <- err
			return
		} else if err != nil {
			s.errs <- err
			continue
		}

		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)

		// find the packets we care about, and print out logging
		// information about them.  All others are ignored.
		if net := packet.NetworkLayer(); net == nil {
			// log.Info("packet has no network layer")
			continue
		} else if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer == nil {
			// log.Info("packet has not ip layer")
			continue
		} else if ip, ok := ipLayer.(*layers.IPv4); !ok {
			continue
		} else if icmpLayer := packet.Layer(layers.LayerTypeICMPv4); icmpLayer == nil {
			// log.Info("packet has not icmp layer")
			continue
		} else if icmp, ok := icmpLayer.(*layers.ICMPv4); !ok {
			// log.Info("packet is not icmp")
			continue
		} else if icmp.TypeCode.Type() == layers.ICMPv4TypeEchoReply {
			fmt.Println("packet is not icmp")
			fmt.Println(ip.SrcIP.String(), "ip.SrcIP.String()")

		}
	}
}
