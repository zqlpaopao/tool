package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/ip/src"
	"github.com/zqlpaopao/tool/raft_lock/pkg"
	"strconv"
	"time"
)

func main() {
	p := pkg.NewProxy(pkg.NewRaft(pkg.NewRedisOption(
		pkg.WithAddr([]string{
			"127.0.0.1:6380",
			"127.0.0.1:6381",
			"127.0.0.1:6382",
		}),
		pkg.WithLockNum(3),
		pkg.WithNodeNum(3),
	)))

	p.Init()

	fmt.Println(p.GetLockInfo())

	for i := 0; i < 10; i++ {
		go makeRaft(p, src.GetEth0()+"_"+strconv.Itoa(i), time.Second*30)
	}

	select {}

}

func makeRaft(
	p *pkg.Proxy,
	lockName string,
	timeout time.Duration,
) {

	var (
		closeTag = make(chan struct{})
	)

	for {
		res, _ := p.Lock(lockName, timeout)
		if res != 1 {
			time.Sleep(timeout - 10*time.Second)
			continue
		}
		fmt.Println(lockName, "加锁成功", res, time.Now().Format("2006-01-02 15:04:05"))
		go func() {
			var (
				timer = time.NewTicker(timeout - 10*time.Second)
			)
			for {
				select {
				case <-closeTag:
					timer.Stop()
					return
				case <-timer.C:
					res, err := p.Renewal(lockName, timeout)
					fmt.Println(1, res, err)
					if res == 1 {
						fmt.Printf("----------->%v 续期成功 %v\n", lockName, time.Now().Format("2006-01-02 15:04:05"))
					} else {
						closeTag <- struct{}{}
						fmt.Printf("%v 续期失败，error:%v\n", lockName, err)
					}
					timer.Reset(timeout - 10*time.Second)
				}
			}
		}()

		//处理数据
		go func() {
			var (
				timer = time.NewTicker(time.Second * 3)
			)
			for {
				select {
				case <-closeTag:
					timer.Stop()
					return
				case <-timer.C:
					fmt.Println(lockName, "处理任务")
					timer.Reset(time.Second * 3)
				}
			}
		}()

	}

}
