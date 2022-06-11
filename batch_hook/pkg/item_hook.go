package pkg

import (
	"fmt"
	"runtime/debug"
	"sync/atomic"
	"time"
)

//IsClose state is close
func (o *option) IsClose() bool {
	return atomic.LoadInt32(&o.close) == CLOSED
}

//check
func (o *option) check() error {
	if o.hookFunc == nil {
		return ERRHookFuncIsEmpty
	}
	return nil
}

//Run tidy other info
func (o *option) Run() {
	o.wg.Add(o.handleGoNum)
	for i := 0; i < o.handleGoNum; i++ {
		go o.doing()
	}
}

//doing Handle by yourself every goroutine process
func (o *option) doing() {
	var (
		task   []interface{}
		timer  = time.NewTimer(o.waitTime)
		//timer1 = time.NewTimer(o.loopTime)
	)
	if o.savePanicFunc != nil {
		defer o.savePanicFunc(task)
	}
	for {
		select {
		case v ,ok := <-o.itemCh:
			if !ok{
				if len(task)<1{
					goto END
				}
				o.hook(&task)
				task = []interface{}{}
				goto END
			}
			task = append(task, v)
			if len(task) >= o.doingSize {
				o.hook(&task)
				task = []interface{}{}
			}
		case <-timer.C:
			if len(task) > 0 {
				o.hook(&task)
				task = []interface{}{}
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
			time.Sleep(time.Second)
		}
	}
END:
	o.wg.Done()
}

//hook Execute specific functions
func (o *option) hook(task *[]interface{}) {
	if o.endHook != nil {
		o.endHook(o.hookFunc(*task))
		return
	}
	o.hookFunc(*task)
}

//defaultSavePanic
func defaultSavePanic(i interface{}) {
	if err := recover(); nil != err {
		fmt.Println(i)
		fmt.Println(err)
		fmt.Println(string(debug.Stack()))
	}
}
