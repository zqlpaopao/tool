package pkg

import (
	"github.com/fatih/color"
	"golang.org/x/net/icmp"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type PingMan struct {
	err               error
	lock              *sync.Mutex
	wait              chan struct{}
	conn4             *icmp.PacketConn
	readyQueue        chan *Ping
	wg                *sync.WaitGroup
	initItems         map[string]*Ping
	conn6             *icmp.PacketConn
	initAddrMappingIp map[string]string
	resultQueue       chan *ResPing
	opt               *Option
	close             chan struct{}
	count             int64
	current           int64
	countCopy         int64
	id                int
	isNotice          atomic.Bool
}

func (p *PingMan) Do() *PingMan {
	if p.err = p.initConn(); nil != p.err {
		return p
	}

	p.RunCustomer()
	p.wg.Add(3)
	go p.receiveIpv4()
	go p.receiveIpv6()
	p.RunProducer()
	go p.waitTime()

	return p
}

// RunCustomer The consumption has been processed by the network card core
// and confirmed to be the data of the current package
func (p *PingMan) RunCustomer() {
	p.wg.Add(p.opt.customerN)
	for i := 0; i < p.opt.customerN; i++ {
		go p.Customer()
	}
}

func (p *PingMan) Customer() {
	defer p.opt.saveRecover()
LOOP:
	for {
		select {
		case _, ok := <-p.close:
			if !ok {
				break LOOP
			}

		case v, ok := <-p.resultQueue:
			if !ok {
				break LOOP
			}
			p.MakePing(v)
		}
	}
	p.wg.Done()
}

func (p *PingMan) MakePing(d *ResPing) {
	var (
		dst string
		ok  bool
	)
	p.lock.Lock()
	defer p.lock.Unlock()

	if dst, ok = p.initAddrMappingIp[d.ip]; !ok {
		return
	}
	p.initItems[dst].lock.Lock()
	p.initItems[dst].AllRtt = append(p.initItems[dst].AllRtt, d.receivedAt)
	p.initItems[dst].PacketsReceive++
	p.initItems[dst].lock.Unlock()

}

func (p *PingMan) RunProducer() {
	p.wg.Add(p.opt.readyPingN)
	for i := 0; i < p.opt.readyPingN; i++ {
		go p.Producer()
	}
}

func (p *PingMan) Producer() {
	defer p.opt.saveRecover()
LOOP:
	for {

		select {
		case _, ok := <-p.close:
			if !ok {
				break LOOP
			}

		case v, ok := <-p.readyQueue:
			if !ok {
				break LOOP
			}
			p.tidyPing(v)
		}
	}
	p.wg.Done()

}

func (p *PingMan) tidyPing(item *Ping) {
	seqID := 0
	item.SetConn(p.conn4, p.conn6)
	for i := 0; i < int(item.Count); i++ {
		if err := item.SendICMP(seqID); nil != err && p.opt.debug {
			color.Red("item.SendICMP  seqId-%d ping-%v err-%v \n", seqID, item, err)
			continue
		}
		if item.Interval > 0 {
			time.Sleep(item.Interval)
		}
		seqID++
	}
}

func (p *PingMan) Submit(pe *Ping) {

	p.lock.Lock()
	p.initItems[pe.Addr] = pe
	p.initItems[pe.Addr].lock.Lock()
	p.initItems[pe.Addr].PacketsSent++
	p.initItems[pe.Addr].lock.Unlock()

	p.initAddrMappingIp[pe.dstIP] = pe.Addr
	p.lock.Unlock()

	p.readyQueue <- pe
	atomic.AddInt64(&p.countCopy, pe.Count)

}

func (p *PingMan) Wait() {
	atomic.SwapInt64(&p.count, p.countCopy)
	<-p.wait
}

func (p *PingMan) GetRes() map[string]*Ping {
	return p.initItems
}

func (p *PingMan) GetError() error {
	return p.err
}

func (p *PingMan) Next() {
	atomic.SwapInt64(&p.count, math.MaxInt64)
	atomic.SwapInt64(&p.countCopy, 0)
	atomic.SwapInt64(&p.current, 0)

	p.isNotice.Store(false)
	p.lock.Lock()
	p.initItems = make(map[string]*Ping, p.opt.initPingMapSiz)
	p.initAddrMappingIp = make(map[string]string, p.opt.initPingMapSiz)
	p.lock.Unlock()

}

func (p *PingMan) Close() {
	close(p.readyQueue)
	close(p.resultQueue)
	close(p.close)
	//close(p.wait)
	p.wg.Wait()
	p.FreeConn()

}

func (p *PingMan) FreeConn() {
	_ = p.conn4.Close()
	_ = p.conn6.Close()
}
