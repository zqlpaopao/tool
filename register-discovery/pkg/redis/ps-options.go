package redis

import (
	"github.com/zqlpaopao/tool/register-discovery/pkg/common"
	"time"
)

type (
	Options struct {
		CallBackErr    func(error)
		CallBackPubSub func(*redis.Message)
		CallBackHash   func(map[string]string)
		Registerer
		Discovery
		Debug bool
	}
	Registerer struct {
		Addr string
		//上报时间
		PushTime   time.Duration
		IsLoopPush bool
	}

	Discovery struct {
		Addr []string
		//发现时间
		PullTime time.Duration
	}
)

type Option interface {
	apply(opt *Options)
}

type OptionFunc func(opts *Options)

func (o OptionFunc) apply(opt *Options) {
	o(opt)
}

// NewOptions creates a redis cli
// This can mask the differences between different versions of Redis
func NewOptions(opts ...OptionFunc) *Options {
	redisOpts := DefaultOptions()
	for _, opt := range opts {
		opt.apply(redisOpts)
	}
	return redisOpts
}

func DefaultOptions() *Options {
	return &Options{
		Registerer: Registerer{
			Addr:       common.RedisDefaultRegisterPub,
			PushTime:   3 * time.Second,
			IsLoopPush: false,
		},
		Discovery: Discovery{
			Addr:     []string{common.RedisDefaultRegisterPub},
			PullTime: 2 * time.Second,
		},
		CallBackErr:    CallBackErr,
		CallBackPubSub: CallBack,
		CallBackHash:   CallBackHash,
	}
}

func WithRegistererAddr(addr string) OptionFunc {
	return func(opts *Options) {
		opts.Registerer.Addr = addr
	}
}

func WithDiscoveryAddr(addr []string) OptionFunc {
	return func(opts *Options) {
		opts.Discovery.Addr = addr
	}
}

func WithRegistererPushTime(time time.Duration) OptionFunc {
	return func(opts *Options) {
		opts.Registerer.PushTime = time
	}
}
func WithRegistererIsLoopPush(loop bool) OptionFunc {
	return func(opts *Options) {
		opts.Registerer.IsLoopPush = loop
	}
}

func WithRCallBackErr(f func(err error)) OptionFunc {
	return func(opts *Options) {
		opts.CallBackErr = f
	}
}

func WithRCallBack(f func(message *redis.Message)) OptionFunc {
	return func(opts *Options) {
		opts.CallBackPubSub = f
	}
}

func WithRCallBackHash(f func(message map[string]string)) OptionFunc {
	return func(opts *Options) {
		opts.CallBackHash = f
	}
}

func WithDebug(debug bool) OptionFunc {
	return func(opts *Options) {
		opts.Debug = debug
	}
}
