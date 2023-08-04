package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/batch_ping_tool/pkg"
	"sync/atomic"
	"time"
)

func main() {
	p := ping.NewPoolWithOptions(
		ping.WithProtocol("icmp"),
		ping.WithOnRevFunc(func(pings *ping.Ping, packet *ping.Packet) {
			pings.AppendRtt(packet.EndTime.Sub(packet.StartTime))

			if atomic.AddInt64(&pings.PacketsReceive, 1) == pings.Count {
				ps := ping.StatisticsLog(pings)
				ping.Debug(ps)
				pings.Reset()
			}
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
		}),
	).Run()
	err := p.Error()
	if err != nil {
		panic(err)
	}
	var i = 0

	for i < 10 {
		pushPingItem(p, "baidu.com")
		pushPingItem(p, "qq.com")

		i++
	}
	p.Close()
	p.CloseRev()
	//select {}

}

func pushPingItem(p *ping.Pool, addr string) {
	ps, err := ping.NewPing(addr)
	if err != nil {
		fmt.Println(err)
	}
	ps.Interval = 100 * time.Millisecond
	ps.Timeout = 3 * time.Second
	ps.Count = 3

	if err != nil {
		fmt.Println(err)
	}
	err = p.Submit(ps)
	if err != nil {
		fmt.Println(err)
	}

}
