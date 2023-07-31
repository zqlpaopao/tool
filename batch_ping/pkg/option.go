package pkg

import (
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type Options interface {
	apply(*Option)
}

type OpFunc func(*Option)

func (o OpFunc) apply(opt *Option) {
	o(opt)
}

type Option struct {
	OnRevFunc       func(*Packet)
	Recover         func()
	errCallBack     func(ping *PingItem, err error)
	network         string
	source          string
	onRevWorkerNum  int
	connReadTimeout time.Duration
	pid             int
	currentMapSize  int
	readyChanSize   int
	customerSize    int
	resCustomerSize int
	readTimeout     time.Duration
	onRevChanSize   int
	ResChanSize     int
	mapTidyInterval time.Duration
	errChanSize     int
	debug           bool
}

func NewPingPoolWithOptions(f ...OpFunc) *PingPool {
	return newPingPool(GetOptionWithOpFunc(f...))
}

func newPingPool(opt *Option) *PingPool {
	return &PingPool{
		option:      opt,
		readyChan:   make(chan *PingItem, opt.readyChanSize),
		errChan:     make(chan *ErrInfo, opt.readyChanSize),
		onRevChan:   make(chan *Packet, opt.onRevChanSize),
		current:     0,
		sendTotal:   0,
		mapTidyLock: &sync.Mutex{},
		wg:          &sync.WaitGroup{},
		revWg:       &sync.WaitGroup{},
		readWg:      &sync.WaitGroup{},
		errWg:       &sync.WaitGroup{},
		done:        make(chan struct{}),
	}

}

func GetOptionWithOpFunc(f ...OpFunc) *Option {
	o := defaultOption()
	for _, fu := range f {
		fu(o)
	}
	return o
}

func defaultOption() *Option {
	return &Option{
		readyChanSize:   ReadyChanSize,
		onRevChanSize:   OnRevChanSize,
		onRevWorkerNum:  OnRevChanSize,
		errChanSize:     ErrChanSize,
		currentMapSize:  CurrentMapSize,
		customerSize:    CustomerSize,
		resCustomerSize: ResCustomerSize,
		connReadTimeout: TimeOut,
		readTimeout:     time.Millisecond * 100,
		pid:             os.Getpid(),
		network:         IP,
		mapTidyInterval: time.Hour,
		errCallBack:     defaultErrCallBack,
		Recover:         defaultRecover,
		OnRevFunc:       OnRevFunc,
	}
}

func OnRevFunc(pkt *Packet) {
	fmt.Println("----------------pkt Packet start---------------------")
	fmt.Printf("%#v\n", pkt)
	fmt.Println("----------------pkt Packet end-----------------------")
}

func defaultErrCallBack(ping *PingItem, err error) {
	if err == nil {
		return
	}
	fmt.Println("----------------ping PingItem start---------------------")
	fmt.Printf("%#v\n", ping)
	fmt.Println(err)
	fmt.Println("----------------ping PingItem end-----------------------")

}

func defaultRecover() {
	if err := recover(); nil != err {
		fmt.Println()
		fmt.Println("----------------Panic start---------------------")
		fmt.Println(string(debug.Stack()))
		fmt.Println(err)
		fmt.Println("----------------Panic end-----------------------")
		fmt.Println()
	}
}
