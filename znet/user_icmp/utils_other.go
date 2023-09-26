//go:build !linux && !windows
// +build !linux,!windows

package user_icmp

import (
	"runtime"
	"syscall"
)

// makeRecConn4 -- ------------------------------
// --> @Describe make the receive conn v4 client
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeRecFDV4() (fd int, err error) {
	if fd,
		err =
		syscall.Socket(
			syscall.AF_INET,
			syscall.SOCK_DGRAM,
			syscall.IPPROTO_ICMP); err != nil {
		panic(err)
	}
	tv := syscall.NsecToTimeval(2000000000)
	if err =
		syscall.SetsockoptTimeval(fd,
			syscall.SOL_SOCKET,
			syscall.SO_RCVTIMEO,
			&tv); err != nil {
		panic(err)
	}
	//write
	if err =
		syscall.SetsockoptTimeval(fd,
			syscall.SOL_SOCKET,
			syscall.SO_SNDTIMEO,
			&tv); err != nil {
		panic(err)
	}
	err =
		syscall.Bind(fd, &syscall.SockaddrInet4{
			Port: 0,
			Addr: [4]byte{},
		})
	if p.option.readBuffer > 0 {
		err = setReadBuffer(fd, p.option.readBuffer)
	}
	return

}

// makeRecConn4 -- ------------------------------
// --> @Describe make the receive conn v4 client
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeRecFDV6() (fd int, err error) {
	if fd,
		err =
		syscall.Socket(syscall.AF_INET6,
			syscall.SOCK_RAW,
			syscall.IPPROTO_ICMPV6); err != nil {
		panic(err)
	}

	tv := syscall.NsecToTimeval(2000000000)
	if err =
		syscall.SetsockoptTimeval(fd,
			syscall.SOL_SOCKET,
			syscall.SO_RCVTIMEO,
			&tv); err != nil {
		panic(err)
	}

	if err =
		syscall.SetsockoptTimeval(fd,
			syscall.SOL_SOCKET,
			syscall.SO_SNDTIMEO,
			&tv); err != nil {
		panic(err)
	}

	if err =
		syscall.Bind(fd, &syscall.SockaddrInet6{
			Port:   0,
			ZoneId: 0,
			Addr:   [16]byte{},
		}); nil != err {
		return
	}
	if p.option.readBuffer > 0 {
		err = setReadBuffer(fd, p.option.readBuffer)
	}
	runtime.KeepAlive(fd)
	return
}
