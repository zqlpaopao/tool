package pkg

import (
	"fmt"
	"os"
	"runtime/debug"
)

type OptionFuncA interface {
	apply(*Option)
}

type OptionA struct {
	srcPort  int
	dstPort  int
	chanBuf  int
	worker   int
	protocol uint16
	netMark  string
	callback func(b []byte) bool
	recover  func()
}

type OPFuncA func(option *OptionA)

func (o OPFuncA) applyA(option *OptionA) {
	o(option)
}

func NewScanWithOption(opt ...OPFuncA) *Scanner {
	return cloneA(opt...)
}

func cloneA(opt ...OPFuncA) *Scanner {
	o := &OptionA{
		srcPort:  SrcPort,
		dstPort:  DstPort,
		protocol: Protocol,
		netMark:  Mark,
		chanBuf:  ChanSize,
		worker:   WorkerNum,
		callback: DefaultCallBack,
		recover:  recovers,
	}

	for i := 0; i < len(opt); i++ {
		opt[i].applyA(o)
	}
	return NewScanner(o)
}

func WithBinBash(sh string) OPFunc {
	return func(option *Option) {
		option.binBash = sh
	}
}

func WithC(c string) OPFunc {
	return func(option *Option) {
		option.C = c
	}
}

func WithCmd(cmdStr string) OPFunc {
	return func(option *Option) {
		option.cmd = cmdStr
	}
}

func WithWorkerNum(num int) OPFunc {
	return func(option *Option) {
		option.workerNum = num
	}
}

func DefaultCallBack(b []byte) bool {
	fmt.Println(string(b))
	return true
}
func recovers() {
	if e := recover(); e != nil {
		fmt.Println(e)
		fmt.Println(string(debug.Stack()))
		os.Exit(2)
	}
}
