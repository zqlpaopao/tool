package pkg

import "os"

// EndFunc Functions that handle callbacks each time
type EndFunc func(res *[]map[string]interface{})

// SavePanic Functions that handle exception panic
type SavePanic func(i interface{})

type OptionFunc interface {
	apply(*Option)
}

type File interface {
	Open(string) (*os.File, error)
	Close() error
}

type PublicFunc interface {
	Error() []error
	Code() int
	GetResp() *Resp
}

type Operation interface {
	PublicFunc
	ParamsCheck()
	Init()
	Doing()
}

// Proxy Agent of device operation ，including port operation
// configuration operation and other device-related operations
type Proxy struct {
	err       error
	operation Operation
}

// NewProxy Initialize agent
func NewProxy(operation Operation) *Proxy {
	return &Proxy{operation: operation}
}

// ParamsCheck init params check
func (d *Proxy) ParamsCheck() {
	d.operation.ParamsCheck()
}

// Init  params init
func (d *Proxy) Init() {
	d.operation.Init()
}

// Doing  is doing
func (d *Proxy) Doing() {
	d.operation.Doing()
}

// GetResp result
func (d *Proxy) GetResp() *Resp {
	return d.operation.GetResp()
}

// Get Device related get operation
func (d *Proxy) Error() []error {
	return d.operation.Error()
}

// Code Device related get operation
func (d *Proxy) Code() int {
	return d.operation.Code()
}

// Do The agent performs real entity operation
func (d *Proxy) Do() *Proxy {
	d.ParamsCheck()
	d.Init()
	d.Doing()
	return d
}
