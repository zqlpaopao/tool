package ping

import (
	"fmt"
	pool "github.com/zqlpaopao/tool/data-any-pool/pkg"
	"golang.org/x/net/icmp"
	"math"
	"math/rand"
	"runtime"
	"runtime/debug"
	"time"
	"unsafe"
)

type Options interface {
	apply(*Option)
}

type OpFunc func(*Option)

func (o OpFunc) apply(opt *Option) {
	o(opt)
}

type Option struct {
	onRevFunc       func(ping *Ping, packet *Packet)
	errCallBack     func(*ErrInfo)
	Recover         func()
	network         string
	source          string
	protocol        string
	filter          Filter
	PacketPoolLen   int
	dataSize        int
	pingPoolLen     int
	uuidPoolLen     int
	uuidBytePoolLen int
	bodyLen         int
	msgLen          int
	msgByteLen      int
	readBuffer      int
	writeBuffer     int
	callbackWorker  int
	packetPoolLen   int
	pid             int
	errPoolLen      int
	resChanLen      int
	mapTidyInterval time.Duration
	errChLen        int
	pingMapLen      int
	prepareChLen    int
	sendWorker      int
	bpf             bool
	v6              bool
}

// NewPoolWithOptions -- --------------------------
// --> @Describe make the *Pool with the Option
// --> @params
// --> @return
// -- ------------------------------------
func NewPoolWithOptions(f ...OpFunc) *Pool {
	return newPool(GetOptionWithOpFunc(f...))
}

// newPool -- --------------------------
// --> @Describe make the default pool
// --> @params
// --> @return
// -- ------------------------------------
func newPool(opt *Option) *Pool {
	return &Pool{
		option:      opt,
		prepareChan: make(chan *Ping, opt.prepareChLen),
		errChan:     make(chan *ErrInfo, opt.errChLen),
		resChan:     make(chan *Packet, opt.resChanLen),
		errP:        pool.NewPool[*ErrInfo](0, opt.errPoolLen, func(_ int) *ErrInfo { return &ErrInfo{} }),
		PacketP:     pool.NewPool[*Packet](0, opt.PacketPoolLen, func(_ int) *Packet { return &Packet{} }),
		packetP:     pool.NewPool[*packet](0, opt.packetPoolLen, func(_ int) *packet { return &packet{} }),
		pingP:       pool.NewPool[*Ping](0, opt.pingPoolLen, func(_ int) *Ping { return &Ping{} }),
		body:        pool.NewPool[*icmp.Echo](0, opt.bodyLen, func(_ int) *icmp.Echo { return &icmp.Echo{} }),
		msg:         pool.NewPool[*icmp.Message](0, opt.msgLen, func(_ int) *icmp.Message { return &icmp.Message{} }),
		msgBytes:    pool.NewPool[unsafe.Pointer](0, opt.msgByteLen, func(_ int) unsafe.Pointer { var b []byte; return unsafe.Pointer(&b) }),
	}

}

// GetOptionWithOpFunc -- --------------------------
// --> @Describe make the Options with OpFunc
// --> @params
// --> @return
// -- ------------------------------------
func GetOptionWithOpFunc(f ...OpFunc) *Option {
	o := defaultOption()
	for _, fu := range f {
		fu(o)
	}
	return o
}

// ///////////////////////////////////////////// defaultOption /////////////////////////////////////
func defaultOption() *Option {
	var r = rand.New(rand.NewSource(getSeed()))
	return &Option{
		pid:             r.Intn(math.MaxUint16),
		prepareChLen:    PrepareChLen,
		resChanLen:      PrepareChLen,
		pingMapLen:      runtime.NumCPU(),
		PacketPoolLen:   PrepareChLen,
		packetPoolLen:   PrepareChLen,
		pingPoolLen:     PrepareChLen,
		bodyLen:         PrepareChLen,
		msgLen:          PrepareChLen,
		msgByteLen:      PrepareChLen,
		uuidPoolLen:     PrepareChLen,
		uuidBytePoolLen: PrepareChLen,
		errChLen:        ErrorInfoSize,
		errPoolLen:      ErrorInfoSize,
		callbackWorker:  PrepareChLen,
		sendWorker:      1,
		readBuffer:      1024 * 1024 * 20,
		writeBuffer:     1024 * 1024 * 20,
		Recover:         defaultRecover,
		errCallBack:     ErrorCallback,
		network:         IP,
		protocol:        Udp,
		source:          "",
		dataSize:        timeSliceLength + trackerLength,
		mapTidyInterval: time.Hour,
		bpf:             true,
		onRevFunc:       OnRevFunc,
	}
}

// WithPid *************************************** option *********************************//
//
//	-- --------------------------
//
// --> @Describe make the icmp Id []uint16
// --> @params
// --> @return
// -- ------------------------------------
func WithPid(pid int) OpFunc {
	return func(option *Option) {
		option.pid = pid
	}
}

// WithRecover -- --------------------------
// --> @Describe the panic with recover
// --> @params
// --> @return
// -- ------------------------------------
func WithRecover(f func()) OpFunc {
	return func(option *Option) {
		option.Recover = f
	}
}

// WithErrorCallback -- --------------------------
// --> @Describe the panic with recover
// --> @params
// --> @return
// -- ------------------------------------
func WithErrorCallback(f func(*ErrInfo)) OpFunc {
	return func(option *Option) {
		option.errCallBack = f
	}
}

// WithOnRevFunc -- --------------------------
// --> @Describe the func is call back func
// --> @params
// --> @return
// -- ------------------------------------
func WithOnRevFunc(f func(ping *Ping, packet2 *Packet)) OpFunc {
	return func(option *Option) {
		option.onRevFunc = f
	}
}

// WithPrepareChLen -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithPrepareChLen(size int) OpFunc {
	return func(option *Option) {
		option.prepareChLen = size
	}
}

// WithBodyLen -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithBodyLen(size int) OpFunc {
	return func(option *Option) {
		option.bodyLen = size
	}
}

// WithMsgLen -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithMsgLen(size int) OpFunc {
	return func(option *Option) {
		option.msgLen = size
	}
}

// WithPingMapLen -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithPingMapLen(size int) OpFunc {
	return func(option *Option) {
		option.pingMapLen = size
	}
}

// WithDataSize -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithDataSize(size int) OpFunc {
	return func(option *Option) {
		option.dataSize = size
	}
}

// WithErrChLen -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithErrChLen(size int) OpFunc {
	return func(option *Option) {
		option.errChLen = size
	}
}

// WithResChanLen -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithResChanLen(size int) OpFunc {
	return func(option *Option) {
		option.resChanLen = size
	}
}

// WithUuidPoolLen -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithUuidPoolLen(size int) OpFunc {
	return func(option *Option) {
		option.uuidPoolLen = size
	}
}

// WithUuidBytePoolLen -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithUuidBytePoolLen(size int) OpFunc {
	return func(option *Option) {
		option.uuidBytePoolLen = size
	}
}

// WithMapTidyInterval -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithMapTidyInterval(time time.Duration) OpFunc {
	return func(option *Option) {
		option.mapTidyInterval = time
	}
}

// WithNetwork -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithNetwork(network string) OpFunc {
	return func(option *Option) {
		option.network = network
	}
}

// WithProtocol -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithProtocol(protocol string) OpFunc {
	return func(option *Option) {
		option.protocol = protocol
	}
}

// WithSource -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithSource(source string) OpFunc {
	return func(option *Option) {
		option.source = source
	}
}

// WithV6 -- --------------------------
// --> @Describe is open v6 client
// --> @params
// --> @return
// -- ------------------------------------
func WithV6(v6 bool) OpFunc {
	return func(option *Option) {
		option.v6 = v6
	}
}

// WithBpf -- --------------------------
// --> @Describe is open v6 client
// --> @params
// --> @return
// -- ------------------------------------
func WithBpf(bpf bool) OpFunc {
	return func(option *Option) {
		option.bpf = bpf
	}
}

// WithCallbackWorker -- --------------------------
// --> @Describe is open v6 client
// --> @params
// --> @return
// -- ------------------------------------
func WithCallbackWorker(callBackNum int) OpFunc {
	return func(option *Option) {
		option.callbackWorker = callBackNum
	}
}

// WithSendWorker -- --------------------------
// --> @Describe is open v6 client
// --> @params
// --> @return
// -- ------------------------------------
func WithSendWorker(num int) OpFunc {
	return func(option *Option) {
		option.sendWorker = num
	}
}

// WithBPFFilter -- --------------------------
// --> @Describe is open v6 client
// --> @params
// --> @return
// -- ------------------------------------
func WithBPFFilter(filter Filter) OpFunc {
	return func(option *Option) {
		option.filter = filter
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

// ErrorCallback -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func ErrorCallback(err *ErrInfo) {
	fmt.Printf("tag %s ip %s err %v", err.Tag, err.Ping, err.Err)
}
