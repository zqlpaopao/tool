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

// //////////////////////////////////////////////packet Pool /////////////////////////////////////
func (p *Pool) getRevPacketRes() (res *packet) {
	select {
	case res, _ = <-p.recPacket:
		return
	default:
		return &packet{}
	}
}

func (p *Pool) SetRevPacketRes(res *packet) {
	select {
	case p.recPacket <- res:
	default:

	}
}

// //////////////////////////////////////////////Packet Pool /////////////////////////////////////
func (p *Pool) getPacketRes() (res *Packet) {
	select {
	case res, _ = <-p.packet:
		return
	default:
		return &Packet{}
	}
}

func (p *Pool) SetPacketRes(res *Packet) {
	select {
	case p.packet <- res:
	default:

	}
}
