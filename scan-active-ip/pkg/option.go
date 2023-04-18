package pkg

import (
	"fmt"
	"os"
	"runtime/debug"
)

type OptionFunc interface {
	apply(*Option)
}

type Option struct {
	binBash   string
	C         string
	cmd       string
	workerNum int
	chanBuf   int
	callback  func(b []byte) bool
	recover   func()
}

type OPFunc func(option *Option)

func (o OPFunc) apply(option *Option) {
	o(option)
}

func NewScanActiveWithOption(opt ...OPFunc) *ScanActive {
	return clone(opt...)
}

func clone(opt ...OPFunc) *ScanActive {
	o := &Option{
		binBash:   BinBash,
		C:         C,
		cmd:       Cmd,
		workerNum: WorkerNum,
		chanBuf:   ChanBuf,
		callback:  DefaultCallBack,
		recover:   recovers,
	}

	for i := 0; i < len(opt); i++ {
		opt[i].apply(o)
	}
	return NewScanActive(o)
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
