package pkg

import "time"

// Config 连接池相关配置
type Config[T any] struct {
	//连接池中拥有的最小连接数
	InitialCap int64
	//最大并发存活连接数
	MaxCap int64
	//最大空闲连接
	MaxIdle int
	//生成连接的方法
	Factory func() (T, error)
	//关闭连接的方法
	Close func(T)
	//检查连接是否有效的方法
	Ping func(T) error
	//连接最大空闲时间，超过该事件则将失效
	IdleTimeout time.Duration
}

// NewPoolWithConfig make new option poll
func NewPoolWithConfig[T any](
	poolSize int64,
	config *Config[T],
) (pool *Pool[T], err error) {

	if err = checkConfig[T](poolSize, config); nil != err {
		return
	}

	return &Pool[T]{
		Conn: make(chan *IdleConn[T], config.MaxCap),
		opt:  config}, nil
}

// checkConfig Verify necessary parameters and set default values
func checkConfig[T any](poolSize int64, c *Config[T]) error {
	if c.InitialCap < 1 {
		c.InitialCap = DefaultInitPoolSize
	}
	if poolSize < 1 {
		c.MaxCap = DefaultPoolSize
	} else {
		c.MaxCap = poolSize
	}
	if c.IdleTimeout < 1 {
		c.IdleTimeout = DefaultIdleTimeout
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
