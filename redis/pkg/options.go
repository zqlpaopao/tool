package pkg

import "time"

type Option interface {
	apply(opt *option)
}

type OpFunc func(*option)

//retry options
type option struct {
	readTimeout time.Duration
	writeTimeout time.Duration
	idleTimeout time.Duration
	poolTimeOut time.Duration
	poolSize int
}


//apply assignment function entity
func(o OpFunc)apply(opt *option){
	o(opt)
}

//NewOption make option
func NewOption(opt ...Option)*option{
	o := &option{
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		idleTimeout:  idleTimeout,
		poolTimeOut:  poolTimeOut,
		poolSize:     poolSize,
	}
	return o.WithOptions(opt...)
}

//clone  new object
func(o *option)clone()*option{
	cp :=*o
	return &cp
}

//WithOptions Execute assignment function entity
func(o option)WithOptions(opt ...Option)*option{
	c := o.clone()
	for _ ,v := range opt{
		v.apply(c)
	}
	return c
}

//WithReadTimeout Set the number of ReadTimeout. The default is  5s
func WithReadTimeout(readTimeout time.Duration)OpFunc{
	return func(o *option) {
		o.readTimeout = readTimeout
	}
}

//WithWriteTimeout Set the WriteTimeout. The default is 5s
func WithWriteTimeout(writeTimeout time.Duration)OpFunc{
	return func(o *option) {
		o.writeTimeout = writeTimeout
	}
}

//WithIdleTimeout Set the IdleTimeout. The default value is 60s
func WithIdleTimeout(idleTimeout time.Duration) OpFunc {
	return func(o *option) {
		o.idleTimeout = idleTimeout
	}
}
//WithPoolTimeOut default 60s
func WithPoolTimeOut(poolTimeOut time.Duration) OpFunc {
	return func(o *option) {
		o.poolTimeOut = poolTimeOut
	}
}

//WithPoolSize set pool size default 20
func WithPoolSize(poolSize int) OpFunc {
	return func(o *option) {
		o.poolSize = poolSize
	}
}
