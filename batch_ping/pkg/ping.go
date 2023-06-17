package pkg

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
	"net"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// Ping returns a new Ping struct pointer
// addr string, pid int, network string
func (p *Ping) Ping() error {
	if p.Addr == "" {
		return DstIpIsErr
	}
	if p.Id < 0 {
		return IdIsErr
	}

	if p.Network == "" {
		return NetworkIsErr
	}

	ipaddr, err := net.ResolveIPAddr("ip", p.Addr)
	if err != nil {
		return err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	p.ipaddr,
		p.AllRtt,
		p.dst,
		p.dstIP,
		p.tb,
		p.Size,
		p.Tracker,
		p.lock =
		ipaddr,
		make([]time.Duration, 0, p.Count),
		ipaddr,
		ipaddr.String(),
		make([]byte, 0, 8),
		timeSliceLength,
		r.Int63n(math.MaxInt64),
		&sync.RWMutex{}

	if isIPv4(ipaddr.IP) {
		p.ipv4, p.tyeIPV =
			true,
			ipv4.ICMPTypeEcho
	} else {
		p.tyeIPV = ipv6.ICMPTypeEchoRequest
	}

	if p.Network == Udp.String() {
		p.dst = &net.UDPAddr{IP: p.ipaddr.IP, Zone: p.ipaddr.Zone}
		p.dstIP = p.dst.String()
	}

	p.tb = append(timeToBytes(time.Now()), intToBytes(p.Tracker)...)
	if remainSize := p.Size - timeSliceLength - trackerLength; remainSize > 0 {
		p.tb = append(p.tb, bytes.Repeat([]byte{1}, remainSize)...)
	}

	return nil
}

// Ping represents ICMP packet sender/receiver
type Ping struct {
	dst            net.Addr
	tyeIPV         icmp.Type
	ipaddr         *net.IPAddr
	lock           *sync.RWMutex
	conn6          *icmp.PacketConn
	conn4          *icmp.PacketConn
	Addr           string
	dstIP          string
	Source         string
	Network        string
	tb             []byte
	AllRtt         []time.Duration
	size           int
	seqID          int64
	TTlIpV6        int
	Id             int
	TTlIpV4        int
	Tracker        int64
	PacketsReceive int64
	Size           int
	PacketsSent    int64
	Count          int64
	Interval       time.Duration
	ipv4           bool
}

type ResPing struct {
	ip         string
	seqID      int
	pid        int
	receivedAt time.Duration
}

type packet struct {
	addr   net.Addr
	proto  string
	bytes  []byte
	nBytes int
	ttl    int
}

// Packet represents a received and processed ICMP echo packet.
type Packet struct {
	IPAddr *net.IPAddr
	Addr   string
	Rtt    time.Duration
	NBytes int
	Seq    int
	Ttl    int
}

// Statistics represent the stats of a currently running or finished
// ping man operation.
type Statistics struct {
	IPAddr         *net.IPAddr
	Ip             string
	Addr           string
	AllRtt         []time.Duration
	PacketsReceive int64
	PacketsSent    int64
	PacketLoss     float64
	MinRtt         time.Duration
	MaxRtt         time.Duration
	AvgRtt         time.Duration
	StdDevRtt      time.Duration
}

// Statistics returns the statistics of the Ping. This can be run while the
// Ping is running  after it is finished. OnFinish calls this function to
// get it's finished statistics.
func (p *Ping) Statistics() *Statistics {
	p.lock.RLock()

	loss := float64(p.PacketsSent-p.PacketsReceive) / float64(p.PacketsSent) * 100
	var min, max, total time.Duration
	if len(p.AllRtt) > 0 {
		min = p.AllRtt[0]
		max = p.AllRtt[0]
	}
	defer p.lock.RUnlock()
	for _, rtt := range p.AllRtt {
		if rtt < min {
			min = rtt
		}
		if rtt > max {
			max = rtt
		}
		total += rtt
	}
	s := Statistics{
		PacketsSent:    p.PacketsSent,
		PacketsReceive: p.PacketsReceive,
		PacketLoss:     loss,
		AllRtt:         p.AllRtt,
		Addr:           p.Addr,
		IPAddr:         p.ipaddr,
		MaxRtt:         max,
		MinRtt:         min,
	}
	if len(p.Addr) > 0 {
		if len(p.AllRtt) > 0 {
			s.AvgRtt = total / time.Duration(len(p.AllRtt))
		}
		var sumSquares time.Duration
		for _, rtt := range p.AllRtt {
			sumSquares += (rtt - s.AvgRtt) * (rtt - s.AvgRtt)
		}
		if len(p.AllRtt) > 0 {
			s.StdDevRtt = time.Duration(math.Sqrt(
				float64(sumSquares / time.Duration(len(p.AllRtt)))))
		}

	}
	return &s
}

// SetConn setup connection
func (p *Ping) SetConn(conn4 *icmp.PacketConn, conn6 *icmp.PacketConn) {
	p.conn4 = conn4
	p.conn6 = conn6
}

func (p *Ping) SendICMP(seqID int) (err error) {
	var typ icmp.Type
	if p.ipv4 {
		typ = ipv4.ICMPTypeEcho
	} else {
		typ = ipv6.ICMPTypeEchoRequest
	}

	var dst net.Addr = p.ipaddr

	if p.Network == "udp" {
		dst = &net.UDPAddr{IP: p.ipaddr.IP, Zone: p.ipaddr.Zone}
	}

	t := append(timeToBytes(time.Now()), intToBytes(p.Tracker)...)
	if remainSize := p.Size - timeSliceLength - trackerLength; remainSize > 0 {
		t = append(t, bytes.Repeat([]byte{1}, remainSize)...)
	}

	body := &icmp.Echo{
		ID:   p.Id,
		Seq:  seqID,
		Data: t,
	}

	msg := &icmp.Message{
		Type: typ,
		Code: 0,
		Body: body,
	}

	if p.ipv4 {
		if p.TTlIpV4 > 0 {
			if err = ipv4.NewPacketConn(p.conn4).SetTTL(p.TTlIpV4); nil != err {
				return
			}
		}

	} else {
		if p.TTlIpV6 > 0 {
			if err = ipv6.NewPacketConn(p.conn6).SetHopLimit(p.TTlIpV6); nil != err {
				return
			}

		}
	}

	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return
	}

	for {
		if p.ipv4 {
			if _, err := p.conn4.WriteTo(msgBytes, dst); err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Err == syscall.ENOBUFS {
						continue
					}
				}
			}
		} else {
			if _, err := p.conn6.WriteTo(msgBytes, dst); err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Err == syscall.ENOBUFS {
						continue
					}
				}
			}
		}
		break
	}

	return
}

func bytesToTime(b []byte) time.Time {
	var nse int64
	for i := uint8(0); i < 8; i++ {
		nse += int64(b[i]) << ((7 - i) * 8)
	}
	return time.Unix(nse/1000000000, nse%1000000000)
}

func isIPv4(ip net.IP) bool {
	return len(ip.To4()) == net.IPv4len
}

func timeToBytes(t time.Time) []byte {
	nse := t.UnixNano()
	b := make([]byte, 8)
	for i := uint8(0); i < 8; i++ {
		b[i] = byte((nse >> ((7 - i) * 8)) & 0xff)
	}
	return b
}

func intToBytes(tracker int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(tracker))
	return b
}
