package src

type ErrorHandle struct {
	tag  int
	msg  string
	args []interface{}
}

var (
	errHandler *ErrorHandle
	//resource pool Reuse objects and avoid frequent affectionate objects
	buffErrSize      *BufferErrPool
	buffCallBackSize *CallBackBufferPool
	callBackFunc     func(level int, tag string, info *CallBack)
)

// Reset Fallback object initialization status
func (e *ErrorHandle) Reset() {
	e.msg, e.args, e.tag = "", nil, 0
}

// Reset Fallback object initialization status
func (e *CallBack) Reset() {
	e.Ip, e.Params, e.Msg = "", "", ""
}

// NewBuffErrSize Initialize the size of the resource pool
func NewBuffErrSize(i int) {
	buffErrSize = NewBufferErrPool(i)
}

// NewCallBackBuffSize Initialize the size of the resource pool
func NewCallBackBuffSize(i int) {
	buffCallBackSize = NewCallBackBufferPool(i)
}

// NewCallBackFunc Initialize the size of the call back func
func NewCallBackFunc() {
	callBackFunc = DefaultCallBackFunc
}

// initParams init info args
func (e *ErrorHandle) initParams(msg string, level int, args ...interface{}) (c *ErrorHandle) {
	c = buffErrSize.Get()
	c.msg, c.tag, c.args = msg, level, args
	return c
}
