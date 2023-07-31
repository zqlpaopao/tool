package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/batch_ping/pkg"
	"time"
)

func main() {
	p := pkg.NewPingPoolWithOptions()
	p.Run()
	err := p.Error()
	if err != nil {
		panic(err)
	}

	//err = p.Submit("baidu.com")

	for {
		item := &pkg.PingItem{
			Interval:    100 * time.Millisecond,
			Timeout:     3 * time.Second,
			Count:       10,
			PacketsSent: 0,
			PacketsRev:  0,
			OnFinish:    nil,
			Size:        80,
			Tracker:     1897654345,
			Source:      "",
			Addr:        "jd.com",
		}
		err = item.PingWith()
		if err != nil {
			fmt.Println(err)
		}
		err = p.SubmitPingItem(item)
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(1 * time.Second)
	}

	//time.Sleep(30 * time.Second)
	//p.Close()
	//p.CloseRev()
	//
	//p.Wait()

	fmt.Println("end")

	select {}
}
