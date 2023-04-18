package main

import "github.com/zqlpaopao/tool/scan-active-ip/pkg"

func main() {

	//usageTime := time.Now()
	//返回端口map
	//openPort := dealPort(*port)

	//isLocalIP := false
	//
	//for _, addr := range newRoute.addrs {
	//	if addr.String() == *ipString {
	//		isLocalIP = true
	//		fmt.Println("要扫描的IP: ", addr.String())
	//		break
	//	}

	//
	////判断是否为本机IP
	//if isLocalIP {
	//	scanLocal(openPort)
	//} else {
	//	syncScan(newRoute, openPort)
	//}
	//fmt.Printf("Usage time: %v\n\n", time.Since(usageTime))
	//return
	pkg.GetRouteInfo()

	//var c = make(chan []string, 10000)
	//
	//scanner := NewScanner()
	//scanner.Scan(c)
	//
	//c <- []string{"8.8.8.8/24"}

	//err := exec.Command("/bin/bash", "ip.sh").Run()
	//if err != nil {
	//	golog.Fatal(err)
	//}
	//golog.Println("11")
	//defer os.Remove("ipv4.txt")
	//defer os.Remove("delegated-apnic-latest")
	//
	//ipList, err := os.Open("ipv4.txt")
	//if err != nil {
	//	panic(err)
	//}
	//
	//scanner := bufio.NewScanner(ipList)
	//for scanner.Scan() {
	//	//netmask := scanner.Text()
	//	//scan-active-ip(netmask)
	//}
}

//// getHwAddr gets the hardware address of the gateway by sending an ARP request.
//func (s *Scanner) getHwAddr() (net.HardwareAddr, error) {
//	start := time.Now()
//	arpDst := s.gw
//
//	// prepare the layers to send for an ARP request.
//	eth := layers.Ethernet{
//		SrcMAC:       s.iface.HardwareAddr,
//		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
//		EthernetType: layers.EthernetTypeARP,
//	}
//	arp := layers.ARP{
//		AddrType:          layers.LinkTypeEthernet,
//		Protocol:          layers.EthernetTypeIPv4,
//		HwAddressSize:     6,
//		ProtAddressSize:   4,
//		Operation:         layers.ARPRequest,
//		SourceHwAddress:   []byte(s.iface.HardwareAddr),
//		SourceProtAddress: []byte(s.src),
//		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
//		DstProtAddress:    []byte(arpDst),
//	}
//	// send a single ARP request packet (we never retry a send)
//	if err := s.sendPackets(&eth, &arp); err != nil {
//		return nil, err
//	}
//	// wait 3 seconds for an ARP reply.
//	for {
//		if time.Since(start) > time.Second*3 {
//			return nil, errors.New("timeout getting ARP reply")
//		}
//		data, _, err := s.handle.ReadPacketData()
//		if err == pcap.NextErrorTimeoutExpired {
//			continue
//		} else if err != nil {
//			return nil, err
//		}
//		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
//		if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
//			arp := arpLayer.(*layers.ARP)
//			if net.IP(arp.SourceProtAddress).Equal(net.IP(arpDst)) {
//				return net.HardwareAddr(arp.SourceHwAddress), nil
//			}
//		}
//	}
//}
//
//// sendPackets sends a packet with the given layers.
//func (s *Scanner) sendPackets(l ...gopacket.SerializableLayer) error {
//	if err := gopacket.SerializeLayers(s.buf, s.opts, l...); err != nil {
//		return err
//	}
//
//	return s.handle.WritePacketData(s.buf.Bytes())
//}
//
//// Scan scans the network and returns a channel that contains the
//// IP addresses of the hosts that respond to ICMP echo requests.
//func (s *Scanner) Scan(input chan []string) (output chan string) {
//	output = make(chan string, 1024*1024)
//	go s.recv(output)
//	go s.send(input)
//
//	return output
//}
//
//// send sends a single ICMP echo request packet for each ip in the input channel.
//func (s *Scanner) send(input chan []string) error {
//	id := uint16(os.Getpid())
//
//	seq := uint16(0)
//	for ips := range input {
//		for _, ip := range ips {
//			dstIP := net.ParseIP(ip)
//			if dstIP == nil {
//				continue
//			}
//			dstIP = dstIP.To4()
//			if dstIP == nil {
//				continue
//			}
//
//			// construct all the network layers we need.
//			eth := layers.Ethernet{
//				SrcMAC:       s.iface.HardwareAddr,
//				DstMAC:       *s.gwHardwareAddr,
//				EthernetType: layers.EthernetTypeIPv4,
//			}
//			ip4 := layers.IPv4{
//				SrcIP:    s.src,
//				DstIP:    dstIP.To4(),
//				Version:  4,
//				TTL:      64,
//				Protocol: layers.IPProtocolICMPv4,
//			}
//			icmpLayer := layers.ICMPv4{
//				TypeCode: layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0),
//				Id:       id,
//				Seq:      seq,
//			}
//			seq++
//
//			err := s.sendPackets(&eth, &ip4, &icmpLayer)
//			if err != nil {
//				log.Fatal(err)
//			}
//		}
//	}
//
//	return nil
//}
//
//// recv receives ICMP echo reply packets and sends the IP addresses
//func (s *Scanner) recv(output chan string) {
//	defer close(output)
//
//	// set the filter to only receive ICMP echo reply packets.
//	s.handle.SetBPFFilter("dst host " + s.src.To4().String() + " and icmp")
//
//	for {
//		// read in the next packet.
//		data, _, err := s.handle.ReadPacketData()
//		if err == pcap.NextErrorTimeoutExpired {
//			continue
//		} else if errors.Is(err, io.EOF) {
//			// log.Infof("error reading packet: %v", err)
//			return
//		} else if err != nil {
//			log.Printf("error reading packet: %v", err)
//			continue
//		}
//
//		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
//
//		// find the packets we care about, and print out logging
//		// information about them.  All others are ignored.
//		if net := packet.NetworkLayer(); net == nil {
//			// log.Info("packet has no network layer")
//			continue
//		} else if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer == nil {
//			// log.Info("packet has not ip layer")
//			continue
//		} else if ip, ok := ipLayer.(*layers.IPv4); !ok {
//			continue
//		} else if icmpLayer := packet.Layer(layers.LayerTypeICMPv4); icmpLayer == nil {
//			// log.Info("packet has not icmp layer")
//			continue
//		} else if icmp, ok := icmpLayer.(*layers.ICMPv4); !ok {
//			// log.Info("packet is not icmp")
//			continue
//		} else if icmp.TypeCode.Type() == layers.ICMPv4TypeEchoReply {
//			fmt.Println("packet is not icmp")
//			fmt.Println(ip.SrcIP.String(), "ip.SrcIP.String()")
//			select {
//			case output <- ip.SrcIP.String():
//			default:
//			}
//		} else {
//			// log.Info("ignoring useless packet")
//		}
//	}
//}
//
////routing
//
//var (
//	ipString     = flag.String("ip", "127.0.0.1", "ip or domain name")
//	port         = flag.String("p", "22-1000", "port")
//	wg           = &sync.WaitGroup{}
//	showOpenPort []int
//)
//
//type rtInfo struct {
//	Dst              net.IPNet
//	Gateway, PrefSrc net.IP
//	OutputIface      uint32
//	Priority         uint32
//}
//
//type routeSlice []*rtInfo
//
//type router struct {
//	ifaces []net.Interface
//	addrs  []net.IP
//	v4     routeSlice
//}
//
//func getRouteInfo() (*router, error) {
//	rtr := &router{}
//
//	tab, err := syscall.NetlinkRIB(syscall.RTM_GETROUTE, syscall.AF_INET)
//	if err != nil {
//		return nil, err
//	}
//
//	msgs, err := syscall.ParseNetlinkMessage(tab)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, m := range msgs {
//		switch m.Header.Type {
//		case syscall.NLMSG_DONE:
//			break
//		case syscall.RTM_NEWROUTE:
//			rtmsg := (*syscall.RtMsg)(unsafe.Pointer(&m.Data[0]))
//			attrs, err := syscall.ParseNetlinkRouteAttr(&m)
//			if err != nil {
//				return nil, err
//			}
//
//			routeInfo := rtInfo{}
//			rtr.v4 = append(rtr.v4, &routeInfo)
//			for _, attr := range attrs {
//				switch attr.Attr.Type {
//				case syscall.RTA_DST:
//					routeInfo.Dst.IP = net.IPv4(attr.Value[0], attr.Value[1], attr.Value[2], attr.Value[3])
//					routeInfo.Dst.Mask = net.CIDRMask(int(rtmsg.Dst_len), len(attr.Value)*8)
//
//				case syscall.RTA_GATEWAY:
//					routeInfo.Gateway = net.IPv4(attr.Value[0], attr.Value[1], attr.Value[2], attr.Value[3])
//
//				case syscall.RTA_OIF:
//					routeInfo.OutputIface = *(*uint32)(unsafe.Pointer(&attr.Value[0]))
//
//				case syscall.RTA_PRIORITY:
//					routeInfo.Priority = *(*uint32)(unsafe.Pointer(&attr.Value[0]))
//
//				case syscall.RTA_PREFSRC:
//					routeInfo.PrefSrc = net.IPv4(attr.Value[0], attr.Value[1], attr.Value[2], attr.Value[3])
//
//				}
//			}
//		}
//	}
//
//	sort.Slice(rtr.v4, func(i, j int) bool {
//		return rtr.v4[i].Priority < rtr.v4[j].Priority
//	})
//
//	ifaces, err := net.Interfaces()
//	if err != nil {
//		return nil, err
//	}
//
//	for i, iface := range ifaces {
//		if i != iface.Index-1 {
//			break
//		}
//		if iface.Flags&net.FlagUp == 0 {
//			continue
//		}
//
//		rtr.ifaces = append(rtr.ifaces, iface)
//		ifaceAddrs, err := iface.Addrs()
//		if err != nil {
//			return nil, err
//		}
//		var addrs net.IP
//		for _, addr := range ifaceAddrs {
//			if inet, ok := addr.(*net.IPNet); ok {
//				if v4 := inet.IP.To4(); v4 != nil {
//					if addrs == nil {
//						addrs = v4
//					}
//				}
//			}
//		}
//		rtr.addrs = append(rtr.addrs, addrs)
//	}
//	return rtr, nil
//}
//
//func (r *router) Route(dst net.IP) (iface net.Interface, gateway, prefsrc net.IP, err error) {
//	for _, rt := range r.v4 {
//		if rt.Dst.IP != nil && !rt.Dst.Contains(dst) {
//			continue
//		}
//		iface = r.ifaces[rt.OutputIface-1]
//		gateway = rt.Gateway.To4()
//		if rt.PrefSrc == nil {
//			prefsrc = r.addrs[rt.OutputIface-1]
//		} else {
//			prefsrc = rt.PrefSrc.To4()
//		}
//		return
//	}
//	err = errors.New("No route found!")
//	return
//}
//
//func dealRepeat(port []int) []int {
//	result := make([]int, 0)
//	tempMap := make(map[int]bool, len(port))
//	for _, v := range port {
//		if tempMap[v] == false {
//			tempMap[v] = true
//			result = append(result, v)
//		}
//	}
//	sort.Slice(result, func(i, j int) bool {
//		return result[i] < result[j]
//	})
//	return result
//}
//
//func getAllPort(port string) []int {
//	var all []int
//	port = strings.Trim(port, ", ")
//	portArr := strings.Split(port, ",")
//	for _, v := range portArr {
//		v = strings.Trim(v, " ")
//
//		if strings.Contains(v, "-") {
//			data := strings.Split(v, "-")
//			firstPort, _ := strconv.Atoi(data[0])
//			lastPort, _ := strconv.Atoi(data[1])
//
//			for i := firstPort; i <= lastPort; i++ {
//				if i < 1 || i > 65535 {
//					log.Fatal("port illegal!")
//				}
//				all = append(all, i)
//			}
//		} else {
//			data, _ := strconv.Atoi(v)
//			if data < 1 || data > 65535 {
//				log.Fatal("port illegal!")
//			}
//			all = append(all, data)
//		}
//	}
//	//fmt.Println(all)
//	return all
//}
//
//func getHwAddr(ip, gateway, srcIP net.IP, networkInterface *net.Interface, handle *pcap.Handle) (net.HardwareAddr, error) {
//	arpDst := ip
//	if gateway != nil {
//		arpDst = gateway
//	}
//
//	// handle, err := pcap.OpenLive(networkInterface.Name, 4096, false, pcap.BlockForever)
//	// if err != nil {
//	//  return nil, err
//	// }
//	// defer handle.Close()
//
//	eth := layers.Ethernet{
//		SrcMAC:       networkInterface.HardwareAddr,
//		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
//		EthernetType: layers.EthernetTypeARP,
//	}
//
//	arp := layers.ARP{
//		AddrType:          layers.LinkTypeEthernet,
//		Protocol:          layers.EthernetTypeIPv4,
//		HwAddressSize:     uint8(6),
//		ProtAddressSize:   uint8(4),
//		Operation:         layers.ARPRequest,
//		SourceHwAddress:   []byte(networkInterface.HardwareAddr),
//		SourceProtAddress: srcIP,
//		DstHwAddress:      net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
//		DstProtAddress:    arpDst,
//	}
//
//	opt := gopacket.SerializeOptions{
//		FixLengths:       true,
//		ComputeChecksums: true,
//	}
//
//	//var opt gopacket.SerializeOptions
//	buf := gopacket.NewSerializeBuffer()
//
//	if err := gopacket.SerializeLayers(buf, opt, &eth, &arp); err != nil {
//		return nil, err
//	}
//	if err := handle.WritePacketData(buf.Bytes()); err != nil {
//		return nil, err
//	}
//
//	start := time.Now()
//	for {
//		if time.Since(start) > time.Millisecond*time.Duration(1000) {
//			return nil, errors.New("timeout getting ARP reply")
//		}
//		data, _, err := handle.ReadPacketData()
//		if err == pcap.NextErrorTimeoutExpired {
//			continue
//		} else if err != nil {
//			return nil, err
//		}
//		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)
//		if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
//			arp := arpLayer.(*layers.ARP)
//			if net.IP(arp.SourceProtAddress).Equal(arpDst) {
//				return net.HardwareAddr(arp.SourceHwAddress), nil
//			}
//		}
//
//	}
//}
//
//func dealPort(port string) map[int][]int {
//	var (
//		openPort         = make(map[int][]int)
//		processCount int = 100
//		interval     int
//	)
//
//	ports := getAllPort(port)
//
//	if len(ports) <= processCount {
//		interval = 1
//		processCount = len(ports)
//	} else {
//		interval = int(math.Ceil(float64(len(ports)) / float64(processCount)))
//	}
//
//	for i := 0; i < processCount; i++ {
//		for j := 0; j < interval; j++ {
//			temp := i*interval + j
//			if temp < len(ports) {
//				openPort[i] = append(openPort[i], ports[temp])
//			}
//		}
//	}
//
//	return openPort
//}
//
//func getOpenPort(ip string, port int) bool {
//	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), time.Millisecond*time.Duration(100))
//	if err != nil {
//		return false
//	}
//	defer conn.Close()
//	return true
//}
//
//func scanLocal(openPort map[int][]int) {
//	for i := 0; i < 1; i++ {
//		for _, value := range openPort {
//			wg.Add(1)
//			go func(v []int) {
//				defer wg.Done()
//				for _, num := range v {
//					if getOpenPort(*ipString, num) {
//						showOpenPort = append(showOpenPort, num)
//					}
//				}
//			}(value)
//		}
//	}
//	wg.Wait()
//	fmt.Println("开放的端口: ", dealRepeat(showOpenPort))
//	return
//}
//
//func syncScan(route *router, openPort map[int][]int) {
//	//获取本机一个没有使用的端口
//	//rawPort, err := freeport.GetFreePort()
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//
//	//realIP, err := net.ResolveIPAddr("ip", *ipString)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//dstIP := (*realIP).IP.To4()
//	//fmt.Println("要扫描的IP: ", dstIP.String())
//	//
//	////获取本地路由
//	//networkInterface, gateway, srcIP, err := route.Route(dstIP)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//
//	//handle, err := pcap.OpenLive(networkInterface.Name, 65535, true, pcap.BlockForever)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//defer handle.Close()
//	//
//	////获取网关MAC或者()目的地址MAC
//	//hwaddr, err := getHwAddr(dstIP, gateway, srcIP, &networkInterface, handle)
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//rand.Seed(time.Now().Unix())
//	//eth := layers.Ethernet{
//	//	SrcMAC:       networkInterface.HardwareAddr,
//	//	DstMAC:       hwaddr,
//	//	EthernetType: layers.EthernetTypeIPv4,
//	//}
//	//ip4 := layers.IPv4{
//	//	SrcIP:    srcIP,
//	//	DstIP:    dstIP,
//	//	Version:  4,
//	//	TTL:      255,
//	//	Id:       uint16(rand.Intn(1 << 16)),
//	//	Protocol: layers.IPProtocolTCP,
//	//	Flags:    layers.IPv4DontFragment,
//	//}
//	//bsTsval := make([]byte, 4)
//	//bsTsecr := make([]byte, 4)
//	//binary.BigEndian.PutUint32(bsTsval, uint32(time.Now().UnixNano()))
//	//bsTime := append(bsTsval, bsTsecr...)
//	//tcp := layers.TCP{
//	//	SrcPort: layers.TCPPort(rawPort),
//	//	DstPort: 0,
//	//	SYN:     true,
//	//	Window:  43690,
//	//	Seq:     uint32(rand.Intn(1 << 32)),
//	//	Options: []layers.TCPOption{
//	//		{
//	//			OptionType:   layers.TCPOptionKindMSS,
//	//			OptionLength: 4,
//	//			OptionData:   []byte{0xff, 0xd7},
//	//		},
//	//		{
//	//			OptionType:   layers.TCPOptionKindSACKPermitted,
//	//			OptionLength: 2,
//	//		},
//	//		{
//	//			OptionType:   layers.TCPOptionKindTimestamps,
//	//			OptionLength: 10,
//	//			OptionData:   bsTime,
//	//		},
//	//	},
//	//}
//	//
//	//tcp.SetNetworkLayerForChecksum(&ip4)
//	//ch := make(chan bool)
//	//quitCh := make(chan bool)
//	//
//	//opt := gopacket.SerializeOptions{
//	//	FixLengths:       true,
//	//	ComputeChecksums: true,
//	//}
//	//
//	//go func() {
//	//	ethRecv := &layers.Ethernet{}
//	//	ip4Recv := &layers.IPv4{}
//	//	tcpRecv := &layers.TCP{}
//	//
//	//	ipFlow := gopacket.NewFlow(layers.EndpointIPv4, dstIP, srcIP)
//	//	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, ethRecv, ip4Recv, tcpRecv)
//	//
//	//	for {
//	//		select {
//	//		case <-quitCh:
//	//			return
//	//		default:
//	//		}
//	//
//	//		data, _, err := handle.ReadPacketData()
//	//
//	//		if err == pcap.NextErrorTimeoutExpired {
//	//			break
//	//		} else if err == io.EOF {
//	//			break
//	//		} else if err != nil {
//	//			fmt.Printf("Packet read error: %s\n", err)
//	//			continue
//	//		}
//	//
//	//		decoded := []gopacket.LayerType{}
//	//		if err := parser.DecodeLayers(data, &decoded); err != nil {
//	//			continue
//	//		}
//	//
//	//		for _, layerType := range decoded {
//	//			if layerType == layers.LayerTypeIPv4 {
//	//				if ip4.NetworkFlow() != ipFlow {
//	//					continue
//	//				}
//	//			}
//	//			if layerType != layers.LayerTypeTCP {
//	//				if tcpRecv.DstPort == layers.TCPPort(rawPort) {
//	//					if tcpRecv.SYN && tcpRecv.ACK {
//	//						ch <- true
//	//						showOpenPort = append(showOpenPort, int(tcpRecv.SrcPort))
//	//						tcp := layers.TCP{
//	//							SrcPort: tcpRecv.DstPort,
//	//							DstPort: tcpRecv.SrcPort,
//	//							RST:     true,
//	//							ACK:     true,
//	//							Window:  0,
//	//						}
//	//						tcp.SetNetworkLayerForChecksum(&ip4)
//	//						buf := gopacket.NewSerializeBuffer()
//	//						if err := gopacket.SerializeLayers(buf, opt, &eth, &ip4, &tcp, gopacket.Payload([]byte{})); err != nil {
//	//							log.Fatal(err)
//	//						}
//	//						handle.WritePacketData(buf.Bytes())
//	//						ch <- false
//	//					} else if tcpRecv.RST && tcpRecv.ACK {
//	//						continue
//	//						// ch<-true
//	//						// showBlockPort = append(showBlockPort, int(tcpRecv.SrcPort))
//	//						// ch<-false
//	//					}
//	//				}
//	//			}
//	//		}
//	//	}
//	//}()
//	//
//	//go func() {
//	//	for _, value := range openPort {
//	//		wg.Add(1)
//	//		go func(v []int) {
//	//			defer wg.Done()
//	//			for _, num := range v {
//	//				tcp.DstPort = layers.TCPPort(num)
//	//				buf := gopacket.NewSerializeBuffer()
//	//				// gopacket.SerializeLayers(buf, opt, &eth, &ip4, &tcp, gopacket.Payload([]byte{1, 2, 3, 4})
//	//				if err := gopacket.SerializeLayers(buf, opt, &eth, &ip4, &tcp, gopacket.Payload([]byte{})); err != nil {
//	//					log.Fatal(err)
//	//				}
//	//				handle.WritePacketData(buf.Bytes())
//	//			}
//	//		}(value)
//	//		time.Sleep(10 * time.Millisecond)
//	//	}
//	//	wg.Wait()
//	//}()
//	//
//	//t := time.NewTicker(8 * time.Second)
//	//exitNow := true
//	//
//	//for {
//	//	select {
//	//	case <-t.C:
//	//		exitNow = true
//	//	case changeTicker := <-ch:
//	//		if changeTicker {
//	//			t.Stop()
//	//		} else {
//	//			t = time.NewTicker(1 * time.Second)
//	//		}
//	//		exitNow = false
//	//	}
//	//	if exitNow {
//	//		break
//	//	}
//	//}
//	//
//	//fmt.Println("开放的端口: ", dealRepeat(showOpenPort))
//	//quitCh <- true
//	//return
//}
//
//func init() {
//	runtime.GOMAXPROCS(runtime.NumCPU())
//	flag.Parse()
//	log.SetFlags(log.Lshortfile)
//}
