package pkg

// initConn Initialize ipv4 and ipv6 monitors for icmp
import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"sync/atomic"
	"time"
)

// initConn Initialize the Connection pool
func (p *PingMan) initConn() (err error) {
	if p.conn4, err = icmp.ListenPacket(ipv4Proto[p.opt.network], p.opt.source); err != nil {
		return
	}
	if err = p.conn4.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true); nil != err {
		return
	}
	if p.conn6, err = icmp.ListenPacket(ipv6Proto[p.opt.network], p.opt.source); err != nil {
		return
	}
	err = p.conn6.IPv6PacketConn().SetControlMessage(ipv6.FlagHopLimit, true)

	return
}

func (p *PingMan) waitTime() {
	defer p.opt.saveRecover()
	var t = time.NewTicker(p.opt.readWait)
	for {

		select {
		case _, ok := <-p.close:
			if !ok {
				goto END
			}
		case <-t.C:
			if !p.isNotice.Load() {
				p.wait <- struct{}{}
				p.isNotice.Store(true)
			}
			t.Reset(p.opt.readWait)
		default:
			if atomic.LoadInt64(&p.current) >= atomic.LoadInt64(&p.count) && !p.isNotice.Load() {
				p.wait <- struct{}{}
				p.isNotice.Store(true)
			}
		}
	}

END:
	t.Stop()
	p.wg.Done()

}
