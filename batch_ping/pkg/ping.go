package pkg

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
	"net"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// NewPingItem returns a new PingItem struct pointer
func NewPingItem(addr string, pid int, network string) (*PingItem, error) {
	ipaddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return nil, err
	}

	var ipv4Tag bool
	if isIPv4(ipaddr.IP) {
		ipv4Tag = true
	} else if isIPv6(ipaddr.IP) {
		ipv4Tag = false
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &PingItem{
		ipaddr:   ipaddr,
		Addr:     addr,
		Interval: time.Second,
		Timeout:  time.Second * 100000,
		Count:    1,
		id:       pid,
		network:  network,
		ipv4:     ipv4Tag,
		Size:     timeSliceLength,
		Tracker:  r.Int63n(math.MaxInt64),
	}, nil
}

// PingWith returns a new PingItem struct pointer
func (p *PingItem) PingWith() (err error) {
	ipaddr, err := net.ResolveIPAddr("ip", p.Addr)
	if err != nil {
		return
	}

	if isIPv4(ipaddr.IP) {
		p.ipv4 = true
	} else if isIPv6(ipaddr.IP) {
		p.ipv4 = false
	}

	p.ipaddr =
		ipaddr
	return
}

// PingItem represents ICMP Packet sender/receiver
type PingItem struct {
	conn4       *icmp.PacketConn
	ipaddr      *net.IPAddr
	conn6       *icmp.PacketConn
	Source      string
	network     string
	Addr        string
	RttS        []time.Duration
	Tracker     int64
	Size        int
	PacketsRev  int64
	size        int
	id          int
	PacketsSent int64
	Count       int64
	Timeout     time.Duration
	Interval    time.Duration
	ipv4        bool
}

type packet struct {
	addr   net.Addr
	proto  string
	bytes  []byte
	nBytes int
	ttl    int
}

// Packet represents a received and processed ICMP echo Packet.
type Packet struct {
	StartTime time.Time
	EndTime   time.Time
	IPAddr    string
	Addr      string
	Bytes     []byte
	Rtt       time.Duration
	NBytes    int
	Seq       int
	Ttl       int
	ID        int
}

// Statistics represent the stats of a currently running or finished
// PingItem operation.
type Statistics struct {
	IPAddr      *net.IPAddr
	Addr        string
	RttS        []time.Duration
	PacketsRev  int64
	PacketsSent int64
	PacketLoss  float64
	MinRtt      time.Duration
	MaxRtt      time.Duration
	AvgRtt      time.Duration
	StdDevRtt   time.Duration
}

// SetConn set ipv4 and ipv6 conn
func (p *PingItem) SetConn(conn4 *icmp.PacketConn, conn6 *icmp.PacketConn) {
	p.conn4 = conn4
	p.conn6 = conn6
}

// SetIPAddr sets the ip address of the target host.
func (p *PingItem) SetIPAddr(ipaddr *net.IPAddr) {
	var ipv4Tag bool
	if isIPv4(ipaddr.IP) {
		ipv4Tag = true
	} else if isIPv6(ipaddr.IP) {
		ipv4Tag = false
	}

	p.ipaddr = ipaddr
	p.Addr = ipaddr.String()
	p.ipv4 = ipv4Tag
}

// IPAddr returns the ip address of the target host.
func (p *PingItem) IPAddr() *net.IPAddr {
	return p.ipaddr
}

// ID returns the id
func (p *PingItem) ID() int {
	return p.id
}

// SetAddr resolves and sets the ip address of the target host, Addr can be a
// DNS name like "www.google.com" or IP like "127.0.0.1".
func (p *PingItem) SetAddr(addr string) error {
	ipaddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return err
	}

	p.SetIPAddr(ipaddr)
	p.Addr = addr
	return nil
}

//// Addr returns the string ip address of the target host.
//func (p *PingItem) Addr() string {
//	return p.Addr
//}

// SetPrivileged sets the type of ping PingItem will send.
// false means PingItem will send an "unprivileged" UDP ping.
// true means PingItem will send a "privileged" raw ICMP ping.
// NOTE: setting to true requires that it be run with super-user privileges.
func (p *PingItem) SetPrivileged(privileged bool) {
	if privileged {
		p.network = "ip"
	} else {
		p.network = "udp"
	}
}

// Privileged returns whether PingItem is running in privileged mode.
func (p *PingItem) Privileged() bool {
	return p.network == "ip"
}

// Statistics returns the statistics of the PingItem. This can be run while the
// PingItem is running, or after it is finished. OnFinish calls this function to
// get it's finished statistics.
func (p *PingItem) Statistics() *Statistics {
	loss := float64(p.PacketsSent-p.PacketsRev) / float64(p.PacketsSent) * 100
	var min, max, total time.Duration
	if len(p.RttS) > 0 {
		min = p.RttS[0]
		max = p.RttS[0]
	}
	for _, rtt := range p.RttS {
		if rtt < min {
			min = rtt
		}
		if rtt > max {
			max = rtt
		}
		total += rtt
	}
	s := Statistics{
		PacketsSent: p.PacketsSent,
		PacketsRev:  p.PacketsRev,
		PacketLoss:  loss,
		RttS:        p.RttS,
		Addr:        p.Addr,
		IPAddr:      p.ipaddr,
		MaxRtt:      max,
		MinRtt:      min,
	}
	if len(p.RttS) > 0 {
		s.AvgRtt = total / time.Duration(len(p.RttS))
		var squares time.Duration
		for _, rtt := range p.RttS {
			squares += (rtt - s.AvgRtt) * (rtt - s.AvgRtt)
		}
		s.StdDevRtt = time.Duration(math.Sqrt(
			float64(squares / time.Duration(len(p.RttS)))))
	}
	return &s
}

func (p *PingItem) SendICMP(seqID int) (err error) {
	var (
		typ      icmp.Type
		dst      net.Addr = p.ipaddr
		msgBytes []byte
	)
	if p.ipv4 {
		typ = ipv4.ICMPTypeEcho
	} else {
		typ = ipv6.ICMPTypeEchoRequest
	}

	if p.network == "udp" {
		dst = &net.UDPAddr{IP: p.ipaddr.IP, Zone: p.ipaddr.Zone}
	}

	t := append(timeToBytes(time.Now()), IntToBytes(p.Tracker)...)
	if remainSize := p.Size - timeSliceLength - trackerLength; remainSize > 0 {
		t = append(t, bytes.Repeat([]byte{1}, remainSize)...)
	}
	body := &icmp.Echo{
		ID:   p.id,
		Seq:  seqID,
		Data: t,
	}
	msg := &icmp.Message{
		Type: typ,
		Code: 0,
		Body: body,
	}

	if msgBytes,
		err = msg.Marshal(nil); err != nil {
		return
	}

	for {
		if p.ipv4 {
			if _, err = p.conn4.WriteTo(msgBytes, dst); err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Err == syscall.ENOBUFS {
						continue
					}
				}
			}
		} else {
			if _, err = p.conn6.WriteTo(msgBytes, dst); err != nil {
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
func (p *PingItem) ResetRttS() {
	p.RttS,
		p.PacketsRev,
		p.PacketsSent =
		make([]time.Duration, 0, p.Count),
		0,
		0
}

func bytesToTime(b []byte) time.Time {
	var nSecond int64
	for i := uint8(0); i < 8; i++ {
		nSecond += int64(b[i]) << ((7 - i) * 8)
	}
	return time.Unix(nSecond/1000000000, nSecond%1000000000)
}

func isIPv4(ip net.IP) bool {
	return len(ip.To4()) == net.IPv4len
}

func isIPv6(ip net.IP) bool {
	return len(ip) == net.IPv6len
}

func timeToBytes(t time.Time) []byte {
	nSecond := t.UnixNano()
	b := make([]byte, 8)
	for i := uint8(0); i < 8; i++ {
		b[i] = byte((nSecond >> ((7 - i) * 8)) & 0xff)
	}
	return b
}

func BytesToInt(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

func IntToBytes(tracker int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(tracker))
	return b
}
