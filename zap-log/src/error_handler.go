package src

var errHandler *ErrorHandle

func (e *ErrorHandle)clone()*ErrorHandle{
	c := *e
	return &c
}

func(e ErrorHandle)initParams(msg string, level int,args ...interface{})*ErrorHandle{
	c := errHandler.clone()
	c.msg,c.tag,c.args  = msg,level,args
	return c
}