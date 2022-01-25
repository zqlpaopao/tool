package pkg

type Options interface {
	apply(*option)
}

type option struct {
	conTime int
	print bool
	tryDump bool
	tryNum int32
}

type OptionFunc func(*option)

func(o OptionFunc)apply(option2 *option){
	o(option2)
}


func (o *option)clone()*option{
	cp := *o
	return &cp
}

func (o *option)WithOption(f... Options)*option{
	for _, v:= range f{
		v.apply(o)
	}
	return o
}


func WithOptConTime(t int)OptionFunc{
	return func(o *option) {
		o.conTime = t
	}
}

func WithOptPrint(b bool)OptionFunc{
	return func(o *option) {
		o.print = b
	}
}

func WithOptTryDump(b bool)OptionFunc{
	return func(o *option) {
		o.tryDump = b
	}
}

func WithOptTryNum(i int32)OptionFunc{
	return func(o *option) {
		o.tryNum = i
	}
}