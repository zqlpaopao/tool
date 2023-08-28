//go:build linux
// +build linux

package ping

import (
	"errors"
	"golang.org/x/net/bpf"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"os"
	"reflect"
	"syscall"

	"golang.org/x/net/icmp"
)

// Returns the length of an ICMP message.
func (p *Ping) getMessageLength() int {
	return p.Size + 8
}

// Attempts to match the ID of an ICMP packet.
func (p *Ping) matchID(ID int) bool {
	// On Linux we can only match ID if we are privileged.
	if p.Protocol == "icmp" {
		return ID == p.Pid
	}
	return true
}

// SetMark sets the SO_MARK socket option on outgoing ICMP packets.
// Setting this option requires CAP_NET_ADMIN.
func (c *icmpConn) SetMark(mark uint) error {
	fd, err := getFD(c.c)
	if err != nil {
		return err
	}
	return os.NewSyscallError(
		"setsockopt",
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, int(mark)),
	)
}

// SetMark sets the SO_MARK socket option on outgoing ICMP packets.
// Setting this option requires CAP_NET_ADMIN.
func (c *icmpV4Conn) SetMark(mark uint) error {
	fd, err := getFD(c.icmpConn.c)
	if err != nil {
		return err
	}
	return os.NewSyscallError(
		"setsockopt",
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, int(mark)),
	)
}

// SetMark sets the SO_MARK socket option on outgoing ICMP packets.
// Setting this option requires CAP_NET_ADMIN.
func (c *icmpV6Conn) SetMark(mark uint) error {
	fd, err := getFD(c.icmpConn.c)
	if err != nil {
		return err
	}
	return os.NewSyscallError(
		"setsockopt",
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_MARK, int(mark)),
	)
}

// SetDoNotFragment sets the do-not-fragment bit in the IP header of outgoing ICMP packets.
func (c *icmpConn) SetDoNotFragment() error {
	fd, err := getFD(c.c)
	if err != nil {
		return err
	}
	return os.NewSyscallError(
		"setsockopt",
		syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_MTU_DISCOVER, syscall.IP_PMTUDISC_DO),
	)
}

// SetDoNotFragment sets the do-not-fragment bit in the IP header of outgoing ICMP packets.
func (c *icmpV4Conn) SetDoNotFragment() error {
	fd, err := getFD(c.icmpConn.c)
	if err != nil {
		return err
	}
	return os.NewSyscallError(
		"setsockopt",
		syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IP, syscall.IP_MTU_DISCOVER, syscall.IP_PMTUDISC_DO),
	)
}

// SetDoNotFragment sets the do-not-fragment bit in the IPv6 header of outgoing ICMPv6 packets.
func (c *icmpV6Conn) SetDoNotFragment() error {
	fd, err := getFD(c.icmpConn.c)
	if err != nil {
		return err
	}
	return os.NewSyscallError(
		"setsockopt",
		syscall.SetsockoptInt(int(fd), syscall.IPPROTO_IPV6, syscall.IPV6_MTU_DISCOVER, syscall.IP_PMTUDISC_DO),
	)
}

// getFD gets the system file descriptor for an icmp.PacketConn
func getFD(c *icmp.PacketConn) (uintptr, error) {
	v := reflect.ValueOf(c).Elem().FieldByName("c").Elem()
	if v.Elem().Kind() != reflect.Struct {
		return 0, errors.New("invalid type")
	}

	fd := v.Elem().FieldByName("conn").FieldByName("fd")
	if fd.Elem().Kind() != reflect.Struct {
		return 0, errors.New("invalid type")
	}

	pfd := fd.Elem().FieldByName("pfd")
	if pfd.Kind() != reflect.Struct {
		return 0, errors.New("invalid type")
	}

	return uintptr(pfd.FieldByName("Sysfd").Int()), nil
}

// makeRecConn4 -- ------------------------------
// --> @Describe make the receive conn v4 client
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeRecConn4() (conn net.PacketConn, err error) {
	var (
		pConn     *ipv4.RawConn
		assembled []bpf.RawInstruction
	)
	if conn, err = net.ListenPacket(ipv4Proto[p.option.protocol], p.option.source); err != nil {
		panic(err)
	}

	cc, ok := conn.(*net.IPConn)
	if !ok {
		return nil, errors.New("makeRecConn4 change net.IPConn is err")
	}

	if err = cc.SetReadBuffer(p.option.readBuffer); nil != err {
		return
	}
	if err = cc.SetWriteBuffer(p.option.writeBuffer); nil != err {
		return
	}

	if pConn, err = ipv4.NewRawConn(conn); nil != err {
		panic(err)
	}

	if !p.option.bpf || len(p.option.filter) < 1 {
		return
	}

	if assembled, err = bpf.Assemble(p.option.filter); err != nil {
		return
	}
	return conn, pConn.SetBPF(assembled)
}

// makeRecConn6 -- ------------------------------
// --> @Describe make the receive conn v4 client
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeRecConn6() (conn net.PacketConn, err error) {
	var (
		pConn     *ipv6.PacketConn
		assembled []bpf.RawInstruction
	)
	if conn, err = net.ListenPacket(ipv6Proto[p.option.protocol], p.option.source); err != nil {
		panic(err)
	}

	cc, ok := conn.(*net.IPConn)
	if !ok {
		return nil, errors.New("makeRecConn4 change net.IPConn is err")
	}

	if err = cc.SetReadBuffer(p.option.readBuffer); nil != err {
		return
	}
	if err = cc.SetWriteBuffer(p.option.writeBuffer); nil != err {
		return
	}

	pConn = ipv6.NewPacketConn(conn)

	if !p.option.bpf {
		return
	}

	if assembled, err = bpf.Assemble(p.option.filter); err != nil {
		return
	}
	if err = pConn.SetBPF(assembled); nil != err {
		return
	}
	return
}
