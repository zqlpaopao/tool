package pkg

import "time"

type Proxy struct {
	CmdAble
	Cmder
	err error
}

func (p *Proxy) Lock(s string, duration time.Duration) (int64, error) {
	return p.Cmder.Lock(s, duration)
}

func (p *Proxy) UnLock(lockName string) (int64, error) {
	return p.Cmder.UnLock(lockName)
}

func (p *Proxy) Renewal(s string, duration time.Duration) (int64, error) {
	return p.Cmder.Renewal(s, duration)
}

func (p *Proxy) Process(cmd Cmder) {
	cmd.SetError(cmd.Ping())
}

func NewProxy(cmd Cmder) *Proxy {
	p := &Proxy{Cmder: cmd}
	p.CmdAble = p.Process
	return p
}

func (p *Proxy) Error() error {
	return p.err
}

func (p *Proxy) Init() *Proxy {
	p.err = p.Run(p.Cmder).Error()
	return p
}
