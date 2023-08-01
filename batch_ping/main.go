package main

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/zqlpaopao/tool/batch_ping/pkg"
	"os"
	"sync"
	"time"
)

type data struct {
	lock *sync.RWMutex
	info *pkg.PingItem
}

var Info = cmap.New[*data]()

func main() {
	p := pkg.NewPingPoolWithOptions(
		pkg.WithOnRevFunc(func(ping *pkg.PingItem, packet *pkg.Packet) {
			fmt.Println("ID", ping.ID())
			fmt.Println("Seq", packet.Seq)
			var (
				p  *data
				ok bool
			)
			if p, ok = Info.Get(ping.Addr); !ok {
				return
			}
			p.lock.Lock()
			if p.info.PacketsRev == p.info.Count {
				p.info.PacketsSent = p.info.Count
				ps := p.info.Statistics()
				pkg.Log(ps)
				p.info.ResetRttS()
				//
			} else {
				p.info.PacketsRev++
				p.info.RttS = append(p.info.RttS, packet.EndTime.Sub(packet.StartTime))
			}
			//
			p.lock.Unlock()

		}),
		pkg.WithErrCallBack(func(ping *pkg.PingItem, err error) {
			if err == nil {
				return
			}
			fmt.Println(ping.Addr, err)
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
		//pushAddr(p, "baidu.com")
		//pushAddr(p, "qq.com")

		i++
	}
	p.Close()
	p.CloseRev()

	//for {
	//
	//	time.Sleep(1 * time.Second)
	//}

	//time.Sleep(30 * time.Second)
	//p.Close()
	//p.CloseRev()
	//
	//p.Wait()

}

func pushPingItem(p *pkg.PingPool, addr string) {
	var err error
	item := &pkg.PingItem{
		Interval:    100 * time.Millisecond,
		Timeout:     3 * time.Second,
		Count:       3,
		PacketsSent: 0,
		PacketsRev:  0,
		Size:        800,
		Tracker:     1897654345,
		Source:      "",
		Addr:        addr,
	}

	if _, ok := Info.Get(item.Addr); !ok {
		Info.Set(item.Addr, &data{
			lock: &sync.RWMutex{},
			info: item,
		})
	}

	err = item.PingWith()

	if err != nil {
		fmt.Println(err)
	}
	err = p.SubmitPingItem(item)
	if err != nil {
		fmt.Println(err)
	}

}

func pushAddr(p *pkg.PingPool, addr string) {
	var (
		err error
		ps  *pkg.PingItem
	)

	if ps, err = p.Submit(addr); nil != err {
		fmt.Println(err)
		os.Exit(1)
	}
	if _, ok := Info.Get(addr); !ok {
		Info.Set(addr, &data{
			lock: &sync.RWMutex{},
			info: ps,
		})
	}
}
