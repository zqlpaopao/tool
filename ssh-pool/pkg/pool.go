package pkg

import (
	"fmt"
	"github.com/zqlpaopao/tool/format/src"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// Pool Configuration of Identical Connection pool
type Pool[T any] struct {
	cli    map[unsafe.Pointer]T
	option *PoolOption
	lock   *sync.RWMutex
	used   int32
}

// check  necessary parameters
func (p *Pool[T]) check() error {
	if p.option.name == "" {
		return ErrPoolTag
	}
	if p.option.item == nil {
		return ErrItemIsNil
	}
	return nil
}

// Submit Submitted to
func (p *Pool[T]) Submit(ssh *Ssh) {
	p.lock.Lock()
	p.cli[unsafe.Pointer(&ssh)] = ssh
	p.lock.Unlock()
}

// Object Managers of multiple Connection pool
type Object[T any] struct {
	Pool          map[string]*Pool[T]
	lock          *sync.RWMutex
	checkCh       chan Item
	ch            chan struct{}
	recover       func()
	checkNum      int
	heartbeatTime time.Duration
}

// NewObject make object
func NewObject[T any]() *Object[T] {
	return &Object[T]{
		Pool:          make(map[string]*Pool[T], 10),
		lock:          &sync.RWMutex{},
		checkCh:       make(chan Item, 100),
		checkNum:      CheckNum,
		ch:            make(chan struct{}),
		heartbeatTime: HeartbeatTime,
		recover:       SavePanic(),
	}
}

// SetCheckNum Set the number of checks
func (o *Object[T]) SetCheckNum(i int) *Object[T] {
	o.checkNum = i
	return o
}

// SetHeartbeatTime Setting the heartbeat time for examination -3
func (o *Object[T]) SetHeartbeatTime(i time.Duration) *Object[T] {
	o.heartbeatTime = i
	return o
}

// Do start work
func (o *Object[T]) Do() *Object[T] {
	for i := 0; i < o.checkNum; i++ {
		go o.check()
	}

	go o.loop()
	return o
}

// check Perform heartbeat detection
func (o *Object[T]) check() {
	defer o.recover()
	for {
		select {
		case _, ok := <-o.ch:
			if !ok {
				goto END
			}
		case v, ok := <-o.checkCh:
			if !ok {
				goto END
			}
			o.Heartbeat(v)
		}
	}

END:
}

// loop Recurrent update heartbeat
func (o *Object[T]) loop() {
	defer o.recover()

	var ts time.Duration
	if o.heartbeatTime-time.Second > 0 {
		ts = o.heartbeatTime - time.Second
	}

	var t = time.NewTicker(ts)
	for {
		select {
		case _, ok := <-o.ch:
			if !ok {
				goto END
			}

		case <-t.C:
		}
		o.push()
		t.Reset(o.heartbeatTime - 3*time.Second)
	}
END:
	t.Stop()

}

func (o *Object[T]) push() {
	o.lock.RLock()
	for k, v := range o.Pool {
		v.lock.Lock()
		for k1, v1 := range v.cli {
			if v1.IsDelete(v.option.maxLifeTime) {
				delete(o.Pool[k].cli, k1)
				continue
			}
			o.checkCh <- v1
		}
		v.lock.Unlock()
	}
	o.lock.RUnlock()
}

func (o *Object[T]) Heartbeat(p Item) {
	p.Heartbeat()
}

func (o *Object[T]) Submit(p *Pool[T]) (err error) {
	if err = p.check(); nil != err {
		return
	}
	o.lock.Lock()
	o.Pool[p.option.name] = p
	o.lock.Unlock()
	return
}

// Put  the Connection pool connection with the specified ID
func (o *Object[T]) Put(name string, p Item) {
	defer SavePanic()()

	o.lock.Lock()
	o.Pool[name].cli[unsafe.Pointer(&p)] = p
	o.lock.Unlock()

	atomic.AddInt32(&o.Pool[name].used, -1)
	p.SetNoUsing()
}

// Get  the Connection pool connection with the specified ID
func (o *Object[T]) Get(name string) (cli T, err error) {
	if _, ok := o.Pool[name]; !ok {
		err = ErrNotExist
		return
	}
	return o.GetItem(name, o.Pool[name])
}

// GetItem Get the Connection pool connection with the specified ID
func (o *Object[T]) GetItem(name string, p *Pool[T]) (cli T, err error) {
	defer SavePanic()()
	for {
		if len(p.cli) < 1 {
			s := p.option.item.Copy()
			if v, ok := s.(T); !ok {
				continue
			} else {
				if o.Pool[name].cli == nil {
					o.Pool[name].cli = make(map[unsafe.Pointer]Item, 10)
				}
				o.Pool[name].cli[unsafe.Pointer(&v)] = s
				s.SetUsing()
				atomic.AddInt32(&o.Pool[name].used, 1)
				return v, nil
			}
		}

		o.lock.Lock()
		for k1, v := range o.Pool[name].cli {
			if v.IsDelete(o.Pool[name].option.maxLifeTime) {
				delete(o.Pool[name].cli, k1)
				continue
			}
			if v.IsUse() {
				continue
			}
			if cli, ok := v.(T); ok {
				atomic.AddInt32(&o.Pool[name].used, 1)
				v.SetUsing()
				o.lock.Unlock()
				return cli, nil
			}
		}
		o.lock.Unlock()

		if int(atomic.LoadInt32(&o.Pool[name].used)) < o.Pool[name].option.maxConnNum {
			s := p.option.item.Copy()
			if v, ok := s.(T); !ok {
				continue
			} else {
				o.lock.Lock()
				o.Pool[name].cli[unsafe.Pointer(&v)] = s
				o.lock.Unlock()

				s.SetUsing()
				atomic.AddInt32(&o.Pool[name].used, 1)
				return v, nil
			}
		}
		//else {
		//	time.Sleep(time.Millisecond)
		//}

	}
}

// Close closing chan and  client
func (o *Object[T]) Close() {

}

// SavePanic default recover
func SavePanic() func() {
	return func() {
		if err := recover(); err != nil {
			src.PrintRed("PANIC-----")
			fmt.Println(err)
			fmt.Println(string(debug.Stack()))
		}
	}
}
