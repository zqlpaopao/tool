package pkg

import "time"

// Config 连接池相关配置
type Config[T any] struct {
	Factory       func() (T, error)
	Close         func(T)
	Ping          func(T) error
	InitialCap    int64
	MaxCap        int64
	MaxIdle       int
	IdleTimeout   time.Duration
	CheckInterval time.Duration
	IsCheck       bool
	Debug         bool
}

// NewPoolWithConfig make new option poll
func NewPoolWithConfig[T any](
	config *Config[T]) (
	pool *Pool[T]) {
	p := &Pool[T]{close: DefaultIsRunning}

	if p.err = checkConfig[T](config); nil != p.err {
		return
	}
	p.Conn,
		p.opt =
		make(chan *IdleConn[T], config.MaxCap),
		config
	return p
}

// checkConfig Verify necessary parameters and set default values
func checkConfig[T any](c *Config[T]) error {
	if c.InitialCap < 1 {
		c.InitialCap = DefaultInitPoolSize
	}
	if c.MaxCap < 1 {
		c.MaxCap = DefaultPoolSize
	}
	if c.MaxCap < c.InitialCap {
		c.MaxCap = c.InitialCap
	}
	if c.IdleTimeout < 1 {
		c.IdleTimeout = DefaultIdleTimeout
	}
	if c.IsCheck && c.CheckInterval < 1 {
		c.CheckInterval = DefaultCheckInterval
	}

	if c.Factory == nil {
		return ErrFactoryIssNil
	}

	if c.Close == nil {
		return ErrCloseIssNil
	}

	if c.Ping == nil {
		return ErrPingIssNil
	}
	return nil
}
