package ping

import (
	"errors"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
)

// New returns a new ping struct pointer.
func New(addr string) *Ping {
	uuid := NewUUid()
	return &Ping{
		Size:     timeSliceLength + trackerLength + 8,
		addr:     addr,
		id:       0,
		ipaddr:   nil,
		ipv4:     false,
		network:  IP,
		protocol: Udp,
		TTL:      64,
		uuid:     &uuid,
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

	uuid *UUID

	addr string

	// protocol is "icmp" or "udp".
	protocol string

	// network is one of "ip", "ip4", or "ip6".
	network string

	// Tracker: Used to uniquely identify packets - Deprecated
	Tracker uint64

	sequence int

	// mark is a SO_MARK (fw-mark) set on outgoing icmp packets
	mark uint

	TTL int
	id  int

	// Size of packet being sent
	Size int

	ipv4 bool

	// df when true sets the do-not-fragment bit in the outer IP or IPv6 header
	df bool
}

type packet struct {
	receivedAt time.Time
	addr       net.Addr
	bytes      []byte
	nBytes     int
	ttl        int
	ipv4       bool
}

// Packet represents a received and processed ICMP echo packet.
type Packet struct {
	StartTime time.Time
	EndTime   time.Time
	Addr      string
	uuid      string
	NBytes    int
	Seq       int
	TTL       int
	ID        int
}

func (p *Packet) Uuid() string {
	return p.uuid
}

// SetIPAddrCheckIpv4 sets the ip address of the target host.
func (p *Ping) SetIPAddrCheckIpv4(ipaddr *net.IPAddr) *Ping {
	p.ipv4 = isIPv4(ipaddr.IP)

	p.ipaddr = ipaddr
	p.addr = ipaddr.String()

	return p
}

// SetUUid returns the ip address of the target host.
func (p *Ping) SetUUid(uuid *UUID) *Ping {
	p.uuid = uuid
	return p
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

// SetDstAddr set packet number
func (p *Ping) SetDstAddr(addr string) *Ping {
	p.addr = addr
	return p
}

// SetSize set packet number
func (p *Ping) SetSize(size int) *Ping {
	p.Size = size
	return p
}

// SetTtl set packet ttl
func (p *Ping) SetTtl(ttl int) *Ping {
	p.TTL = ttl
	return p
}

// SetIpV4 set the ip is ipv4
func (p *Ping) SetIpV4() *Ping {
	p.ipv4 = true
	return p
}

// SetResolveIpAddr set the ip is ipaddr
func (p *Ping) SetResolveIpAddr(ipaddr *net.IPAddr) *Ping {
	p.ipaddr = ipaddr
	return p
}

// SetNetwork allows configuration of DNS resolution.
// * "ip" will automatically select IPv4 or IPv6.
// * "ip4" will select IPv4.
// * "ip6" will select IPv6.
func (p *Ping) SetNetwork(n string) *Ping {
	switch n {
	case "ip4":
		p.network = "ip4"
	case "ip6":
		p.network = "ip6"
	default:
		p.network = "ip"
	}
	return p
}

// SetPrivileged sets the type of ping will send.
// false means ping will send an "unprivileged" UDP ping.
// true means ping will send a "privileged" raw ICMP ping.
// NOTE: setting to true requires that it be run with super-user privileges.
func (p *Ping) SetPrivileged(privileged bool) *Ping {
	if privileged {
		p.protocol = "icmp"
	} else {
		p.protocol = "udp"
	}
	return p
}

// SetMark sets a mark intended to be set on outgoing ICMP packets.
func (p *Ping) SetMark(m uint) *Ping {
	p.mark = m
	return p
}

// SetID sets the ICMP identifier.
func (p *Ping) SetID(id int) *Ping {
	p.id = id
	return p
}

// SetDoNotFragment sets the do-not-fragment bit in the outer IP header to the desired value.
func (p *Ping) SetDoNotFragment(df bool) *Ping {
	p.df = df
	return p
}

///////////////////////////////////////////////// get /////////////////////////////////////

// Privileged returns whether ping is running in privileged mode.
func (p *Ping) Privileged() bool {
	return p.protocol == "icmp"
}

// Addr returns the string ip address of the target host.
func (p *Ping) Addr() string {
	return p.addr
}

// IPAddr returns the ip address of the target host.
func (p *Ping) IPAddr() *net.IPAddr {
	return p.ipaddr
}

// ID returns the ICMP identifier.
func (p *Ping) ID() int {
	return p.id
}

// Uuid returns the uuid to be set on outgoing ICMP packets.
func (p *Ping) Uuid() *UUID {
	return p.uuid
}

// Statistics /////////////////////////////////////////////// Statistics /////////////////////////////////////
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

// StatisticsLog returns the statistics of the ping. This can be run while the
// ping is running, or after it is finished. OnFinish calls this function to
// get it's finished statistics.
func StatisticsLog(ps *Ping, packet *Packet) *Statistics {
	loss, rtt := 0.00, packet.EndTime.Sub(packet.StartTime)
	s := Statistics{
		PacketsSent:    1,
		PacketsReceive: 1,
		PacketLoss:     loss,
		RttS:           []time.Duration{rtt},
		Addr:           ps.Addr(),
		IPAddr:         ps.ipaddr,
		MaxRtt:         rtt,
		MinRtt:         rtt,
		AvgRtt:         rtt,
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
func (p *Pool) getPacketUUID(pkt []byte) (packetUUID UUID, err error) {
	packetUUID = GetUUID()
	if err = packetUUID.UnmarshalBinary(pkt[timeSliceLength : timeSliceLength+trackerLength]); err != nil {
		return
	}
	return
}

// getCurrentTrackerUUID grabs the latest tracker UUID.
func (p *Ping) getCurrentTrackerUUID() *UUID {
	return p.uuid
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
