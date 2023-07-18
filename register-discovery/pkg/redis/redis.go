package redis

import (
	"context"
	"crypto/tls"
	"github.com/go-redis/redis/v8"
	"net"
	"time"
)

// NewRedis creates a redis cli
// This can mask the differences between different versions of Redis
func NewRedis(opts ...OptionFuncRedis) *redis.Client {
	redisOpts := defaultRedisCOnf()
	for _, opt := range opts {
		opt.apply(redisOpts)
	}
	return redis.NewClient(redisOpts)
}

func defaultRedisCOnf() *redis.Options {
	return &redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}
}

type OptionRedis interface {
	apply(opt *redis.Options)
}

type OptionFuncRedis func(opts *redis.Options)

func (o OptionFuncRedis) apply(opt *redis.Options) {
	o(opt)
}

func WithPassword(password string) OptionFuncRedis {
	return func(opts *redis.Options) {
		opts.Password = password
	}
}

func WithDB(db int) OptionFuncRedis {
	return func(opts *redis.Options) {
		opts.DB = db
	}
}

func WithTLSConfig(t *tls.Config) OptionFuncRedis {
	return func(opts *redis.Options) {
		opts.TLSConfig = t
	}
}

func WithDialer(dialer func(ctx context.Context, network, addr string) (net.Conn, error)) OptionFuncRedis {
	return func(opts *redis.Options) {
		opts.Dialer = dialer
	}
}

func WithReadTimeout(t time.Duration) OptionFuncRedis {
	return func(opts *redis.Options) {
		opts.ReadTimeout = t
	}
}

func WithWriteTimeout(t time.Duration) OptionFuncRedis {
	return func(opts *redis.Options) {
		opts.WriteTimeout = t
	}
}
func WithAddr(addr string) OptionFuncRedis {
	return func(opts *redis.Options) {
		opts.Addr = addr
	}
}

func WithPoolSize(num int) OptionFuncRedis {
	return func(opts *redis.Options) {
		opts.PoolSize = num
	}
}
