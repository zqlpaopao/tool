package pkg

import (
	"sync"
	"time"
)

//HookFunc Callbacks need to handle specific functions
type HookFunc func(task []interface{}) bool

//EndFunc Functions that handle callbacks each time
type EndFunc func(b bool, i ...interface{})

//SavePanic Functions that handle exception panic
type SavePanic func(i interface{})

type Option interface {
	apply(*option)
}

type option struct {
	close         int32
	doingSize     int
	handleGoNum   int
	chanSize      int
	waitTime      time.Duration
	loopTime      time.Duration
	itemCh        chan interface{}
	hookFunc      HookFunc
	endHook       EndFunc
	savePanicFunc SavePanic
	wg            sync.WaitGroup
}

type OpFunc func(*option)

func NewOption(opt ...Option) *option {
	return clone().WithOptions(opt...)
}

//apply assignment function entity
func (o OpFunc) apply(opt *option) {
	o(opt)
}

//clone  new object
func clone() *option {
	return &option{
		close:         OPENED,
		doingSize:     DoingSize,
		handleGoNum:   HandleGoNum,
		chanSize:      ChanSize,
		waitTime:      WaitTime,
		loopTime:      LoopTime,
		itemCh:        nil,
		hookFunc:      nil,
		endHook:       nil,
		savePanicFunc: defaultSavePanic,
		wg:            sync.WaitGroup{},
	}

}

//WithOptions Execute assignment function entity
func (o *option) WithOptions(opt ...Option) *option {
	for _, v := range opt {
		v.apply(o)
	}
	o.initParams()
	return o
}

//initParams Initialization parameters
func (o *option) initParams() {
	o.itemCh = make(chan interface{}, o.chanSize)
}

//WithDoingSize How much to start processing default 100
func WithDoingSize(size int) OpFunc {
	return func(o *option) {
		o.doingSize = size
	}
}

//WithHandleGoNum Number of goroutine processed default 100
func WithHandleGoNum(num int) OpFunc {
	return func(o *option) {
		o.handleGoNum = num
	}
}

//WithWaitTime How often to wait default 2s
func WithWaitTime(waitTime time.Duration) OpFunc {
	return func(o *option) {
		o.waitTime = waitTime
	}
}

//WithChanSize chan size default 100
func WithChanSize(size int) OpFunc {
	return func(o *option) {
		o.chanSize = size
	}
}

//WithLoopTime How often is the length checked and whether it is implemented default 1s
func WithLoopTime(loopTime time.Duration) OpFunc {
	return func(o *option) {
		o.loopTime = loopTime
	}
}

//WithHookFunc callback func
func WithHookFunc(hookFunc HookFunc) OpFunc {
	return func(o *option) {
		o.hookFunc = hookFunc
	}
}

//WithEndHook callback end func
func WithEndHook(endHook EndFunc) OpFunc {
	return func(o *option) {
		o.endHook = endHook
	}
}

//WithSavePanic save panic
func WithSavePanic(savePanicFunc SavePanic) OpFunc {
	return func(o *option) {
		o.savePanicFunc = savePanicFunc
	}
}
