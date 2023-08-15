package ping

import (
	"fmt"
	"math"
	"math/rand"
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
	OnRevFunc   func(*Ping, *Packet)
	Recover     func()
	errCallBack func(ping *Ping, err error)
	// network is one of "ip", "ip4", or "ip6".
	network string
	// protocol is "icmp" or "udp".
	protocol        string
	source          string
	onRevWorkerNum  int
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
	dataSize        int
	conn4Size       int
	conn6Size       int
	errInfoSize     int
	recPacketSize   int
	packetSize      int
	pingPoolSize    int
	uuidPool        struct {
		arraySize int
		uuidSize  int
	}
}

func NewPoolWithOptions(f ...OpFunc) *Pool {
	return newPool(GetOptionWithOpFunc(f...))
}

func newPool(opt *Option) *Pool {
	return &Pool{
		option:      opt,
		readyChan:   make(chan *Ping, opt.readyChanSize),
		errChan:     make(chan *ErrInfo, opt.readyChanSize),
		onRevChan:   make(chan *Packet, opt.onRevChanSize),
		wg:          &sync.WaitGroup{},
		revWg:       &sync.WaitGroup{},
		readWg:      &sync.WaitGroup{},
		errWg:       &sync.WaitGroup{},
		done:        make(chan struct{}),
		conn4:       make(chan packetConn, opt.conn4Size),
		conn6:       make(chan packetConn, opt.conn6Size),
		errInfoPool: make(chan *ErrInfo, opt.errInfoSize),
		recPacket:   make(chan *packet, opt.recPacketSize),
		packet:      make(chan *Packet, opt.packetSize),
		ping:        make(chan *Ping, opt.pingPoolSize),
	}

}

func GetOptionWithOpFunc(f ...OpFunc) *Option {
	o := defaultOption()
	for _, fu := range f {
		fu(o)
	}
	return o
}

// ///////////////////////////////////////////// defaultOption /////////////////////////////////////
func defaultOption() *Option {
	r := rand.New(rand.NewSource(getSeed()))
	return &Option{
		readyChanSize:   ReadyChanSize,
		conn4Size:       Conn4,
		conn6Size:       Conn6,
		onRevChanSize:   OnRevChanSize,
		onRevWorkerNum:  OnRevChanSize,
		errChanSize:     ErrChanSize,
		recPacketSize:   RevPacketSize,
		packetSize:      PacketSize,
		currentMapSize:  CurrentMapSize,
		customerSize:    CustomerSize,
		resCustomerSize: ResCustomerSize,
		errInfoSize:     ErrorInfoSize,
		pingPoolSize:    PoolPingSize,
		readTimeout:     time.Millisecond * 100,
		pid:             r.Intn(math.MaxUint16),
		network:         IP,
		protocol:        Udp,
		dataSize:        timeSliceLength + trackerLength,
		mapTidyInterval: time.Hour,
		errCallBack:     defaultErrCallBack,
		Recover:         defaultRecover,
		OnRevFunc:       OnRevFunc,
		uuidPool: struct {
			arraySize int
			uuidSize  int
		}{arraySize: PoolArraySize, uuidSize: PoolUUIDSize},
	}
}

/////////////////////////////////////////////// set options /////////////////////////////////////

func WithOnRevFunc(f func(*Ping, *Packet)) OpFunc {
	return func(option *Option) {
		option.OnRevFunc = f
	}
}
func WithRecover(f func()) OpFunc {
	return func(option *Option) {
		option.Recover = f
	}
}

func WithErrCallBack(f func(ping *Ping, err error)) OpFunc {
	return func(option *Option) {
		option.errCallBack = f
	}
}

func WithPingPoolSize(pingPoolSize int) OpFunc {
	return func(option *Option) {
		option.pingPoolSize = pingPoolSize
	}
}

func WithUUIDPoolArraySize(arraySize int) OpFunc {
	return func(option *Option) {
		option.uuidPool.arraySize = arraySize
	}
}

func WithUUIDPoolUUIDSize(uuidSize int) OpFunc {
	return func(option *Option) {
		option.uuidPool.uuidSize = uuidSize
	}
}

func WithErrInfoSize(errPoolSize int) OpFunc {
	return func(option *Option) {
		option.errInfoSize = errPoolSize
	}
}
func WithPacketSize(packetSize int) OpFunc {
	return func(option *Option) {
		option.packetSize = packetSize
	}
}
func WithRevPacketSIze(revPtSi int) OpFunc {
	return func(option *Option) {
		option.recPacketSize = revPtSi
	}
}

func WithDataSize(size int) OpFunc {
	return func(option *Option) {
		option.dataSize = size
	}
}

func WithConn4Size(size int) OpFunc {
	return func(option *Option) {
		option.conn4Size = size
	}
}
func WithConn6Size(size int) OpFunc {
	return func(option *Option) {
		option.conn6Size = size
	}
}
func WithNetwork(s string) OpFunc {
	return func(option *Option) {
		option.network = s
	}
}
func WithProtocol(protocol string) OpFunc {
	return func(option *Option) {
		option.protocol = protocol
	}
}
func WithSource(s string) OpFunc {
	return func(option *Option) {
		option.source = s
	}
}
func WithOnRevWorkerNum(num int) OpFunc {
	return func(option *Option) {
		option.onRevWorkerNum = num
	}
}

func WithPid(pid int) OpFunc {
	return func(option *Option) {
		option.pid = pid
	}
}
func WithCurrentMapSize(num int) OpFunc {
	return func(option *Option) {
		option.currentMapSize = num
	}
}

func WithReadyChanSize(num int) OpFunc {
	return func(option *Option) {
		option.readyChanSize = num
	}
}
func WithCustomerSize(num int) OpFunc {
	return func(option *Option) {
		option.customerSize = num
	}
}
func WithResCustomerSize(num int) OpFunc {
	return func(option *Option) {
		option.resCustomerSize = num
	}
}
func WithReadTimeout(t time.Duration) OpFunc {
	return func(option *Option) {
		option.readTimeout = t
	}
}
func WithOnRevChanSize(num int) OpFunc {
	return func(option *Option) {
		option.onRevChanSize = num
	}
}
func WithResChanSize(num int) OpFunc {
	return func(option *Option) {
		option.ResChanSize = num
	}
}
func WithMapTidyInterval(t time.Duration) OpFunc {
	return func(option *Option) {
		option.mapTidyInterval = t
	}
}
func WithErrChanSize(num int) OpFunc {
	return func(option *Option) {
		option.errChanSize = num
	}
}

/////////////////////////////////////////////// default /////////////////////////////////////

func OnRevFunc(ping *Ping, pkt *Packet) {
	fmt.Println("----------------pkt Packet start---------------------")
	fmt.Printf("%#v\n", pkt)
	fmt.Println("----------------pkt Packet end-----------------------")
}

func defaultErrCallBack(ping *Ping, err error) {
	if err == nil {
		return
	}
	fmt.Println("----------------ping Ping start---------------------")
	fmt.Printf("%#v\n", ping)
	fmt.Println(err)
	fmt.Println("----------------ping Ping end-----------------------")

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
