package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/zqlpaopao/tool/batch_ping_agent/pkg"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "net/http/pprof"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	var pool = ping.NewPoolData(10)

	p := ping.NewPoolWithOptions(
		//ping.WithProtocol("icmp"),
		ping.WithProtocol("udp"),
		ping.WithOnRevFunc(func(pings *ping.Ping, packet *ping.Packet) {
			fmt.Println()
			fmt.Printf("%#v\n", packet)
			st := ping.StatisticsLog(pings, packet)
			ping.Debug(st)
			pool.SetPing(pings)
			//os.Exit(1)

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
	for {
		IPV4Addr(pool, p, "110.242.68.66", 110, 242, 68, 66)
		IPV4Addr(pool, p, "123.151.137.18", 123, 151, 137, 18)
		IPV4Addr(pool, p, "203.205.254.157", 203, 205, 254, 157)

		time.Sleep(100 * time.Millisecond)
	}
	//p.Close()
	//p.CloseRev()
	http.ListenAndServe("0.0.0.0:6061", nil)

}

func IPV4Addr(pool *ping.PoolData, p *ping.Pool, addr string, a, b, c, d byte) {
	var err error
	bd := pool.GetPing()
	u := uuid.New()
	fmt.Println("------->uuid", u)
	bd.SetDstAddr(addr).SetIpV4().SetResolveIpAddr(&net.IPAddr{IP: net.IPv4(a, b, c, d), Zone: ""}).SetSize(60).SetTtl(60).SetUUid(u)

	if err = p.Submit(bd); nil != err {
		panic(err)
	}

}

func hostIpV4(pool *ping.PoolData, p *ping.Pool) {
	var err error
	qq := pool.GetPing()
	if err = qq.SetAddr("qq.com"); nil != err {
		panic(err)
	}
	qq.SetSize(40).SetTtl(60).SetUUid(uuid.New())

	if err = p.Submit(qq); nil != err {
		panic(err)
	}

}
