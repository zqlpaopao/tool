package user_icmp

import (
	"syscall"
	"time"
)

// revIpv4 -- --------------------------
// --> @Describe revIpv4 rev the ipv4 packet
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) revIpv4(fd int) {
	defer p.option.Recover()
	var err error
	for {
		select {
		case <-p.done:
			goto END
		default:
			msg := p.receiveMMsgPool.Get()
			msg.N, msg.Dest, err = syscall.Recvfrom(fd, msg.Data, syscall.MSG_WAITALL)
			msg.RXTime,
				msg.V4 =
				time.Now(),
				true
			if err != nil {
				if err == syscall.EAGAIN ||
					err == syscall.EINTR {
					continue
				}
				e := p.errP.Get()
				e.Tag, e.Ping, e.Err = "revIpv4", nil, err
				p.errChan <- e
			}
			p.resChan <- msg
		}

	}
END:
	_ = syscall.Close(fd)
	p.wgRec.Done()
}

// revIpv6 -- --------------------------
// --> @Describe revIpv6 rev the ipv6 packet
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) revIpv6(fd int) {
	defer p.option.Recover()
	var err error
	for {
		select {
		case <-p.done:
			goto END
		default:
			msg := p.receiveMMsgPool.Get()
			msg.N, msg.Dest, err = syscall.Recvfrom(fd, msg.Data, syscall.MSG_WAITALL)
			msg.RXTime,
				msg.V4 =
				time.Now(),
				false
			if err != nil {
				if err == syscall.EAGAIN ||
					err == syscall.EINTR {
					continue
				}
				e := p.errP.Get()
				e.Tag, e.Ping, e.Err = "revIpv4", nil, err
				p.errChan <- e
			}
			p.resChan <- msg
		}

	}
END:
	_ = syscall.Close(fd)
	p.wgRec.Done()
}
