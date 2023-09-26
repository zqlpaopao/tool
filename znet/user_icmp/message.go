package user_icmp

func checksum(b []byte) uint16 {
	cSumCv := len(b) - 1 // checksum coverage
	s := uint32(0)
	for i := 0; i < cSumCv; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if cSumCv&1 == 0 {
		s += uint32(b[cSumCv])
	}
	s = s>>16 + (s & 0xffff)
	s = s + s>>16
	return ^uint16(s)
}
