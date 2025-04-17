package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
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
	//pyroscope.Start(pyroscope.Config{
	//	ApplicationName: "master.app",
	//
	//	// replace this with the address of pyroscope server
	//	ServerAddress:   "http://127.0.0.1:4040",
	//
	//	// you can disable logging by setting this to nil
	//	Logger:          pyroscope.StandardLogger,
	//
	//	// optionally, if authentication is enabled, specify the API key:
	//	// AuthToken: os.Getenv("PYROSCOPE_AUTH_TOKEN"),
	//
	//	// by default all profilers are enabled,
	//	// but you can select the ones you want to use:
	//	ProfileTypes: []pyroscope.ProfileType{
	//		pyroscope.ProfileCPU,
	//		pyroscope.ProfileAllocObjects,
	//		pyroscope.ProfileAllocSpace,
	//		pyroscope.ProfileInuseObjects,
	//		pyroscope.ProfileInuseSpace,
	//	},
	//})

	defer func() {
		if err := recover(); nil != err {
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
		//持续争夺还是只是一次
		pkg.WithSeizeTag(true),

		//过多久检测一次是否可以获得master，如果是master 就跳过，
		//前提是 WithSeizeTag is true
		pkg.WithSeizeCycle(2*time.Second),

		//争夺的的标识 存如hash 中
		// key member:group:
		// WithLockKey 就是每个争夺成员的自己的标识
		// val 是master 或者slave ,可以获取到所有的争夺的成员，
		// 一般设置为ip 可以通过提供的默认方法获取当前ip，如果同一个程序就要设置为自己可以识别的标识
		// 例如 member:group: key slave
		pkg.WithLockKey("key"),

		//redis的操作超时时间,默认3s
		pkg.WithRedisTimeout(3*time.Second),

		//master的超时时间
		pkg.WithExpireTime(5),

		//如果抢到master，续期多长时间,默认expire的一半
		pkg.WithRenewalOften(pkg.DefaultRenewalTime),

		pkg.WithRedisClient(redis),

		// 此标识是表示 此包在项目中可以使用多次，因为一个项目可能存在多个不同组的争夺master操作
		// 例如 第一组 a 、b、c 争夺 master1 谁拿到谁就是master
		// 但是项目中还需要使用 分布式锁  n、 f、 j 争夺 master-1 谁拿到 谁就是主  ，分为不同组的
		// master key争夺
		//master标识，一个项目可能需要多个不同master的抢锁操作，有默认值 member:master
		pkg.WithMasterKey("master"),

		//抢锁失败回调函数
		pkg.WithLockFailFunc(func(i ...interface{}) {
			for _, v := range i {
				fmt.Println("传入的参数", v)
			}
		}),

		//抢锁成功回调函数
		pkg.WithLockSuccessFunc(func(i ...interface{}) {
			for _, v := range i {
				fmt.Println("传入的参数成功", v)
			}
		}),
	)

	gLock.Lock(i, i1)

	//定时获取所有竞争的成员 和 谁是master
	for {

		fmt.Println(gLock.GetMembers()) //全部竞争锁的成员

		fmt.Println(gLock.IsMaster()) //是否是master
		fmt.Println(gLock.Error())    //获取失败原因

		time.Sleep(5 * time.Second)
		break
	}

	gLock.UnLock()

}
