package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/zqlpaopao/tool/register-discovery/pkg"
	"log"
	"sync/atomic"
	"time"
)

type PubHashRedis struct {
	err    error
	redis  *redis.Client
	option *Options
	close  atomic.Bool
}

func (p *PubHashRedis) Register(msg *pkg.RegisterInfo) {
	if p.err != nil {
		return
	}
	p.Push(msg.RegisterInfo, msg.PushTime)
	if !p.option.Registerer.IsLoopPush {
		return
	}
	go p.LoopPush(msg.RegisterInfo, msg.PushTime)
}

func (p *PubHashRedis) LoopPush(addr string, val int) {
	for {
		time.Sleep(p.option.PushTime)
		p.Push(addr, val)
		if p.close.Load() {
			return
		}
	}
}

func (p *PubHashRedis) Push(addr string, val int) {
	err := p.redis.HSet(context.Background(), p.option.Registerer.Addr, addr, val).Err()
	p.Debug("HSet-->"+p.option.Registerer.Addr+"---->"+addr, p.err)
	if err != nil {
		p.option.CallBackErr(err)
	}
	if !p.option.Registerer.IsLoopPush {
		return
	}
}

func (p *PubHashRedis) UnRegister() {
	return
}

func (p *PubHashRedis) Discovery() {
	go p.LoopDiscovery()
}

func (p *PubHashRedis) LoopDiscovery() {

	pubSub := p.redis.Subscribe(context.Background(), p.option.Registerer.Addr)
	// 处理订阅接收到的消息
	for {

		msg, err := p.redis.HGetAll(context.Background(), p.option.Registerer.Addr).Result()
		if err != nil {
			p.option.CallBackErr(err)
			continue
		}
		p.option.CallBackHash(msg)
		time.Sleep(p.option.PullTime)
		if p.close.Load() {
			goto END
		}
	}
END:
	p.err = pubSub.Close()
}

func (p *PubHashRedis) UnDiscovery() {
	p.close.Store(true)
}

func (p *PubHashRedis) Error() error {
	return p.err
}

func NewPubHashRedis(redis *redis.Client, opt ...OptionFunc) pkg.Abstract {
	o := &PubHashRedis{redis: redis, option: NewOptions(opt...), close: atomic.Bool{}}
	o.Ping()
	return o
}

func (p *PubHashRedis) Ping() {
	p.err = p.redis.Ping(context.Background()).Err()
}

func (p *PubHashRedis) Debug(info string, err error) {
	if !p.option.Debug {
		return
	}
	if err == nil {
		log.Println("DEBUG " + info)
		return
	}
	log.Println("DEBUG "+info, err.Error())
}
