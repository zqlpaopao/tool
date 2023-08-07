package ping

type PoolData struct {
	ping chan *Ping
}

func NewPoolData(pingSi int) *PoolData {
	return &PoolData{ping: make(chan *Ping, pingSi)}
}

func (p *PoolData) GetPing() (res *Ping) {
	select {
	case res, _ = <-p.ping:
		return
	default:
		return &Ping{}
	}
}

func (p *PoolData) SetPing(res *Ping) {
	select {
	case p.ping <- res:
	default:

	}
}

// //////////////////////////////////////////////Error Pool /////////////////////////////////////
func (p *Pool) getErrRes() (res *ErrInfo) {
	select {
	case res, _ = <-p.errInfoPool:
		return
	default:
		return &ErrInfo{}
	}
}

func (p *Pool) SetErrRes(res *ErrInfo) {
	select {
	case p.errInfoPool <- res:
	default:

	}
}
