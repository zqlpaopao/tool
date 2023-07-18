package etcd

import (
	"github.com/zqlpaopao/tool/register-discovery/pkg/common"
	etcd "go.etcd.io/etcd/client/v3"
	"time"
)

type (
	Options struct {
		Registerer
		Discovery
		Debug         bool
		CallBackErr   func(string, error)
		CallBackWatch func(*etcd.Event)
		CallBack      func(*etcd.GetResponse)
	}
	Registerer struct {
		Addr string
		//上报时间
		PushTime int64
	}

	Discovery struct {
		Addr string
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
			Addr:     common.EtcdDefaultRegisterSub,
			PushTime: 3,
		},
		Discovery: Discovery{
			Addr:     common.EtcdDefaultRegisterSub,
			PullTime: 2 * time.Second,
		},
		CallBackErr:   CallBackErr,
		CallBackWatch: CallBackWatch,
		CallBack:      CallBack,
	}
}

func WithRegistererAddr(addr string) OptionFunc {
	return func(opts *Options) {
		opts.Registerer.Addr = addr
	}
}

func WithDiscoveryAddr(addr string) OptionFunc {
	return func(opts *Options) {
		opts.Discovery.Addr = addr
	}
}

func WithRegistererLeaseTime(ttl int64) OptionFunc {
	return func(opts *Options) {
		opts.Registerer.PushTime = ttl
	}
}

func WithRCallBackErr(f func(funcName string, err error)) OptionFunc {
	return func(opts *Options) {
		opts.CallBackErr = f
	}
}

func WithRCallBackWatch(f func(message *etcd.Event)) OptionFunc {
	return func(opts *Options) {
		opts.CallBackWatch = f
	}
}

func WithRCallBack(f func(*etcd.GetResponse)) OptionFunc {
	return func(opts *Options) {
		opts.CallBack = f
	}
}

func WithDebug(debug bool) OptionFunc {
	return func(opts *Options) {
		opts.Debug = debug
	}
}
