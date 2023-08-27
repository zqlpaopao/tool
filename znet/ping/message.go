package ping

const (
	ProtocolIPv6ICMP = 58 // ICMP for IPv6
)

func checksum(b []byte) uint16 {
	csumcv := len(b) - 1 // checksum coverage
	s := uint32(0)
	for i := 0; i < csumcv; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if csumcv&1 == 0 {
		s += uint32(b[csumcv])
	}
	s = s>>16 + (s & 0xffff)
	s = s + s>>16
	return ^uint16(s)
}
