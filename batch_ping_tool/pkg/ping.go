package ping

import (
	"errors"
	"github.com/google/uuid"
	"github.com/zqlpaopao/tool/rand-string/pkg"
	"math"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// New returns a new ping struct pointer.
func New(addr string) *Ping {
	r := rand.New(rand.NewSource(getSeed()))
	firstUUID := uuid.New()
	var firstSequence = map[uuid.UUID]map[int]struct{}{}
	firstSequence[firstUUID] = make(map[int]struct{})
	return &Ping{
		Count:        -1,
		Interval:     time.Second,
		Size:         timeSliceLength + trackerLength,
		Timeout:      time.Duration(math.MaxInt64),
		addr:         addr,
		id:           r.Intn(math.MaxUint16),
		trackerUUIDs: []uuid.UUID{firstUUID},
		ipaddr:       nil,
		ipv4:         false,
		network:      IP,
		protocol:     Udp,
		TTL:          64,
		uuid:         pkg.RandGenString(pkg.RandSourceLetterAndNumber, 5),
	}
}

// NewPing returns a new ping and resolves the address.
func NewPing(addr string) (*Ping, error) {
	p := New(addr)
	return p, p.Resolve()
}

// Ping represents a packet sender/receiver.
type Ping struct {
	ipaddr *net.IPAddr
	uuid   string

	// protocol is "icmp" or "udp".
	protocol string

	// network is one of "ip", "ip4", or "ip6".
	network string

	addr string

	// RttS is all the Rtt
	RttS []time.Duration

	// trackerUUIDs is the list of UUIDs being used for sending packets.
	trackerUUIDs []uuid.UUID

	// Round trip time statistics
	maxRtt time.Duration

	sequence int

	// Round trip time statistics
	stdDevRtt time.Duration

	// Size of packet being sent
	Size int

	// Tracker: Used to uniquely identify packets - Deprecated
	Tracker uint64

	// Round trip time statistics
	avgRtt time.Duration

	// Interval is the wait time between each packet send. Default is 1s.
	Interval time.Duration

	// Round trip time statistics
	minRtt time.Duration

	// mark is a SO_MARK (fw-mark) set on outgoing icmp packets
	mark uint

	// Timeout specifies a timeout before ping exits, regardless of how many
	// packets have been received.
	Timeout time.Duration

	// Number of packets received
	PacketsReceive int64

	TTL int
	id  int

	// Round trip time statistics
	stDDEvm2 time.Duration

	// Number of packets sent
	PacketsSent int64

	// Count tells ping to stop after sending (and receiving) Count echo
	// packets. If this option is not specified, ping will operate until
	// interrupted.
	Count int64

	rttLock sync.RWMutex
	ipv4    bool

	// df when true sets the do-not-fragment bit in the outer IP or IPv6 header
	df bool
}

type packet struct {
	addr   net.Addr
	bytes  []byte
	nBytes int
	ttl    int
	ipv4   bool
}

// Packet represents a received and processed ICMP echo packet.
type Packet struct {
	StartTime time.Time
	EndTime   time.Time
	IPAddr    *net.IPAddr
	Addr      string
	uuid      string
	Rtt       time.Duration
	NBytes    int
	Seq       int
	TTL       int
	ID        int
}

// Statistics represent the stats of a currently running or finished
// ping operation.
type Statistics struct {
	IPAddr         *net.IPAddr
	Addr           string
	RttS           []time.Duration
	PacketsReceive int
	PacketsSent    int
	PacketLoss     float64
	MinRtt         time.Duration
	MaxRtt         time.Duration
	AvgRtt         time.Duration
	StdDevRtt      time.Duration
}

// SetIPAddr sets the ip address of the target host.
func (p *Ping) SetIPAddr(ipaddr *net.IPAddr) {
	p.ipv4 = isIPv4(ipaddr.IP)

	p.ipaddr = ipaddr
	p.addr = ipaddr.String()
}

// IPAddr returns the ip address of the target host.
func (p *Ping) IPAddr() *net.IPAddr {
	return p.ipaddr
}

// Resolve does the DNS lookup for the Ping address and sets IP protocol.
func (p *Ping) Resolve() error {
	if len(p.addr) == 0 {
		return errors.New("addr cannot be empty")
	}
	ipaddr, err := net.ResolveIPAddr(p.network, p.addr)
	if err != nil {
		return err
	}

	p.ipv4 = isIPv4(ipaddr.IP)

	p.ipaddr = ipaddr

	return nil
}

// SetAddr resolves and sets the ip address of the target host, addr can be a
// DNS name like "www.google.com" or IP like "127.0.0.1".
func (p *Ping) SetAddr(addr string) error {
	oldAddr := p.addr
	p.addr = addr
	err := p.Resolve()
	if err != nil {
		p.addr = oldAddr
		return err
	}
	return nil
}

// Addr returns the string ip address of the target host.
func (p *Ping) Addr() string {
	return p.addr
}

// SetNetwork allows configuration of DNS resolution.
// * "ip" will automatically select IPv4 or IPv6.
// * "ip4" will select IPv4.
// * "ip6" will select IPv6.
func (p *Ping) SetNetwork(n string) {
	switch n {
	case "ip4":
		p.network = "ip4"
	case "ip6":
		p.network = "ip6"
	default:
		p.network = "ip"
	}
}

// SetPrivileged sets the type of ping will send.
// false means ping will send an "unprivileged" UDP ping.
// true means ping will send a "privileged" raw ICMP ping.
// NOTE: setting to true requires that it be run with super-user privileges.
func (p *Ping) SetPrivileged(privileged bool) {
	if privileged {
		p.protocol = "icmp"
	} else {
		p.protocol = "udp"
	}
}

// Privileged returns whether ping is running in privileged mode.
func (p *Ping) Privileged() bool {
	return p.protocol == "icmp"
}

// SetID sets the ICMP identifier.
func (p *Ping) SetID(id int) {
	p.id = id
}

// ID returns the ICMP identifier.
func (p *Ping) ID() int {
	return p.id
}

// SetMark sets a mark intended to be set on outgoing ICMP packets.
func (p *Ping) SetMark(m uint) {
	p.mark = m
}

// Mark returns the mark to be set on outgoing ICMP packets.
func (p *Ping) Mark() uint {
	return p.mark
}

// SetDoNotFragment sets the do-not-fragment bit in the outer IP header to the desired value.
func (p *Ping) SetDoNotFragment(df bool) {
	p.df = df
}

// StatisticsLog returns the statistics of the ping. This can be run while the
// ping is running, or after it is finished. OnFinish calls this function to
// get it's finished statistics.
func StatisticsLog(ps *Ping) *Statistics {
	pn, rn := atomic.LoadInt64(&ps.PacketsSent), atomic.LoadInt64(&ps.PacketsReceive)
	loss := float64(pn-rn) / float64(pn) * 100
	var min, max, total time.Duration
	ps.rttLock.RLock()
	defer ps.rttLock.RUnlock()
	if len(ps.RttS) > 0 {
		min = ps.RttS[0]
		max = ps.RttS[0]
	}
	for _, rtt := range ps.RttS {
		if rtt < min {
			min = rtt
		}
		if rtt > max {
			max = rtt
		}
		total += rtt
	}
	s := Statistics{
		PacketsSent:    int(pn),
		PacketsReceive: int(rn),
		PacketLoss:     loss,
		RttS:           ps.RttS,
		Addr:           ps.Addr(),
		IPAddr:         ps.ipaddr,
		MaxRtt:         max,
		MinRtt:         min,
	}
	if len(ps.RttS) > 0 {
		s.AvgRtt = total / time.Duration(len(ps.RttS))
		var squares time.Duration
		for _, rtt := range ps.RttS {
			squares += (rtt - s.AvgRtt) * (rtt - s.AvgRtt)
		}
		s.StdDevRtt = time.Duration(math.Sqrt(
			float64(squares / time.Duration(len(ps.RttS)))))
	}
	return &s
}

type expBackoff struct {
	baseDelay time.Duration
	maxExp    int64
	c         int64
}

func (b *expBackoff) Get() time.Duration {
	if b.c < b.maxExp {
		b.c++
	}

	return b.baseDelay * time.Duration(rand.Int63n(1<<b.c))
}

func newExpBackoff(baseDelay time.Duration, maxExp int64) expBackoff {
	return expBackoff{baseDelay: baseDelay, maxExp: maxExp}
}

// getPacketUUID scans the tracking slice for matches.
func (p *Pool) getPacketUUID(pkt []byte) (packetUUID uuid.UUID, err error) {
	if err = packetUUID.UnmarshalBinary(pkt[timeSliceLength : timeSliceLength+trackerLength]); err != nil {
		return
	}
	return
}

// getCurrentTrackerUUID grabs the latest tracker UUID.
func (p *Ping) getCurrentTrackerUUID() uuid.UUID {
	return p.trackerUUIDs[len(p.trackerUUIDs)-1]
}

func bytesToTime(b []byte) time.Time {
	var nSec int64
	for i := uint8(0); i < 8; i++ {
		nSec += int64(b[i]) << ((7 - i) * 8)
	}
	return time.Unix(nSec/1000000000, nSec%1000000000)
}

func isIPv4(ip net.IP) bool {
	return len(ip.To4()) == net.IPv4len
}

func timeToBytes(t time.Time) []byte {
	nSec := t.UnixNano()
	b := make([]byte, 8)
	for i := uint8(0); i < 8; i++ {
		b[i] = byte((nSec >> ((7 - i) * 8)) & 0xff)
	}
	return b
}

var seed = time.Now().UnixNano()

// getSeed returns a goroutine-safe unique seed
func getSeed() int64 {
	return atomic.AddInt64(&seed, 1)
}

// stripIPv4Header strips IPv4 header bytes if present
// https://github.com/golang/go/commit/3b5be4522a21df8ce52a06a0c4ba005c89a8590f
func stripIPv4Header(n int, b []byte) int {
	if len(b) < 20 {
		return n
	}
	l := int(b[0]&0x0f) << 2
	if 20 > l || l > len(b) {
		return n
	}
	if b[0]>>4 != 4 {
		return n
	}
	copy(b, b[l:])
	return n - l
}

func (p *Ping) Reset() {
	atomic.SwapInt64(&p.PacketsSent, 0)
	atomic.SwapInt64(&p.PacketsReceive, 0)
	p.rttLock.Lock()
	p.RttS = make([]time.Duration, 0, p.Count)
	p.rttLock.Unlock()
}

func (p *Ping) SwapPacketsSent(new int64) {
	atomic.SwapInt64(&p.PacketsSent, new)

}

func (p *Ping) AppendRtt(sub time.Duration) {
	p.rttLock.Lock()
	p.RttS = append(p.RttS, sub)
	p.rttLock.Unlock()
}
