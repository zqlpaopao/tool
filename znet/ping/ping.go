package ping

import (
	"encoding/binary"
	"errors"
	"github.com/zqlpaopao/tool/string-byte/src"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
)

// New returns a new ping struct pointer.
func New(addr string) *Ping {
	return &Ping{
		Size:     timeSliceLength + trackerLength + 8,
		Addr:     addr,
		Pid:      0,
		Ipaddr:   nil,
		Ipv4:     false,
		Network:  IP,
		Protocol: Udp,
		TTL:      64,
		Uuid:     NewUUid().String(),
	}
}

// NewPing returns a new ping and resolves the address.
func NewPing(addr string) (*Ping, error) {
	p := New(addr)
	return p, p.Resolve()
}

// Ping represents a packet sender/receiver.
type Ping struct {
	Ipaddr     *net.IPAddr
	RevPacket  *Packet
	Uuid       string
	Addr       string
	Protocol   string
	Network    string
	Data       []byte
	Id         int64
	Pid        int
	Size       int
	Mark       uint
	Sequence   int
	Tracker    uint64
	TTL        int
	Ipv4       bool
	Df         bool
	IsSetMap   bool
	IsOrderSeq bool
}

// Packet represents a received and processed ICMP echo packet.
type Packet struct {
	StartTime time.Time
	EndTime   time.Time
	Addr      string
	uuid      string
	NBytes    int
	Seq       int
	ID        int
	//TTL       int
}

func (p *Packet) Uuid() string {
	return p.uuid
}

// SetIPAddrCheckIpv4 sets the ip address of the target host.
func (p *Ping) SetIPAddrCheckIpv4(ipaddr *net.IPAddr) *Ping {
	p.Ipv4 = isIPv4(ipaddr.IP)

	p.Ipaddr = ipaddr
	p.Addr = ipaddr.String()

	return p
}

// SetUUid returns the ip address of the target host.
func (p *Ping) SetUUid(uuid string) *Ping {
	p.Uuid = uuid
	return p
}

// Resolve does the DNS lookup for the Ping address and sets IP protocol.
func (p *Ping) Resolve() error {
	if len(p.Addr) == 0 {
		return errors.New("addr cannot be empty")
	}
	ipaddr, err := net.ResolveIPAddr(p.Network, p.Addr)
	if err != nil {
		return err
	}

	p.Ipv4 = isIPv4(ipaddr.IP)

	p.Ipaddr = ipaddr

	return nil
}

// SetAddr resolves and sets the ip address of the target host, addr can be a
// DNS name like "www.google.com" or IP like "127.0.0.1".
func (p *Ping) SetAddr(addr string) error {
	oldAddr := p.Addr
	p.Addr = addr
	err := p.Resolve()
	if err != nil {
		p.Addr = oldAddr
		return err
	}
	return nil
}

// SetDstAddr set packet number
func (p *Ping) SetDstAddr(addr string) *Ping {
	p.Addr = addr
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
	p.Ipv4 = true
	return p
}

// SetResolveIpAddr set the ip is Ipaddr
func (p *Ping) SetResolveIpAddr(ipaddr *net.IPAddr) *Ping {
	p.Ipaddr = ipaddr
	return p
}

// SetNetwork allows configuration of DNS resolution.
// * "ip" will automatically select IPv4 or IPv6.
// * "ip4" will select IPv4.
// * "ip6" will select IPv6.
func (p *Ping) SetNetwork(n string) *Ping {
	switch n {
	case "ip4":
		p.Network = "ip4"
	case "ip6":
		p.Network = "ip6"
	default:
		p.Network = "ip"
	}
	return p
}

// SetPrivileged sets the type of ping will send.
// false means ping will send an "unprivileged" UDP ping.
// true means ping will send a "privileged" raw ICMP ping.
// NOTE: setting to true requires that it be run with super-user privileges.
func (p *Ping) SetPrivileged(privileged bool) *Ping {
	if privileged {
		p.Protocol = "icmp"
	} else {
		p.Protocol = "udp"
	}
	return p
}

// SetMark sets a mark intended to be set on outgoing ICMP packets.
func (p *Ping) SetMark(m uint) *Ping {
	p.Mark = m
	return p
}

// SetPID sets the ICMP identifier.
func (p *Ping) SetPID(pid int) *Ping {
	p.Pid = pid
	return p
}

// SetOrderSeq sets the ICMP identifier.
func (p *Ping) SetOrderSeq(order bool) *Ping {
	p.IsOrderSeq = true
	return p
}

// SetSeq sets the ICMP identifier.
func (p *Ping) SetSeq(seq int) *Ping {
	p.Sequence = seq
	return p
}

// SetIcmp sets the ICMP identifier.
func (p *Ping) SetIcmp() *Ping {

	p.initTimeToBytes()

	return p
}

// SetDoNotFragment sets the do-not-fragment bit in the outer IP header to the desired value.
func (p *Ping) SetDoNotFragment(df bool) *Ping {
	p.Df = df
	return p
}

// SetID set the incr id.
func (p *Ping) SetID(id int64) *Ping {
	p.Id = id
	return p
}

// SetIsMap set the incr id.
func (p *Ping) SetIsMap(t bool) *Ping {
	p.IsSetMap = t
	return p
}

///////////////////////////////////////////////// get /////////////////////////////////////

// Privileged returns whether ping is running in privileged mode.
func (p *Ping) Privileged() bool {
	return p.Protocol == "icmp"
}

// GetAddr returns the string ip address of the target host.
func (p *Ping) GetAddr() string {
	return p.Addr
}

// IPAddr returns the ip address of the target host.
func (p *Ping) IPAddr() *net.IPAddr {
	return p.Ipaddr
}

// ID returns the ICMP identifier.
func (p *Ping) ID() int64 {
	return p.Id
}

// GetUuid returns the uuid to be set on outgoing ICMP packets.
func (p *Ping) GetUuid() string {
	return p.Uuid
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
		Addr:           ps.Addr,
		IPAddr:         ps.Ipaddr,
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
func (p *Pool) getPacketUUID(pkt []byte) string {
	return src.Bytes2String(pkt[timeSliceLength : timeSliceLength+trackerLength])
}

// bytesToTime -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
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

// initTimeToBytes -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Ping) initTimeToBytes() {
	var (
		uuidEncoded = src.String2Bytes(p.Uuid)[:]
		proto       byte
	)
	if p.Ipv4 {
		proto = byte(ipv4.ICMPTypeEcho)
	} else {
		proto = byte(ipv6.ICMPTypeEchoRequest)
	}

	//test
	p.Data = make([]byte, p.Size+8)
	p.Data[0],
		p.Data[1] =
		proto,
		byte(0)

	b := make([]byte, 2)

	binary.BigEndian.PutUint16(b[:2], uint16(p.Pid))
	p.Data[4],
		p.Data[5] =
		b[0],
		b[1]

	for i := timeSliceLength + 8; i < timeSliceLength+trackerLength+8; i++ {
		p.Data[i] = uuidEncoded[i-16]
	}
	if remainSize := p.Size + 8 - timeSliceLength - trackerLength; remainSize > 0 {
		for i := timeSliceLength + trackerLength + 8; i < p.Size+8; i++ {
			p.Data[i] = 1
		}
	}

	if proto == ProtocolIPv6ICMP {
		//if psh == nil { // cannot calculate checksum here
		//	return b, nil
		//}
		//off, l := 2*net.IPv6len, len(b)-len(psh)
		//binary.BigEndian.PutUint32(b[off:off+4], uint32(l))
	}

	//return b[len(psh):], nil
}

// ReplaceTimeToBytes -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Ping) ReplaceTimeToBytes() {
	var nSec = time.Now().UnixNano()
	for i := uint8(8); i < 16; i++ {
		p.Data[i] = byte((nSec >> ((15 - i) * 8)) & 0xff)
	}
}

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
