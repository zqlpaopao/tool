package pkg

import (
	"github.com/orcaman/concurrent-map/v2"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"sync"
	"time"
)

type ErrInfo struct {
	ping *PingItem
	err  error
}

type PingPool struct {
	dstAddrOfSrcAddr cmap.ConcurrentMap[string, string]
	pingMap          cmap.ConcurrentMap[string, *PingItem]
	err              error
	wg               *sync.WaitGroup
	errChan          chan *ErrInfo
	done             chan struct{}
	onRevChan        chan *Packet
	mapTidyLock      *sync.Mutex
	conn4            *icmp.PacketConn
	conn6            *icmp.PacketConn
	errWg            *sync.WaitGroup
	readyChan        chan *PingItem
	option           *Option
	readWg           *sync.WaitGroup
	revWg            *sync.WaitGroup
	current          int32
	sendTotal        int32
}

func (p *PingPool) Run() *PingPool {
	p.init()
	p.listen()
	p.loopErr()
	p.loopRev()
	p.rev()
	p.startWorker()
	go p.loopTidyMap()
	return p
}

func (p *PingPool) init() {
	cmap.SHARD_COUNT = p.option.currentMapSize
	p.pingMap, p.dstAddrOfSrcAddr = cmap.New[*PingItem](), cmap.New[string]()
}

func (p *PingPool) listen() {

	if p.conn4, p.err = icmp.ListenPacket(ipv4Proto[p.option.network], p.option.source); p.err != nil {
		return
	}
	if p.err = p.conn4.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true); nil != p.err {
		return
	}

	if p.conn6, p.err = icmp.ListenPacket(ipv6Proto[p.option.network], p.option.source); p.err != nil {
		return
	}
	p.err = p.conn6.IPv6PacketConn().SetControlMessage(ipv6.FlagHopLimit, true)
}

func (p *PingPool) loopTidyMap() {
	var t = time.NewTicker(p.option.mapTidyInterval)
	for {
		select {
		case <-t.C:
			p.tidyMap()
			t.Reset(p.option.mapTidyInterval)
		case _, ok := <-p.done:
			if !ok {
				goto END
			}
		}

	}
END:
	t.Stop()
}

func (p *PingPool) tidyMap() {
	its := p.pingMap.Items()
	itDst := p.dstAddrOfSrcAddr.Items()
	p.init()

	for k, v := range its {
		p.pingMap.Set(k, v)
	}
	for k, v := range itDst {
		p.dstAddrOfSrcAddr.Set(k, v)
	}
}

func (p *PingPool) startWorker() {
	if p.err != nil {
		return
	}
	p.wg.Add(p.option.customerSize)
	for i := 0; i < p.option.customerSize; i++ {
		go p.loop()
	}
}

func (p *PingPool) loop() {
LOOP:
	for {
		select {
		case v, ok := <-p.readyChan:
			if !ok {
				break LOOP
			}
			p.Ping(v)
		}
	}
	p.wg.Done()
}

func (p *PingPool) Ping(ps *PingItem) {
	ps.SetConn(p.conn4, p.conn6)
	for i := 1; i <= int(ps.Count); i++ {
		if err := ps.SendICMP(i); nil != err {
			p.errChan <- &ErrInfo{ping: ps, err: err}
			continue
		}
		time.Sleep(ps.Interval)
		//atomic.AddInt32(&p.sendTotal, 1)
	}
}

func (p *PingPool) Submit(addr string) (ps *PingItem, err error) {
	if p.err != nil {
		return nil, p.err
	}
	if ps, err = NewPingItem(addr, p.option.pid, p.option.network); err != nil {
		return
	}
	p.readyChan <- ps
	p.set(ps)

	return
}

func (p *PingPool) SubmitPingItem(ps *PingItem) (err error) {
	if p.err != nil {
		return p.err
	}
	ps.id,
		ps.network =
		p.option.pid,
		p.option.network

	p.readyChan <- ps
	p.set(ps)
	return
}

func (p *PingPool) set(ps *PingItem) {
	p.dstAddrOfSrcAddr.Set(ps.IPAddr().String(), ps.Addr)
	p.pingMap.Set(ps.Addr, ps)

}

func (p *PingPool) loopErr() {
	p.errWg.Add(1)
	for i := 0; i < 1; i++ {
		go p.callBackErrFunc()
	}

}
func (p *PingPool) callBackErrFunc() {
	for {
		select {
		case v, ok := <-p.errChan:
			if !ok {
				goto END
			}
			p.option.errCallBack(v.ping, v.err)
		}
	}
END:
	p.errWg.Done()
}

func (p *PingPool) loopRev() {
	p.revWg.Add(p.option.onRevWorkerNum)
	for i := 0; i < p.option.onRevWorkerNum; i++ {
		go p.revCallBack()
	}
}

func (p *PingPool) revCallBack() {
	for {
		select {
		case v, ok := <-p.onRevChan:
			if !ok {
				goto END
			}
			p.callBack(v)
		}
	}
END:
	p.revWg.Done()
}

func (p *PingPool) callBack(pkt *Packet) {
	var (
		ip   string
		ping *PingItem
		ok   bool
	)
	if ip, ok = p.dstAddrOfSrcAddr.Get(pkt.IPAddr); !ok {
		return
	}

	if ping, ok = p.pingMap.Get(ip); !ok {
		return
	}
	p.option.OnRevFunc(ping, pkt)
}

func (p *PingPool) rev() {
	if p.err != nil {
		return
	}
	p.readWg.Add(2)
	go p.revIpv4()
	go p.revIpv6()
}

func (p *PingPool) Error() error {
	return p.err
}

func (p *PingPool) Close() *PingPool {
	close(p.readyChan)
	p.wg.Wait()
	return p
}

func (p *PingPool) CloseRev() {
	close(p.done)
	p.readWg.Wait()
	close(p.onRevChan)
	p.revWg.Wait()
	close(p.errChan)
	p.errWg.Wait()
	p.CloseConn()
}

func (p *PingPool) CloseConn() {
	_ = p.conn4.Close()
	_ = p.conn6.Close()
}
