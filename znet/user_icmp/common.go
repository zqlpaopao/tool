package user_icmp

import (
	"golang.org/x/net/bpf"
	"time"
)

type Filter []bpf.Instruction

const (
	//Icmp Protocol = "user_icmp"

	//IcmpProtocol = 1

	IpHeader = 20

	//PrepareChLen is the object *Ping of ready to send
	PrepareChLen = 10

	//ErrorInfoSize is the pool error chan
	ErrorInfoSize = 3

	NoIcmpAndDarwin = 20

	timeSliceLength = 8
	Ttl2AndSeq2     = 4

	//trackerLength = len(uuid.UUID{})

	//IPv4ICMPTypeEchoReply = 0
	//IPv6ICMPTypeEchoReply = 129
)

var (
	//protocol    = map[Protocol]string{Icmp: "ip4:1", Udp: "ip4:17", Tcp: "ip4:6"}
	//protocolInt = map[Protocol]int{Icmp: 1, Udp: 17, Tcp: 6}
	seed = time.Now().UnixNano()
)
