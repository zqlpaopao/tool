package pkg

import (
	"context"
	"errors"
	"sync"
)

type ErrInfo[P any] struct {
	Err    error
	Params P
	Url    string
}

type PageHandler[P, R any] struct {
	params    P
	res       R
	ctx       context.Context
	opt       *Option
	err       *ErrInfo[P]
	totalFunc func(url string, p P) int
	resFunc   func(url string, p P, page, limit int, res R) *ErrInfo[P]
	loopChan  chan int
	cancel    context.CancelFunc
	errLock   *sync.Mutex
	total     int
}

// DO -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *PageHandler[P, R]) DO() *ErrInfo[P] {
	p.Check()
	p.MakeLoop()
	return p.err
}

// Check -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *PageHandler[P, R]) Check() {
	if p.opt.url == "" {
		p.err.Err = errors.New("url is empty")
		return
	}
	p.total = p.totalFunc(p.opt.url, p.params)
	if p.total < 1 {
		p.err.Err = errors.New("get the total is less 1")
		return
	}
	p.ctx,
		p.cancel =
		context.WithCancel(context.Background())

}

// MakeLoop -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *PageHandler[P, R]) MakeLoop() {
	if p.err.Err != nil {
		return
	}
	var (
		wg   = &sync.WaitGroup{}
		loop = (p.total + p.opt.limit) / p.opt.limit
	)
	p.CustomerLoop(wg)
	for i := 0; i < loop; i++ {
		p.loopChan <- i
	}
	close(p.loopChan)
	wg.Wait()
}

// CustomerLoop -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *PageHandler[P, R]) CustomerLoop(wg *sync.WaitGroup) {
	wg.Add(p.opt.customerGO)
	for i := 0; i < p.opt.customerGO; i++ {
		go p.Loop(wg)
	}
}

// Loop -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *PageHandler[P, R]) Loop(wg *sync.WaitGroup) {
	for {
		select {
		case v, ok := <-p.loopChan:
			if !ok {
				goto END
			}
			p.ItemRequest(v)
		case <-p.ctx.Done():
			goto END
		}
	}
END:
	wg.Done()
}

// ItemRequest -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *PageHandler[P, R]) ItemRequest(page int) {
	var err *ErrInfo[P]
	if err = p.resFunc(
		p.opt.url,
		p.params,
		page,
		p.opt.limit,
		p.res); err != nil && err.Err != nil && p.opt.errStop {
		p.errLock.Lock()
		p.err = err
		p.errLock.Unlock()
		p.cancel()
	}
}
