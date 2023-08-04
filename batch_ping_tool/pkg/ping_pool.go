package ping

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/orcaman/concurrent-map/v2"
	"golang.org/x/net/icmp"
	"net"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type ErrInfo struct {
	ping *Ping
	err  error
}

type Pool struct {
	option      *Option
	conn4       chan packetConn
	conn4Rev    packetConn
	conn6Rev    packetConn
	conn6       chan packetConn
	pingMap     cmap.ConcurrentMap[string, *Ping]
	err         error
	errChan     chan *ErrInfo
	done        chan struct{}
	onRevChan   chan *Packet
	mapTidyLock *sync.Mutex
	wg          *sync.WaitGroup
	errWg       *sync.WaitGroup
	readyChan   chan *Ping
	readWg      *sync.WaitGroup
	revWg       *sync.WaitGroup
}

func (p *Pool) Run() *Pool {
	defer p.option.Recover()
	p.init()
	p.listen()
	p.loopErr()
	p.loopRev()
	p.rev()
	p.startWorker()
	go p.loopTidyMap()
	return p
}

func (p *Pool) init() {
	cmap.SHARD_COUNT = p.option.currentMapSize
	p.pingMap = cmap.New[*Ping]()
}

func (p *Pool) listen() {
	if p.conn4Rev, p.err = p.makeConn4(); p.err != nil {
		return
	}

	if p.conn6Rev, p.err = p.makeConn6(); p.err != nil {
		return
	}

	for i := 0; i < p.option.conn4Size; i++ {
		if conn4, err := p.makeConn4(); err == nil {
			p.conn4 <- conn4
		}
	}

	for i := 0; i < p.option.conn6Size; i++ {
		if conn6, err := p.makeConn6(); err == nil {
			p.conn6 <- conn6
		}
	}
}

func (p *Pool) loopTidyMap() {
	defer p.option.Recover()
	var t = time.NewTicker(p.option.mapTidyInterval)
	for {
		select {
		case <-t.C:
			p.tidyMap()
			t.Reset(p.option.mapTidyInterval)
		case <-p.done:
			goto END
		}

	}
END:
	t.Stop()
}

func (p *Pool) tidyMap() {
	p.mapTidyLock.Lock()
	defer p.mapTidyLock.Unlock()
	its := p.pingMap.Items()
	p.init()

	for k, v := range its {
		p.pingMap.Set(k, v)
	}

}

func (p *Pool) startWorker() {
	if p.err != nil {
		return
	}
	p.wg.Add(p.option.customerSize)
	for i := 0; i < p.option.customerSize; i++ {
		go p.loop()
	}
}

func (p *Pool) loop() {
	defer p.option.Recover()
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

func (p *Pool) Ping(ps *Ping) {
	var (
		conn     packetConn
		err      error
		timeout  *time.Ticker
		interval *time.Ticker
	)
	if conn, err = p.getConn(ps.ipv4); conn == nil || err != nil {
		p.errChan <- &ErrInfo{
			ping: ps,
			err:  err,
		}
		return
	}
	if err = p.sendICMP(ps, conn); err != nil {
		goto END
	}
	if ps.Count == 1 {
		return
	}

	timeout = time.NewTicker(ps.Timeout)
	interval = time.NewTicker(ps.Interval)
	defer func() {
		interval.Stop()
		timeout.Stop()
	}()

	for {
		select {
		case <-p.done:
			return

		case <-timeout.C:
			return

		case <-interval.C:
			if ps.Count > 0 && atomic.LoadInt64(&ps.PacketsSent) >= ps.Count {
				interval.Stop()
				continue
			}
			err = p.sendICMP(ps, conn)
			if err != nil {
				p.errChan <- &ErrInfo{ping: ps, err: err}
			}
		}
		if ps.Count > 0 && atomic.LoadInt64(&ps.PacketsSent) >= ps.Count {
			return
		}
	}
END:
	if err != nil {
		p.errChan <- &ErrInfo{ping: ps, err: err}
	}
}

func (p *Pool) getConn(ipv4 bool) (packetConn, error) {
	if ipv4 {
		return p.getConn4()
	}
	return p.getConn6()
}

func (p *Pool) Submit(ps *Ping) (err error) {
	if p.err != nil {
		return p.err
	}
	if ps.Size < timeSliceLength+trackerLength {
		return fmt.Errorf("size %d is less than minimum required size %d", ps.Size, timeSliceLength+trackerLength)
	}
	ps.id,
		ps.network,
		ps.protocol, ps.Size =
		p.option.pid,
		p.option.network,
		p.option.protocol,
		p.option.dataSize

	p.readyChan <- ps
	return
}

func (p *Pool) set(uuid string, ps *Ping) {
	p.pingMap.Set(uuid, ps)
}

func (p *Pool) get(uuid string) (ping *Ping, ok bool) {
	return p.pingMap.Get(uuid)
}

func (p *Pool) remove(uuid string) {
	p.pingMap.Remove(uuid)
}

func (p *Pool) loopErr() {
	p.errWg.Add(1)
	for i := 0; i < 1; i++ {
		go p.callBackErrFunc()
	}

}
func (p *Pool) callBackErrFunc() {
	defer p.option.Recover()
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

func (p *Pool) loopRev() {
	p.revWg.Add(p.option.onRevWorkerNum)
	for i := 0; i < p.option.onRevWorkerNum; i++ {
		go p.revCallBack()
	}
}

func (p *Pool) revCallBack() {
	defer p.option.Recover()
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

func (p *Pool) callBack(pkt *Packet) {
	var (
		ping *Ping
		ok   bool
	)
	if ping, ok = p.get(pkt.uuid); ping == nil || !ok {
		return
	}
	p.option.OnRevFunc(ping, pkt)
}

func (p *Pool) rev() {
	if p.err != nil {
		return
	}
	p.readWg.Add(2)
	go p.revIpv4()
	go p.revIpv6()
}

func (p *Pool) Error() error {
	return p.err
}

func (p *Pool) Close() *Pool {
	close(p.readyChan)
	p.wg.Wait()
	return p
}

func (p *Pool) CloseRev() {
	close(p.done)
	p.readWg.Wait()
	close(p.onRevChan)
	p.revWg.Wait()
	close(p.errChan)
	p.errWg.Wait()
	close(p.conn6)
	close(p.conn4)
	p.CloseConn()
}

func (p *Pool) CloseConn() {
	for conn := range p.conn4 {
		_ = conn.Close()
	}
	for conn := range p.conn6 {
		_ = conn.Close()
	}

}

// sendICMP send the packet
func (p *Pool) sendICMP(ps *Ping, conn packetConn) (err error) {

	if ps.mark != 0 {
		if err = conn.SetMark(ps.mark); err != nil {
			return fmt.Errorf("error setting mark: %v", err)
		}
	}

	if ps.df {
		if err = conn.SetDoNotFragment(); err != nil {
			return fmt.Errorf("error setting do-not-fragment: %v", err)
		}
	}

	conn.SetTTL(ps.TTL)

	if err = conn.SetFlagTTL(); nil != err {
		return err
	}

	var dst net.Addr = ps.ipaddr
	if ps.protocol == Udp {
		dst = &net.UDPAddr{IP: ps.ipaddr.IP, Zone: ps.ipaddr.Zone}
	}
	currentUUID := ps.getCurrentTrackerUUID()
	uuidEncoded, err := currentUUID.MarshalBinary()
	if err != nil {
		return fmt.Errorf("unable to marshal UUID binary: %w", err)
	}

	t := append(timeToBytes(time.Now()), uuidEncoded...)
	if remainSize := ps.Size - timeSliceLength - trackerLength; remainSize > 0 {
		t = append(t, bytes.Repeat([]byte{1}, remainSize)...)
	}
	body := &icmp.Echo{
		//ID:   p.id & 0xffff,
		ID:   ps.id,
		Seq:  ps.sequence,
		Data: t,
	}
	p.set(currentUUID.String(), ps)

	msg := &icmp.Message{
		Type: conn.ICMPRequestType(),
		Code: 0,
		Body: body,
	}
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return err
	}
	for {
		if _, err = conn.WriteTo(msgBytes, dst); err != nil {
			p.errChan <- &ErrInfo{ping: ps, err: err}
			if netErr, ok := err.(*net.OpError); ok {
				if netErr.Err == syscall.ENOBUFS {
					continue
				}
			}
			return err
		}
		// mark this sequence as in-flight
		//p.awaitingSequences[currentUUID][p.sequence] = struct{}{}
		atomic.AddInt64(&ps.PacketsSent, 1)
		ps.sequence++
		if ps.sequence > 65535 {
			newUUID := uuid.New()
			ps.trackerUUIDs = append(ps.trackerUUIDs, newUUID)
			ps.sequence = 0
		}
		break
	}
	return nil
}
