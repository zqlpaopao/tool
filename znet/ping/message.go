package ping

import (
	"encoding/binary"
)

const (
	ProtocolIPv6ICMP = 58 // ICMP for IPv6
)

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

// ParseMessage parses b as an ICMP message.
// The provided proto must be either the ICMPv4 or ICMPv6 protocol
// number.
func ParseMessage(proto int, b []byte) (bool, int, int, []byte, error) {
	if len(b) < 4 {
		return false, 0, 0, nil, errMessageTooShort
	}
	var Type int
	//m := &Message{Code: int(b[1]), Checksum: int(binary.BigEndian.Uint16(b[2:4]))}
	//switch proto {
	//case iana.ProtocolICMP:
	//	m.Type = ipv4.ICMPType(b[0])
	//case iana.ProtocolIPv6ICMP:
	//	m.Type = ipv6.ICMPType(b[0])
	//default:
	//	return nil, errInvalidProtocol
	//}

	switch proto {
	case ProtocolICMP:
		Type = int(b[0])
	case ProtocolIPv6ICMP:
		Type = int(b[0])
	default:
		return false, 0, 0, nil, errInvalidProtocol
	}
	if Type != ICMPTypeEchoV4Reply && Type != ICMPTypeEchoV6Reply {
		return parseRawBody(b[4:])
	} else {
		return parseEcho(b[4:])
	}

	//if fn, ok := parseFns[m.Type]; !ok {
	//	m.Body, err = parseRawBody(proto, b[4:])
	//} else {
	//	m.Body, err = fn(proto, b[4:])
	//}
	//if err != nil {
	//	return nil, err
	//}
	//return m, nil
}

// parseEcho parses b as an ICMP echo request or reply message body.
func parseEcho(b []byte) (bool, int, int, []byte, error) {
	bodyLen := len(b)
	if bodyLen < 4 {
		return false, 0, 0, nil, errBodyTooShort
	}
	//p := &Echo{ID: int(binary.BigEndian.Uint16(b[:2])), Seq: int(binary.BigEndian.Uint16(b[2:4]))}
	if bodyLen > 4 {
		return true, int(binary.BigEndian.Uint16(b[:2])), int(binary.BigEndian.Uint16(b[2:4])), b[4:], nil
		//p.Data = make([]byte, bodyLen-4)
		//copy(p.Data, b[4:])
	}

	return false, 0, 0, nil, errBodyTooShort
}

// parseRawBody parses b as an ICMP message body.
func parseRawBody(b []byte) (bool, int, int, []byte, error) {
	//p := &RawBody{Data: make([]byte, len(b))}
	//copy(p.Data, b)
	//return p, nil
	return false, 0, 0, nil, nil
}
