package icmp

import (
	"encoding/binary"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"sync/atomic"
	"syscall"
	"time"
)

// Ping represents a packet sender/receiver.
type Ping struct {
	Data         []byte
	TTL          uint16
	Sequence     uint16
	Size         int
	pid          int
	Ipv4         bool
	IsOrderSeq   bool
	SocketAddrV4 *syscall.SockaddrInet4
	SocketAddrV6 *syscall.SockaddrInet6
}

// Packet represents a received and processed ICMP echo packet.
type Packet struct {
	TXTime time.Time
	RXTime time.Time
	Dest   net.IP
	NBytes int
	Seq    uint16
	Ttl    uint16
}

type ReceiveMMsg struct {
	RXTime time.Time
	Dest   syscall.Sockaddr
	N      int
	Data   []byte
	Seq    int
	Ttl    int
	V4     bool
}

// SetIpv6 returns the ip address of the target host.
func (p *Ping) SetIpv6() *Ping {
	p.Ipv4 = false
	return p
}

// SetSize set packet number
func (p *Ping) SetSize(size int) *Ping {
	if size < 32 {
		return p
	}
	p.Size = size
	return p
}

// SetTtl set packet ttl
func (p *Ping) SetTtl(ttl uint16) *Ping {
	p.TTL = ttl
	return p
}

// SetDestAddrV4 set the ip is Ipaddr
func (p *Ping) SetDestAddrV4(ipaddr *syscall.SockaddrInet4) *Ping {
	p.SocketAddrV4 = ipaddr
	return p
}

// SetDestAddrV6 set the ip is Ipaddr
func (p *Ping) SetDestAddrV6(ipaddr *syscall.SockaddrInet6) *Ping {
	p.SocketAddrV6 = ipaddr
	return p
}

// SetOrderSeq sets the ICMP identifier.
func (p *Ping) SetOrderSeq(order bool) *Ping {
	p.IsOrderSeq = order
	return p
}

// SetSeq sets the ICMP identifier.
func (p *Ping) SetSeq(seq uint16) *Ping {
	p.Sequence = seq
	return p
}

// SetIcmp sets the ICMP identifier.
func (p *Ping) SetIcmp() *Ping {
	p.initTimeToBytes()
	return p
}

// Statistics /////////////////////////////////////////////// Statistics /////////////////////////////////////
// Statistics represent the stats of a currently running or finished
// icmp operation.
type Statistics struct {
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

// StatisticsLog returns the statistics of the icmp. This can be run while the
// icmp is running, or after it is finished. OnFinish calls this function to
// get it's finished statistics.
func StatisticsLog(packet *Packet) *Statistics {
	loss, rtt := 0.00, packet.RXTime.Sub(packet.TXTime)
	s := Statistics{
		PacketsSent:    1,
		PacketsReceive: 1,
		PacketLoss:     loss,
		RttS:           []time.Duration{rtt},
		Addr:           packet.Dest.String(),
		MaxRtt:         rtt,
		MinRtt:         rtt,
		AvgRtt:         rtt,
	}
	return &s
}

// initTimeToBytes -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Ping) initTimeToBytes() {
	p.icmpMsg()
}

// icmpMsg -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Ping) icmpMsg() {
	if p.Ipv4 {
		p.icmpV4Message(int(ipv4.ICMPTypeEcho))
		return
	}
	p.icmpV4Message(int(ipv6.ICMPTypeEchoRequest))
}

// -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Ping) icmpV4Message(protocol int) {
	p.Data = make([]byte, p.Size+20)
	p.Data[0],
		p.Data[1] =
		byte(protocol),
		byte(0)

	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b[:2], uint16(p.pid))
	p.Data[4],
		p.Data[5] =
		b[0],
		b[1]

	binary.BigEndian.PutUint16(b[:2], p.TTL)
	p.Data[16],
		p.Data[17] =
		b[0],
		b[1]

	binary.BigEndian.PutUint16(b[:2], p.Sequence)
	p.Data[18],
		p.Data[19] =
		b[0],
		b[1]
	//var uuidEncoded = src.String2Bytes(p.Uuid)[:]
	//for i := timeSliceLength + 8; i < timeSliceLength+trackerLength+8; i++ {
	//	p.Data[i] = uuidEncoded[i-16]
	//}
	if remainSize := p.Size - timeSliceLength - Ttl2AndSeq2 - 8; remainSize > 0 {
		for i := timeSliceLength + Ttl2AndSeq2 + 8; i < p.Size; i++ {
			p.Data[i] = 1
		}
	}

}

// ReplaceTimeToBytes -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (p *Ping) ReplaceTimeToBytes() {
	var nSec = time.Now().UnixNano()

	//nSec = 1695289891422141000
	for i := uint8(8); i < 16; i++ {
		p.Data[i] = byte((nSec >> ((15 - i) * 8)) & 0xff)
	}
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

// getSeed returns a goroutine-safe unique seed
func getSeed() int64 {
	return atomic.AddInt64(&seed, 1)
}
