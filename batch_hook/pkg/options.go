package pkg

import (
	"sync"
	"time"
)

// HookFunc Callbacks need to handle specific functions
type HookFunc[T any] func(task []T) bool

// EndFunc Functions that handle callbacks each time
type EndFunc[T any] func(b bool, i ...T)

// SavePanic Functions that handle exception panic
type SavePanic func()

type Option[T any] interface {
	apply(*OptionItem[T])
}

type OptionItem[T any] struct {
	close         int32
	doingSize     int
	handleGoNum   int
	chanSize      int
	waitTime      time.Duration
	loopTime      time.Duration
	itemCh        chan T
	hookFunc      HookFunc[T]
	endHook       EndFunc[T]
	savePanicFunc SavePanic
	wg            sync.WaitGroup
}

type OpFunc[T any] func(*OptionItem[T])

func NewOption[T any](opt ...Option[T]) *OptionItem[T] {
	return clone[T]().WithOptions(opt...)
}

// apply assignment function entity
func (o OpFunc[T]) apply(opt *OptionItem[T]) {
	o(opt)
}

// clone  new object
func clone[T any]() *OptionItem[T] {
	return &OptionItem[T]{
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

// WithOptions Execute assignment function entity
func (o *OptionItem[T]) WithOptions(opt ...Option[T]) *OptionItem[T] {
	for _, v := range opt {
		v.apply(o)
	}
	o.initParams()
	return o
}

// initParams Initialization parameters
func (o *OptionItem[T]) initParams() {
	o.itemCh = make(chan T, o.chanSize)
}

// WithDoingSize How much to start processing default 100
func WithDoingSize[T any](size int) OpFunc[T] {
	return func(o *OptionItem[T]) {
		o.doingSize = size
	}
}

// WithHandleGoNum Number of goroutine processed default 100
func WithHandleGoNum[T any](num int) OpFunc[T] {
	return func(o *OptionItem[T]) {
		o.handleGoNum = num
	}
}

// WithWaitTime How often to wait default 2s
func WithWaitTime[T any](waitTime time.Duration) OpFunc[T] {
	return func(o *OptionItem[T]) {
		o.waitTime = waitTime
	}
}

// WithChanSize chan size default 100
func WithChanSize[T any](size int) OpFunc[T] {
	return func(o *OptionItem[T]) {
		o.chanSize = size
	}
}

// WithLoopTime How often is the length checked and whether it is implemented default 1s
func WithLoopTime[T any](loopTime time.Duration) OpFunc[T] {
	return func(o *OptionItem[T]) {
		o.loopTime = loopTime
	}
}

// WithHookFunc callback func
func WithHookFunc[T any](hookFunc HookFunc[T]) OpFunc[T] {
	return func(o *OptionItem[T]) {
		o.hookFunc = hookFunc
	}
}

// WithEndHook callback end func
func WithEndHook[T any](endHook EndFunc[T]) OpFunc[T] {
	return func(o *OptionItem[T]) {
		o.endHook = endHook
	}
}

// WithSavePanic save panic
func WithSavePanic[T any](savePanicFunc SavePanic) OpFunc[T] {
	return func(o *OptionItem[T]) {
		o.savePanicFunc = savePanicFunc
	}
}
