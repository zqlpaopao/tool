package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/register-discovery/pkg"
	"github.com/zqlpaopao/tool/register-discovery/pkg/etcd"
	"github.com/zqlpaopao/tool/register-discovery/pkg/redis"
	clientV3 "go.etcd.io/etcd/client/v3"
	"strconv"
	"time"
)

func redisPs() {
	redisP := pkg.NewProxy(
		redis.NewPubSubRedis(
			redis.NewRedis(
				redis.WithAddr("127.0.0.1:6379")),
			redis.WithDebug(true),
			redis.WithRegistererIsLoopPush(true)))

	fmt.Println(redisP.Error())

	redisP.Discovery()

	redisP.Register(&pkg.RegisterInfo{
		RegisterInfo: "{\n\"addr\":\"127.0.0.1:8080\",\n\"time\":123456\n}",
		PushTime:     0,
	})

	time.Sleep(time.Second * 20)

}

func redisHash() {

	redisP := pkg.NewProxy(
		redis.NewPubHashRedis(
			redis.NewRedis(
				redis.WithAddr("127.0.0.1:6379")),
			redis.WithDebug(true),
			redis.WithRegistererIsLoopPush(true)))

	fmt.Println(redisP.Error())

	redisP.Discovery()

	redisP.Register(&pkg.RegisterInfo{
		RegisterInfo: "127.0.0.1:8080",
		PushTime:     123456,
	})

	time.Sleep(time.Second * 20)
}

func main() {
	//redisPs()
	//redisHash()
	configByEtcd()
	select {}
}

type Cl *clientV3.Client
type CC clientV3.Config

func configByEtcd() {

	etcdCli, err := etcd.NewEtcd[CC, Cl](CC{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}, func(config CC) (Cl, error) {
		return clientV3.New(clientV3.Config(config))
	})
	if err != nil {
		panic(err)
	}

	p := pkg.NewProxy(etcd.NewPDoEtcd(etcdCli))

	p.Discovery()
	fmt.Println("start")
	time.Sleep(2 * time.Second)
	p.UnDiscovery()
	//
	for i := 0; i < 100; i++ {
		p.Register(&pkg.RegisterInfo{
			RegisterInfo: "127.0.0.1:8002" + strconv.Itoa(i)})
	}

}
