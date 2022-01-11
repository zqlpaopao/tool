package src

type Options interface {
	apply(*option)
}

type option struct{
	notAuth bool
	witBranch bool
	print bool
	tag string
}

type OptionFunc func(option2 *option)

func(f OptionFunc)apply(option2 *option){
	f(option2)
}

func(o *option)clone()*option{
	c := *o
	return &c
}

func(o option)WithOptions(f... Options)*option{
	c := o.clone()
	for _, opt := range f {
		opt.apply(c)
	}
	return c
}

func NewOptions(f...Options)*option{
	o := &option{}
	return o.WithOptions(f...)
}

//WithBranch get branch info
func WithBranch(b bool)OptionFunc{
	return func(option2 *option) {
		option2.witBranch = b
	}
}
//WithNotAuth auth info
func WithNotAuth(b bool)OptionFunc{
	return func(option2 *option) {
		option2.notAuth = b
	}
}

//WithTag describe info
func WithTag(tag string)OptionFunc{
	return func(option2 *option) {
		option2.tag = tag
	}
}

//WithPrint print version
func WithPrint(b bool)OptionFunc{
	return func(option2 *option) {
		option2.print = b
	}
}