package pkg

import (
	"math"
	"sync"
	"sync/atomic"
	"time"
)

/////////////////////////////////////////ping man option ///////////////////////////////

type OptionInter interface {
	apply(*Option)
}

type Option struct {
	OnFinish        Finish
	saveRecover     SavePanic
	source          string
	network         string
	readWait        time.Duration
	connReadTimeout time.Duration
	readyPingN      int
	customerN       int
	customerCS      int
	readyPingCS     int
	initPingMapSiz  int
	debug           bool
}

// Finish The function that calls back after each execution is completed
type Finish func(*Ping)

type SavePanic func()

type OPFunc func(*Option)

func (o OPFunc) apply(opt *Option) {
	o(opt)
}

func NewPingManWithOptions(f ...OPFunc) *PingMan {
	opt := clone(f...)
	return &PingMan{
		opt:               clone(f...),
		id:                GetPId(),
		readyQueue:        make(chan *Ping, opt.readyPingCS),
		resultQueue:       make(chan *ResPing, opt.customerCS),
		initItems:         make(map[string]*Ping, opt.initPingMapSiz),
		initAddrMappingIp: make(map[string]string, opt.initPingMapSiz),
		count:             math.MaxInt64,
		wait:              make(chan struct{}),
		close:             make(chan struct{}),
		lock:              &sync.Mutex{},
		wg:                &sync.WaitGroup{},
		isNotice:          atomic.Bool{},
	}
}

func clone(f ...OPFunc) *Option {
	o := NewDefaultOption()
	for _, v := range f {
		v(o)
	}
	return o
}

func NewDefaultOption() *Option {
	return &Option{
		readWait:        ReadTimeout,
		connReadTimeout: connReadTimeout,
		initPingMapSiz:  InitPingMap,
		source:          "",
		network:         Ip.String(),
		debug:           true,
		readyPingN:      ReadyPingNum,
		customerN:       CustomerNum,
		customerCS:      CustomerChanSize,
		readyPingCS:     ReadyPingChanSize,
		saveRecover:     SaveRecover,
	}
}

// WithReadWait Set the read time for finish
func (o *Option) WithReadWait(duration time.Duration) OPFunc {
	return func(o *Option) {
		o.readWait = duration
	}
}

// WithConnReadTimeout Set the read time for ipv4 and ipv6
func (o *Option) WithConnReadTimeout(duration time.Duration) OPFunc {
	return func(o *Option) {
		o.readWait = duration
	}
}

// WithInitPingMapSiz Set the number of pings to initialize storage
func (o *Option) WithInitPingMapSiz(si int) OPFunc {
	return func(o *Option) {
		o.initPingMapSiz = si
	}
}

// WithSource Set isp
func (o *Option) WithSource(isp string) OPFunc {
	return func(o *Option) {
		o.source = isp
	}
}

// WithDebug Print and record all error messages
func (o *Option) WithDebug(debug bool) OPFunc {
	return func(o *Option) {
		o.debug = debug
	}
}

// WithReadyPingN Goroutine for ping
func (o *Option) WithReadyPingN(pingN int) OPFunc {
	return func(o *Option) {
		o.readyPingN = pingN
	}
}

// WithReadyPingCS chan size for ping
func (o *Option) WithReadyPingCS(pingCSN int) OPFunc {
	return func(o *Option) {
		o.readyPingCS = pingCSN
	}
}

// WithCustomerN Set the number of goroutines that receive message processing from the network card
func (o *Option) WithCustomerN(CsN int) OPFunc {
	return func(o *Option) {
		o.customerN = CsN
	}
}

// WithCustomerCS Set the length of chan for receiving message processing from the network card
func (o *Option) WithCustomerCS(csN int) OPFunc {
	return func(o *Option) {
		o.customerCS = csN
	}
}

// WithSaveRecover Set the recover for the panic
func (o *Option) WithSaveRecover(f func()) OPFunc {
	return func(o *Option) {
		o.saveRecover = f
	}
}
