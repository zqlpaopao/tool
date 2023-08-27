package ping

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/net/bpf"
	"time"
)

type Filter []bpf.Instruction

const (
	Udp = "udp"
	IP  = "ip"

	//PrepareChLen is the object *Ping of ready to send
	PrepareChLen = 10

	//ErrorInfoSize is the pool error chan
	ErrorInfoSize = 3

	NoIcmpAndDarwin = 20

	timeSliceLength  = 8
	trackerLength    = len(uuid.UUID{})
	protocolICMP     = 1
	protocolIPv6ICMP = 58
)

var (
	seed = time.Now().UnixNano()

	ipv4Proto = map[string]string{"icmp": "ip4:icmp", "udp": "udp4"}
	ipv6Proto = map[string]string{"icmp": "ip6:ipv6-icmp", "udp": "udp6"}

	ErrMarkNotSupported = errors.New("setting SO_MARK socket option is not supported on this platform")
	ErrDFNotSupported   = errors.New("setting do-not-fragment bit is not supported on this platform")
)
