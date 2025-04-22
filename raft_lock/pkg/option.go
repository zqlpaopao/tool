package pkg

import (
	"fmt"
	"time"
)

type RedisOption struct {
	addr         []string
	password     string
	groupName    string
	lockNum      int
	nodeNum      int
	db           int
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
	poolTimeOut  time.Duration
	poolSize     int
	isCluster    bool
}

type Option interface {
	apply(opt *RedisOption)
}

type OpFunc func(*RedisOption)

// apply assignment function entity
func (o OpFunc) apply(opt *RedisOption) {
	o(opt)
}

// NewRedisOption make RedisOption
func NewRedisOption(opt ...Option) *RedisOption {
	o := &RedisOption{
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		idleTimeout:  idleTimeout,
		poolTimeOut:  poolTimeOut,
		poolSize:     poolSize,
		groupName:    DefaultSameSlot + groupDefaultName,
		lockNum:      lockNum,
		nodeNum:      nodeNum,
		isCluster:    true,
	}
	return o.WithOptions(opt...)
}

// clone  new object
func (o RedisOption) clone() *RedisOption {
	cp := o
	return &cp
}

// WithOptions Execute assignment function entity
func (o RedisOption) WithOptions(opt ...Option) *RedisOption {
	c := o.clone()
	for _, v := range opt {
		v.apply(c)
	}
	return c
}

// WithGroupName Set the number of ReadTimeout. The default is  5s
func WithGroupName(groupName string) OpFunc {
	return func(o *RedisOption) {
		o.groupName = fmt.Sprintf("%v%v", DefaultSameSlot, groupName)
	}
}

// WithLockNum Set the number of ReadTimeout. The default is  5s
func WithLockNum(lockNum int) OpFunc {
	return func(o *RedisOption) {
		o.lockNum = lockNum
	}
}

// WithNodeNum Set the number of ReadTimeout. The default is  5s
func WithNodeNum(nodeNum int) OpFunc {
	return func(o *RedisOption) {
		o.nodeNum = nodeNum
	}
}

// WithIsCluster Set the number of ReadTimeout. The default is  5s
func WithIsCluster(isCluster bool) OpFunc {
	return func(o *RedisOption) {
		o.isCluster = isCluster
	}
}

// WithReadTimeout Set the number of ReadTimeout. The default is  5s
func WithReadTimeout(readTimeout time.Duration) OpFunc {
	return func(o *RedisOption) {
		o.readTimeout = readTimeout
	}
}

// WithAddr Set the number of ReadTimeout. The default is  5s
func WithAddr(addr []string) OpFunc {
	return func(o *RedisOption) {
		o.addr = addr
	}
}

// WithPassword Set the number of ReadTimeout. The default is  5s
func WithPassword(password string) OpFunc {
	return func(o *RedisOption) {
		o.password = password
	}
}

// WithDB Set the number of ReadTimeout. The default is  5s
func WithDB(db int) OpFunc {
	return func(o *RedisOption) {
		o.db = db
	}
}

// WithWriteTimeout Set the WriteTimeout. The default is 5s
func WithWriteTimeout(writeTimeout time.Duration) OpFunc {
	return func(o *RedisOption) {
		o.writeTimeout = writeTimeout
	}
}

// WithIdleTimeout Set the IdleTimeout. The default value is 60s
func WithIdleTimeout(idleTimeout time.Duration) OpFunc {
	return func(o *RedisOption) {
		o.idleTimeout = idleTimeout
	}
}

// WithPoolTimeOut default 60s
func WithPoolTimeOut(poolTimeOut time.Duration) OpFunc {
	return func(o *RedisOption) {
		o.poolTimeOut = poolTimeOut
	}
}

// WithPoolSize set pool size default 20
func WithPoolSize(poolSize int) OpFunc {
	return func(o *RedisOption) {
		o.poolSize = poolSize
	}
}
