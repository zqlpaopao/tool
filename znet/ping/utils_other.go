//go:build !linux && !windows
// +build !linux,!windows

package ping

import (
	"net"
)

// Returns the length of an ICMP message.
func (p *Ping) getMessageLength() int {
	return p.Size + 8
}

// Attempts to match the ID of an ICMP packet.
func (p *Ping) matchID(ID int) bool {
	return ID == p.Pid
}

// SetMark sets the SO_MARK socket option on outgoing ICMP packets.
// Setting this option requires CAP_NET_ADMIN.
func (c *icmpConn) SetMark(mark uint) error {
	return ErrMarkNotSupported
}

// SetMark sets the SO_MARK socket option on outgoing ICMP packets.
// Setting this option requires CAP_NET_ADMIN.
func (c *icmpV4Conn) SetMark(mark uint) error {
	return ErrMarkNotSupported
}

// SetMark sets the SO_MARK socket option on outgoing ICMP packets.
// Setting this option requires CAP_NET_ADMIN.
func (c *icmpV6Conn) SetMark(mark uint) error {
	return ErrMarkNotSupported
}

// SetDoNotFragment sets the do-not-fragment bit in the IP header of outgoing ICMP packets.
func (c *icmpConn) SetDoNotFragment() error {
	return ErrDFNotSupported
}

// SetDoNotFragment sets the do-not-fragment bit in the IP header of outgoing ICMP packets.
func (c *icmpV4Conn) SetDoNotFragment() error {
	return ErrDFNotSupported
}

// SetDoNotFragment sets the do-not-fragment bit in the IPv6 header of outgoing ICMPv6 packets.
func (c *icmpV6Conn) SetDoNotFragment() error {
	return ErrDFNotSupported
}

// makeRecConn4 -- ------------------------------
// --> @Describe make the receive conn v4 client
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeRecConn4() (conn net.PacketConn, err error) {
	return net.ListenPacket(ipv4Proto[p.option.protocol], p.option.source)

}

// makeRecConn4 -- ------------------------------
// --> @Describe make the receive conn v4 client
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeRecConn6() (conn net.PacketConn, err error) {
	return net.ListenPacket(ipv6Proto[p.option.protocol], p.option.source)
}
