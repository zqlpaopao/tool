package pkg

import "time"

type Proxy struct {
	CmdAble
	Cmder
	err error
}

// Lock -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Proxy) Lock(s string, duration time.Duration) (int64, error) {
	return p.Cmder.Lock(s, duration)
}

// UnLock -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Proxy) UnLock(lockName string) (int64, error) {
	return p.Cmder.UnLock(lockName)
}

// Renewal -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Proxy) Renewal(s string, duration time.Duration) (int64, error) {
	return p.Cmder.Renewal(s, duration)
}

// GetLockInfo  -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Proxy) GetLockInfo() (map[string]string, error) {
	return p.Cmder.GetLockInfo()
}

// Process -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Proxy) Process(cmd Cmder) {
	cmd.SetError(cmd.Ping())
}

// NewProxy -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func NewProxy(cmd Cmder) *Proxy {
	p := &Proxy{Cmder: cmd}
	p.CmdAble = p.Process
	return p
}

// Error -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Proxy) Error() error {
	return p.err
}

// Init -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Proxy) Init() *Proxy {
	p.err = p.Run(p.Cmder).Error()
	return p
}
