package pkg

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"sync"
	"time"
	"unsafe"
)

type OpFunc interface {
	apply(*SshOption)
}

type OPFunc func(*SshOption)

func (o OPFunc) apply(option *SshOption) {
	o(option)
}

type SshOption struct {
	config       *ssh.ClientConfig
	callLog      func(string)
	Addr         string
	network      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	dialTimeout  time.Duration
}

// NewSsh ssh object
func NewSsh(f ...OPFunc) *Ssh {
	return &Ssh{option: clone(f...)}
}

// clone make new sshOption
func clone(f ...OPFunc) *SshOption {
	o := &SshOption{
		readTimeout:  ReadTimeout * time.Second,
		writeTimeout: WriteTimeout * time.Second,
		dialTimeout:  DialTimeout * time.Second,
		network:      Tcp,
		callLog:      defaultCallLog,
	}
	for i := 0; i < len(f); i++ {
		f[i](o)
	}
	return o
}

// defaultCallLog Default connection error callback function
func defaultCallLog(s string) {
	fmt.Println(s)
}

// WithCallLog connection error callback function
func WithCallLog(f func(string)) OPFunc {
	return func(o *SshOption) {
		o.callLog = f
	}
}

// WithAddr must Connect address
func WithAddr(s string) OPFunc {
	return func(o *SshOption) {
		o.Addr = s
	}
}

// WithReadTimeout Read timeout
func WithReadTimeout(t time.Duration) OPFunc {
	return func(o *SshOption) {
		o.readTimeout = t
	}
}

// WithWriteTimeout  Read timeout
func WithWriteTimeout(t time.Duration) OPFunc {
	return func(o *SshOption) {
		o.writeTimeout = t
	}
}

// WithDialTimeout Maximum time for connection establishment
func WithDialTimeout(t time.Duration) OPFunc {
	return func(o *SshOption) {
		o.dialTimeout = t
	}
}

// WithNetwork type of network
func WithNetwork(s string) OPFunc {
	return func(o *SshOption) {
		o.network = s
	}
}

// WithSshConfig Connection Configuration
func WithSshConfig(conf *ssh.ClientConfig) OPFunc {
	return func(o *SshOption) {
		o.config = conf
	}
}

// PoolInter //////////////////////////////////////// Pool option ////////////////////////////////////////
type PoolInter interface {
	apply(*PoolOption)
}

type PoolOPFunc func(*PoolOption)

func (o PoolOPFunc) apply(option *PoolOption) {
	o(option)
}

type PoolOption struct {
	item        Item
	name        string
	maxConnNum  int
	maxIdleNum  int
	maxLifeTime time.Duration
	maxIdleTime time.Duration
}

func NewPool[T any](f ...PoolOPFunc) *Pool[T] {
	return &Pool[T]{cli: make(map[unsafe.Pointer]Item, 10), option: clonePool(f...), lock: &sync.RWMutex{}}
}

func clonePool(f ...PoolOPFunc) *PoolOption {
	o := &PoolOption{
		maxConnNum:  MaxConnNum,
		maxIdleNum:  MaxIdleNum,
		maxLifeTime: MaxLifeTime,
		maxIdleTime: MaxIdleTime,
	}
	for i := 0; i < len(f); i++ {
		f[i](o)
	}
	return o
}

// WithPoolName must Connection pool ID
func WithPoolName(s string) PoolOPFunc {
	return func(o *PoolOption) {
		o.name = s
	}
}

// WithItem must Connect replicated objects
func WithItem(i Item) PoolOPFunc {
	return func(o *PoolOption) {
		o.item = i
	}
}

// WithMaxConnNum max conn
func WithMaxConnNum(i int) PoolOPFunc {
	return func(o *PoolOption) {
		o.maxConnNum = i
	}
}

// WithMaxIdleNum  Connect max idle
func WithMaxIdleNum(i Item) PoolOPFunc {
	return func(o *PoolOption) {
		o.item = i
	}
}

// WithMaxLifeTime  conn life
func WithMaxLifeTime(i Item) PoolOPFunc {
	return func(o *PoolOption) {
		o.item = i
	}
}

// WithMaxIdleTime must Connect idle
func WithMaxIdleTime(i Item) PoolOPFunc {
	return func(o *PoolOption) {
		o.item = i
	}
}
