package pkg

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

//func Test_gLock_Lock(t *testing.T) {
//	redis := redis.NewClient(&redis.Options{
//		Addr: "127.0.0.1:6379",
//	})
//
//	var i int = 1
//	var i1 string = "gg"
//	gLock := NewGlock(
//		WithSeizeTag(true),                   //持续争夺还是只是一次
//		WithSeizeCycle(2*time.Second),        //持续争夺还是只是一次
//		WithLockKey("key666"),                   //争多的标识
//		WithRedisTimeout(3*time.Second),      //redis的操作超时时间,默认3s
//		WithExpireTime(5),                    //master的超时时间
//		WithRenewalOften(DefaultRenewalTime), //如果抢到master，续期多长时间,默认expire的一半
//		WithRedisClient(redis),
//		WithLockFailFunc(func(i ...interface{}) { //抢锁失败回调函数
//			for _, v := range i {
//				fmt.Println("传入的参数", v)
//			}
//		}),
//		WithLockSuccessFunc(func(i ...interface{}) { //抢锁成功回调函数
//			for _, v := range i {
//				fmt.Println("传入的参数成功", v)
//			}
//		}),
//	)
//	gLock.Lock(i, i1)
//
//	fmt.Println(gLock.GetMembers())
//
//	gLock.UnLock()
//	fmt.Println(gLock.IsMaster())
//	fmt.Println(gLock.Error())
//
//}

func BenchmarkGLock_Lock(b *testing.B) {
	redis := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	var i2 int = 1
	var i1 string = "gg"
	gLock := NewGlock(
		WithSeizeTag(true),                   //持续争夺还是只是一次
		WithSeizeCycle(2*time.Second),        //持续争夺还是只是一次
		WithLockKey("key233"),                   //争多的标识
		WithRedisTimeout(3*time.Second),      //redis的操作超时时间,默认3s
		WithExpireTime(2),                    //master的超时时间
		WithRenewalOften(DefaultRenewalTime), //如果抢到master，续期多长时间,默认expire的一半
		WithRedisClient(redis),
		WithLockFailFunc(func(i ...interface{}) { //抢锁失败回调函数
			for _, v := range i {
				fmt.Println("传入的参数", v)
			}
		}),
		WithLockSuccessFunc(func(i ...interface{}) { //抢锁成功回调函数
			for _, v := range i {
				fmt.Println("传入的参数成功", v)
			}
		}),
	)
	for i := 0; i < b.N; i++ {

		gLock.Lock(i2, i1)

		fmt.Println(gLock.GetMembers())

		gLock.UnLock()
		fmt.Println(gLock.IsMaster())
		fmt.Println(gLock.Error())
		time.Sleep(time.Millisecond*500)
	}
}