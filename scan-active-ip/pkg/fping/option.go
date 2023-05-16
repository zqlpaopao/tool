package fping

import (
	"fmt"
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

// WithBinBash Set the bash for the system
func WithBinBash(sh string) OPFunc {
	return func(option *Option) {
		option.binBash = sh
	}
}

// WithC System parameters executed
func WithC(c string) OPFunc {
	return func(option *Option) {
		option.C = c
	}
}

// WithCmd cmd
func WithCmd(cmdStr string) OPFunc {
	return func(option *Option) {
		option.cmd = cmdStr
	}
}

// WithWorkerNum scan workers
func WithWorkerNum(num int) OPFunc {
	return func(option *Option) {
		option.workerNum = num
	}
}

// WithDefaultCallBack DefaultCallBack
func WithDefaultCallBack(f func([]byte) bool) OPFunc {
	return func(option *Option) {
		option.callback = f
	}
}

// DefaultCallBack Default callback function
func DefaultCallBack(b []byte) bool {
	fmt.Println(string(b))
	return true
}

// panic recovers
func recovers() {
	if e := recover(); e != nil {
		fmt.Println(e)
		fmt.Println(string(debug.Stack()))
	}
}
