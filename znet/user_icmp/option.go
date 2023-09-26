package user_icmp

import (
	"fmt"
	pool "github.com/zqlpaopao/tool/data-any-pool/pkg"
	"math"
	"math/rand"
	"runtime"
	"runtime/debug"
)

type Options interface {
	apply(*Option)
}

type OpFunc func(*Option)

func (o OpFunc) apply(opt *Option) {
	o(opt)
}

type Option struct {
	onRevFunc              func(packet *Packet)
	errCallBack            func(*ErrInfo)
	Recover                func()
	source                 string
	filter                 Filter
	PacketPoolLen          int
	ReceiveMMsgLen         int
	errChLen               int
	prepareV4ChLen         int
	prepareV6ChLen         int
	resChanLen             int
	ttl                    int
	dataSize               int
	poolSize               int
	readBuffer             int
	writeBuffer            int
	callbackWorker         int
	packetPoolLen          int
	pid                    uint16
	errPoolLen             int
	sendWorker             int
	fdReadWriteTimeOutNesc int64
	bpf                    bool
	v6                     bool
	everyTTL               bool
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
		option:        opt,
		prepareChanV4: make(chan *Ping, opt.prepareV4ChLen),
		errChan:       make(chan *ErrInfo, opt.errChLen),
		resChan:       make(chan *ReceiveMMsg, opt.resChanLen),
		errP:          pool.NewPool[*ErrInfo](0, opt.errPoolLen, func(_ int) *ErrInfo { return &ErrInfo{} }),
		PacketP:       pool.NewPool[*Packet](0, opt.callbackWorker*2, func(_ int) *Packet { return &Packet{} }),
		receiveMMsgPool: pool.NewPool[*ReceiveMMsg](opt.dataSize,
			opt.ReceiveMMsgLen,
			func(size int) *ReceiveMMsg {
				if runtime.GOOS == "darwin" {
					size = size + IpHeader + NoIcmpAndDarwin
				} else {
					size = size + IpHeader
				}
				return &ReceiveMMsg{Data: make([]byte, size)}
			}),
		pingP: pool.NewPool[*Ping](0, opt.poolSize, func(size int) *Ping { return &Ping{Size: 20, Ipv4: true, pid: opt.pid} }),
		seqP:  pool.NewPool[[]byte](2, opt.poolSize, func(size int) []byte { return make([]byte, size) }),
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
	//mac read 8192
	//linux 8388608 write 8388608
	//default  12
	return &Option{
		pid:                    uint16(r.Intn(math.MaxUint16)),
		poolSize:               PrepareChLen,
		ReceiveMMsgLen:         PrepareChLen * 2,
		errPoolLen:             ErrorInfoSize,
		errChLen:               ErrorInfoSize,
		callbackWorker:         PrepareChLen,
		prepareV4ChLen:         PrepareChLen,
		resChanLen:             PrepareChLen,
		fdReadWriteTimeOutNesc: 2000000000,
		sendWorker:             1,
		ttl:                    60,
		readBuffer:             0,
		writeBuffer:            0,
		Recover:                defaultRecover,
		errCallBack:            ErrorCallback,
		source:                 "0.0.0.0",
		dataSize:               0,
		onRevFunc:              OnRevFunc,
		bpf:                    true,
	}
}

// WithPid *************************************** option *********************************//
//
//	-- --------------------------
//
// --> @Describe make the user_icmp Id []uint16
// --> @params
// --> @return
// -- ------------------------------------
func WithPid(pid uint16) OpFunc {
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
func WithOnRevFunc(f func(packet2 *Packet)) OpFunc {
	return func(option *Option) {
		option.onRevFunc = f
	}
}

// WithNoBpf -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithNoBpf() OpFunc {
	return func(option *Option) {
		option.bpf = false
	}
}

// WithV6 -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithV6() OpFunc {
	return func(option *Option) {
		option.v6 = true
	}
}

// WithEveryTTL -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithEveryTTL() OpFunc {
	return func(option *Option) {
		option.everyTTL = true
	}
}

// WithReceiveMMsgLen -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithReceiveMMsgLen(size int) OpFunc {
	return func(option *Option) {
		option.ReceiveMMsgLen = size
	}
}

// WithFdReadWriteTimeOutNesc -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithFdReadWriteTimeOutNesc(size int64) OpFunc {
	return func(option *Option) {
		option.fdReadWriteTimeOutNesc = size
	}
}

// WithReadBuffer -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithReadBuffer(size int) OpFunc {
	return func(option *Option) {
		option.readBuffer = size
	}
}

// WithWriteBuffer -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithWriteBuffer(size int) OpFunc {
	return func(option *Option) {
		option.writeBuffer = size
	}
}

// WithReadWriteTime -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithReadWriteTime(size int) OpFunc {
	return func(option *Option) {
		option.writeBuffer = size
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

// WithPoolSize -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithPoolSize(size int) OpFunc {
	return func(option *Option) {
		option.poolSize = size
	}
}

// WithPrepareChV4Len -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithPrepareChV4Len(size int) OpFunc {
	return func(option *Option) {
		option.prepareV4ChLen = size
	}
}

// WithPrepareChV6Len -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithPrepareChV6Len(size int) OpFunc {
	return func(option *Option) {
		option.prepareV6ChLen = size
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

func OnRevFunc(pkt *Packet) {
	fmt.Println("----------------pkt Packet start---------------------")
	fmt.Printf("%#v\n", pkt)
	fmt.Println("----------------pkt Packet end-----------------------")
}

func defaultErrCallBack(ping *Ping, err error) {
	if err == nil {
		return
	}
	fmt.Println("----------------user_icmp Ping start---------------------")
	fmt.Printf("%#v\n", ping)
	fmt.Println(err)
	fmt.Println("----------------user_icmp Ping end-----------------------")

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
