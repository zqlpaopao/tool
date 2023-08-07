package ping

import (
	"errors"
	"github.com/google/uuid"
)

const (
	Udp = "udp"
	IP  = "ip"

	ReadyChanSize = 20
	OnRevChanSize = 40

	ErrChanSize = 20

	CurrentMapSize = 20

	CustomerSize = 6

	ErrorInfoSize = 6

	ResCustomerSize = 6

	Conn4 = 4
	Conn6 = 6

	NoIcmpAndDarwin = 20
)

const (
	timeSliceLength  = 8
	trackerLength    = len(uuid.UUID{})
	protocolICMP     = 1
	protocolIPv6ICMP = 58
)

var (
	ipv4Proto = map[string]string{"icmp": "ip4:icmp", "udp": "udp4"}
	ipv6Proto = map[string]string{"icmp": "ip6:ipv6-icmp", "udp": "udp6"}

	ErrMarkNotSupported = errors.New("setting SO_MARK socket option is not supported on this platform")
	ErrDFNotSupported   = errors.New("setting do-not-fragment bit is not supported on this platform")
)
