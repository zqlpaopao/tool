package pkg

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"
)

// IsClose state is close
func (o *OptionItem[T]) IsClose() bool {
	return atomic.LoadInt32(&o.close) == CLOSED
}

// check
func (o *OptionItem[T]) check() error {
	if o.hookFunc == nil {
		return ERRHookFuncIsEmpty
	}
	return nil
}

// Run tidy other info
func (o *OptionItem[T]) Run() {
	//o.wg.Add(o.handleGoNum)
	for i := 0; i < o.handleGoNum; i++ {
		go o.doing()
	}
}

// doing Handle by yourself every goroutine process
func (o *OptionItem[T]) doing() {
	var (
		task  = make([]T, 0, o.doingSize)
		timer = time.NewTimer(o.waitTime)
		//timer1 = time.NewTimer(o.loopTime)
	)
	defer o.savePanicFunc()
	for {
		select {
		case v, ok := <-o.itemCh:
			if !ok {
				if len(task) < 1 {
					goto END
				}
				o.hook(&task)
				task = []T{}
				goto END
			}
			task = append(task, v)
			if len(task) >= o.doingSize {
				o.hook(&task)
				task = make([]T, 0, o.doingSize)
			}
		case <-timer.C:
			if len(task) > 0 {
				o.hook(&task)
				task = make([]T, 0, o.doingSize)
			}
			timer.Reset(o.waitTime)
		//case <-timer1.C:
		//	if len(task) >= o.doingSize {
		//		o.hook(&task, &params)
		//		task, params = []interface{}{}, []interface{}{}
		//	}
		//	timer1.Reset(o.loopTime)
		default:
			//if atomic.LoadInt32(&o.close) == CLOSED {
			//	goto END
			//}
			time.Sleep(300 * time.Millisecond)
		}
	}
END:
	timer.Stop()
	o.wg.Done()
}

// hook Execute specific functions
func (o *OptionItem[T]) hook(task *[]T) {
	if o.endHook != nil {
		o.endHook(o.hookFunc(*task))
		return
	}
	o.hookFunc(*task)
}

// defaultSavePanic
func defaultSavePanic() {
	if err := recover(); nil != err {
		fmt.Println(err)
		fmt.Println(string(debug.Stack()))
	}
}
