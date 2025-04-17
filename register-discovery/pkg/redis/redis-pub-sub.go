package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	format "github.com/zqlpaopao/tool/format/src"
	"github.com/zqlpaopao/tool/register-discovery/pkg"
	"log"
	"sync/atomic"
	"time"
)

type PubSubRedis struct {
	err    error
	redis  *redis.Client
	option *Options
	close  atomic.Bool
}

func (p *PubSubRedis) Register(msg *pkg.RegisterInfo) {
	if p.err != nil {
		return
	}
	p.Push(msg.RegisterInfo)
	if !p.option.Registerer.IsLoopPush {
		return
	}
	go p.LoopPush(msg.RegisterInfo)
}

func (p *PubSubRedis) LoopPush(addr string) {
	for {
		time.Sleep(p.option.PushTime)
		p.Push(addr)
		if p.close.Load() {
			return
		}
	}
}

func (p *PubSubRedis) Push(addr string) {
	err := p.redis.Publish(context.Background(), p.option.Registerer.Addr, addr).Err()
	p.Debug("Publish-->"+p.option.Registerer.Addr+"---->"+addr, p.err)
	if err != nil {
		p.option.CallBackErr(err)
	}
}

func (p *PubSubRedis) UnRegister() {
	return
}

func (p *PubSubRedis) Discovery() {
	go p.LoopDiscovery()
}

func (p *PubSubRedis) LoopDiscovery() {

	pubSub := p.redis.Subscribe(context.Background(), p.option.Registerer.Addr)
	// 处理订阅接收到的消息
	for {
		msg, err := pubSub.ReceiveMessage(context.Background())
		if err != nil {
			p.option.CallBackErr(err)
			continue
		}
		p.option.CallBackPubSub(msg)
		time.Sleep(p.option.PullTime)
		if p.close.Load() {
			goto END
		}
	}
END:
	p.err = pubSub.Close()
}

func (p *PubSubRedis) UnDiscovery() {
	p.close.Store(true)
}

func (p *PubSubRedis) Error() error {
	return p.err
}

func NewPubSubRedis(redis *redis.Client, opt ...OptionFunc) pkg.Abstract {
	o := &PubSubRedis{redis: redis, option: NewOptions(opt...), close: atomic.Bool{}}
	o.Ping()
	return o
}

func (p *PubSubRedis) Ping() {
	p.err = p.redis.Ping(context.Background()).Err()
}

func (p *PubSubRedis) Debug(info string, err error) {
	if !p.option.Debug {
		return
	}
	if err == nil {
		log.Println("DEBUG " + info)
		return
	}
	log.Println("DEBUG "+info, err.Error())
}

func CallBackErr(err error) {
	if err != nil {
		format.PrintRed(err.Error())
	}
}

func CallBack(message *redis.Message) {
	format.PrintGreen(fmt.Sprintf("callback ---> %s--->%s", message.Channel, message.Payload))
}

func CallBackHash(message map[string]string) {
	fmt.Println(message)
}
