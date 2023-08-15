package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/batch_ping_agent/pkg"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	go Loop()
	http.ListenAndServe("0.0.0.0:6061", nil)

}

func Loop() {

	p := ping.NewPoolWithOptions(
		//ping.WithProtocol("icmp"),
		ping.WithProtocol("icmp"),
		ping.WithOnRevFunc(func(pings *ping.Ping, packet *ping.Packet) {
			//fmt.Println()
			//fmt.Printf("%#v\n", packet)
			st := ping.StatisticsLog(pings, packet)
			ping.Debug(st)
			//fmt.Println("uuid1--pings", pings.Uuid().String())
			//fmt.Println("uuid1--packs", packet.Uuid())

			//回收资源

		}),
		ping.WithErrCallBack(func(ping *ping.Ping, err error) {
			if err == nil {
				return
			}

			if ping == nil {
				return
			}
			fmt.Printf("%#v", ping)
			fmt.Println(ping.Addr(), err)
			os.Exit(1)
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
		
	}

	for {

		for k, v := range mapAddr {
			IPV4Addr(p, k, v)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func IPV4Addr(p *ping.Pool, addr string, addrNet *net.IPAddr) {
	var err error
	bd := p.GetPing()
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
