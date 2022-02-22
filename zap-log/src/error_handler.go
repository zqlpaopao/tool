package src

type ErrorHandle struct {
	tag  int
	msg  string
	args []interface{}
}

var (
	errHandler *ErrorHandle
	//resource pool Reuse objects and avoid frequent affectionate objects
	buffErrSize *BufferErrPool
)

//Reset Fallback object initialization status
func (e *ErrorHandle) Reset() {
	e.msg, e.args, e.tag = "", nil, 0
}

//NewBuffErrSize Initialize the size of the resource pool
func NewBuffErrSize(i int) {
	buffErrSize = NewBufferErrPool(i)
}

//initParams init info args
func (e ErrorHandle) initParams(msg string, level int, args ...interface{}) (c *ErrorHandle) {
	c = buffErrSize.Get()
	c.msg, c.tag, c.args = msg, level, args
	return c
}
