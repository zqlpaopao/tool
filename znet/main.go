package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/znet/ping"
	"golang.org/x/net/bpf"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sync/atomic"
	"time"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	go Loop()
	http.ListenAndServe("0.0.0.0:6061", nil)

}

var num int32

func Loop() {

	p := ping.NewPoolWithOptions(
		ping.WithPrepareChLen(100),
		ping.WithProtocol("icmp"),
		ping.WithDataSize(60),
		ping.WithPid(6001),
		ping.WithBPFFilter(ping.Filter{
			// Load "EtherType" field from the ethernet header.
			bpf.LoadAbsolute{Off: 24, Size: 2},
			//bpf.JumpIf{Cond: bpf.JumpNotEqual, Val: 6001, SkipFalse: 1},
			bpf.JumpIf{Cond: bpf.JumpEqual, Val: uint32(6001), SkipFalse: 1},
			// Verdict is "send up to 4k of the packet to userspace."
			bpf.RetConstant{Val: 4096},
			// Verdict is "ignore packet."
			bpf.RetConstant{Val: 0},
		}),
		ping.WithOnRevFunc(func(pings *ping.Ping, packet *ping.Packet) {
			//fmt.Println()
			//fmt.Printf("%#v\n", packet)
			//st := ping.StatisticsLog(pings, packet)
			//ping.Debug(st)
			//fmt.Println("uuid1--pings", pings.Uuid().String())
			//fmt.Println("uuid1--packs", packet.Uuid())
			//os.Exit(1)
			//回收资源
			//fmt.Printf("u := pings.uuid %s %p\n", pings.IPAddr(), *pings.Uuid())
			atomic.AddInt32(&num, 1)

		}),
		ping.WithErrorCallback(func(err *ping.ErrInfo) {
			fmt.Println(err.Tag, err.Ping, err.Err)
		}),
	).Run()
	err := p.Error()
	if err != nil {
		panic(err)
	}

	var mapAddr = map[string]*net.IPAddr{

		"110.242.68.66":   {IP: net.IPv4(110, 242, 68, 66)},
		"123.151.137.18":  {IP: net.IPv4(123, 151, 137, 18)},
		"203.205.254.157": {IP: net.IPv4(203, 205, 254, 157)},
		"172.20.114.51":   {IP: net.IPv4(172, 20, 114, 51)},
		"11.97.22.67 ":    {IP: net.IPv4(11, 97, 22, 67)},
		"11.97.22.76 ":    {IP: net.IPv4(11, 97, 22, 76)},
		"11.97.22.115":    {IP: net.IPv4(11, 97, 22, 115)},
		"11.97.22.167":    {IP: net.IPv4(11, 97, 22, 167)},
		"11.97.22.164":    {IP: net.IPv4(11, 97, 22, 164)},
		"11.97.22.117":    {IP: net.IPv4(11, 97, 22, 117)},
	}
	t := time.NewTimer(60 * time.Second)

	for {

		select {
		case <-t.C:
			fmt.Println("end------------------------", atomic.LoadInt32(&num))
			t.Reset(time.Second * 60)

		default:
			for k, v := range mapAddr {
				IPV4Addr(p, k, v)

			}
			time.Sleep(100 * time.Millisecond)
		}

	}
}

// 228085 6001, 6002, 6003, 6004, 6005
// 2874707
func IPV4Addr(p *ping.Pool, addr string, addrNet *net.IPAddr) {
	var err error
	bd := p.GetObjectPingPool()
	u := ping.NewUUid()
	bd.SetDstAddr(addr).SetIpV4().SetResolveIpAddr(addrNet).SetSize(60).SetTtl(60).SetUUid(&u)
	if err = p.Submit(bd); nil != err {
		panic(err)
	}

}

func hostIpV4(p *ping.Pool) {
	//var err error
	//qq := pool.GetPing()
	//if err = qq.SetAddr("qq.com"); nil != err {
	//	panic(err)
	//}
	//qq.SetSize(40).SetTtl(60).SetUUid(uuid.New())
	//
	//if err = p.Submit(qq); nil != err {
	//	panic(err)
	//}

}
