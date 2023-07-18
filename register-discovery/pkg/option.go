package pkg

type Register interface {
	Register(*RegisterInfo)
	UnRegister()
}

type Discovery interface {
	Discovery()
	UnDiscovery()
}

type Abstract interface {
	Register
	Discovery
	Error() error
}

type RegisterInfo struct {
	RegisterInfo string
	PushTime     int
}

type Proxy struct {
	ety Abstract
}

func NewProxy(proxy Abstract) Abstract {
	return &Proxy{ety: proxy}
}

func (p *Proxy) Register(msg *RegisterInfo) {
	p.ety.Register(msg)
}

func (p *Proxy) UnRegister() {
	p.ety.UnRegister()
}
func (p *Proxy) Discovery() {
	p.ety.Discovery()
}
func (p *Proxy) UnDiscovery() {
	p.ety.UnDiscovery()
}

func (p *Proxy) Error() error {
	return p.ety.Error()
}
