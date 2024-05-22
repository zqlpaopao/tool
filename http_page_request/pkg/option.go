package pkg

import (
	"sync"
	"time"
)

type OptionI interface {
	apply(*Option)
}

type Option struct {
	url        string
	limit      int
	customerGO int
	delayTime  time.Duration
	errStop    bool
}

type OPFunc func(*Option)

// apply -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (o OPFunc) apply(opt *Option) {
	o(opt)
}

// NewOption -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func NewOption(f ...OPFunc) *Option {
	opt := DefaultOption()
	for _, v := range f {
		v(opt)
	}
	return opt
}

// DefaultOption -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func DefaultOption() *Option {
	return &Option{
		limit:      10000,
		url:        "",
		customerGO: 10,
		delayTime:  0,
		errStop:    true,
	}
}

// WithUrl -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithUrl(url string) OPFunc {
	return func(option *Option) {
		option.url = url
	}
}

// WithLimit -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithLimit(limit int) OPFunc {
	return func(option *Option) {
		option.limit = limit
	}
}

// WithCustomerGO -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithCustomerGO(customerGO int) OPFunc {
	return func(option *Option) {
		option.customerGO = customerGO
	}
}

// WithDelayTime -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithDelayTime(delayTime time.Duration) OPFunc {
	return func(option *Option) {
		option.delayTime = delayTime
	}
}

// WithErrStop -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func WithErrStop(errStop bool) OPFunc {
	return func(option *Option) {
		option.errStop = errStop
	}
}

// NewPageHandlerWithOptions -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func NewPageHandlerWithOptions[P, R any](
	params P,
	resFn func(url string, p P, page, limit int, res R) *ErrInfo[P],
	totalFn func(url string, p P) int,
	res R,
	f ...OPFunc,
) *PageHandler[P, R] {
	opt := NewOption(f...)
	return &PageHandler[P, R]{
		params:    params,
		opt:       opt,
		totalFunc: totalFn,
		resFunc:   resFn,
		res:       res,
		loopChan:  make(chan int, opt.customerGO*2),
		errLock:   &sync.Mutex{},
		err: &ErrInfo[P]{
			Err:    nil,
			Url:    "",
			Params: params,
		},
	}
}
