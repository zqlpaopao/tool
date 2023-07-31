package pkg

import (
	"time"
)

const (
	Udp = "udp"
	IP  = "ip"

	ReadyChanSize = 20
	OnRevChanSize = 40

	ErrChanSize = 20

	CurrentMapSize = 20

	CustomerSize = 6

	ResCustomerSize = 6

	TimeOut = 2 * time.Second
)

const (
	timeSliceLength  = 8
	trackerLength    = 8
	protocolICMP     = 1
	protocolIPv6ICMP = 58
	protoIpv4        = "ipv4"
	protoIpv6        = "ipv6"
)

var (
	ipv4Proto = map[string]string{"ip": "ip4:icmp", "udp": "udp4"}
	ipv6Proto = map[string]string{"ip": "ip6:ipv6-icmp", "udp": "udp6"}
)
