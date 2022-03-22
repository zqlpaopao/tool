package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pyroscope-io/client/pyroscope"
	"github.com/zqlpaopao/tool/glock/pkg"
	"runtime/debug"
	"time"
)

/*
1、加的锁是自己的，过期时间
2、续期操作 过期时间的一半时间去加，每次加一半，直到，锁被释放(建议默认)，也可以指定加锁续期周期和时间
3、释放的是自己的锁
4、加锁成功回调函数
5、加锁失败回调函数
6、释放锁失败重试，默认三次
7、支持查看当前竞争者
8、支持短锁
*/

func main() {
	pyroscope.Start(pyroscope.Config{
		ApplicationName: "master.app",

		// replace this with the address of pyroscope server
		ServerAddress:   "http://127.0.0.1:4040",

		// you can disable logging by setting this to nil
		Logger:          pyroscope.StandardLogger,

		// optionally, if authentication is enabled, specify the API key:
		// AuthToken: os.Getenv("PYROSCOPE_AUTH_TOKEN"),

		// by default all profilers are enabled,
		// but you can select the ones you want to use:
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
		},
	})

	defer func() {
		if err := recover();nil != err{
			fmt.Println(err)
			fmt.Println(string(debug.Stack()))
		}
	}()

	redis := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	var i int = 1
	var i1 string = "gg"
	gLock := pkg.NewGlock(
		pkg.WithSeizeTag(true),                       //持续争夺还是只是一次
		pkg.WithSeizeCycle(2*time.Second),            //持续争夺还是只是一次
		pkg.WithLockKey("key"),                       //争多的标识
		pkg.WithRedisTimeout(3*time.Second),          //redis的操作超时时间,默认3s
		pkg.WithExpireTime(5),                        //master的超时时间
		pkg.WithRenewalOften(pkg.DefaultRenewalTime), //如果抢到master，续期多长时间,默认expire的一半
		pkg.WithRedisClient(redis),
		pkg.WithLockFailFunc(func(i ...interface{}) { //抢锁失败回调函数
			for _, v := range i {
				fmt.Println("传入的参数", v)
			}
		}),
		pkg.WithLockSuccessFunc(func(i ...interface{}) { //抢锁成功回调函数
			for _, v := range i {
				fmt.Println("传入的参数成功", v)
			}
		}),
	)


	for {
		gLock.Lock(i, i1)

		fmt.Println(gLock.GetMembers())

		fmt.Println(gLock.IsMaster())
		fmt.Println(gLock.Error())
		gLock.UnLock()

		time.Sleep(5 * time.Second)
	}


}
