package pkg

import "time"

type Option interface {
	apply(*option)
}

//options public
type option struct {
	//addr        string
	//port        string
	maxNum      int
	connTimeout time.Duration
	factory     Factory
}

type OpFunc func(option2 *option)

//apply assignment function entity
func (f OpFunc) apply(opt *option) {
	f(opt)
}

// make new option
func (o *option) clone() *option {
	c := *o
	return &c
}

//WhitOptions Execute assignment function entity
func (o option) WhitOptions(f ...Option) *option {
	c := o.clone()
	for _, v := range f {
		v.apply(c)
	}
	return c
}

//NewOption make new option
func NewOption(f ...Option) *option {
	c := &option{
		maxNum:      maxNum,
		connTimeout: timeOut,
	}
	return c.WhitOptions(f...)
}

////WithAddr Assignment address
//func WithAddr(addr string) OpFunc {
//	return func(opt *option) {
//		opt.addr = addr
//	}
//}
//
////WithPort Assignment port
//func WithPort(port string) OpFunc {
//	return func(opt *option) {
//		opt.port = port
//	}
//}

//WithMaxNum Maximum number of connections assigned default 10
func WithMaxNum(maxNum int) OpFunc {
	return func(opt *option) {
		opt.maxNum = maxNum
	}
}

//WithMakeConn Maximum number of connections assigned default 10
func WithMakeConn(f Factory) OpFunc {
	return func(opt *option) {
		opt.factory = f
	}
}

//WithTimeout conn timeout default 10s
func WithTimeout(t time.Duration) OpFunc {
	return func(opt *option) {
		opt.connTimeout = t
	}
}
