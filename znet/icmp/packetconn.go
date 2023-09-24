package icmp

import (
	"os"
	"runtime"
	"syscall"
)

// makeSendConn4-- --------------------------
// --> @Describe make the conn4 client
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeSendConn4() (fd int, err error) {
	if fd,
		err =
		syscall.Socket(syscall.AF_INET,
			syscall.SOCK_RAW,
			syscall.IPPROTO_ICMP); nil != err {
		return
	}
	if err =
		syscall.SetNonblock(fd,
			true); nil != err {
		return
	}
	if p.option.writeBuffer > 0 {
		err =
			setWriteBuffer(fd,
				p.option.writeBuffer)
	}

	return

}

// setReadBuffer-- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func setReadBuffer(fd int, bytes int) (err error) {
	if err =
		syscall.SetsockoptInt(fd,
			syscall.SOL_SOCKET,
			syscall.SO_RCVBUF,
			bytes); err != nil {
		return os.NewSyscallError("setReadBuffer:setSockOpt", err)
	}
	runtime.KeepAlive(fd)
	return
}

// setWriteBuffer -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func setWriteBuffer(fd int, bytes int) (err error) {
	if err =
		syscall.SetsockoptInt(fd,
			syscall.SOL_SOCKET,
			syscall.SO_SNDBUF,
			bytes); err != nil {
		return os.NewSyscallError("setWriteBuffer:setSockOpt", err)
	}
	runtime.KeepAlive(fd)
	return
}

// makeConn6 -- --------------------------
// --> @Describe  make the conn6 client
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeSendConn6() (fd int, err error) {
	if fd,
		err =
		syscall.Socket(syscall.AF_INET6,
			syscall.SOCK_RAW,
			syscall.IPPROTO_ICMPV6); nil != err {
		return
	}
	if err =
		syscall.SetNonblock(fd,
			true); nil != err {
		return
	}
	if p.option.writeBuffer > 0 {
		err =
			setWriteBuffer(fd,
				p.option.writeBuffer)
	}

	return
}

// MakeRevFd -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func MakeRevFd(
	v4 bool,
	timeout int64,
	readBuffer int,
	sockAddr syscall.Sockaddr) (
	fd int,
	err error) {
	var family, proto int
	if v4 {
		family,
			proto =
			syscall.AF_INET,
			syscall.IPPROTO_ICMP
	} else {
		family,
			proto =
			syscall.AF_INET6,
			syscall.IPPROTO_ICMPV6
	}

	if fd,
		err =
		syscall.Socket(family,
			syscall.SOCK_RAW,
			proto); nil != err {
		return
	}

	tv := syscall.NsecToTimeval(timeout)
	if err =
		syscall.SetsockoptTimeval(fd,
			syscall.SOL_SOCKET,
			syscall.SO_RCVTIMEO,
			&tv); err != nil {
		return
	}
	if err =
		syscall.SetsockoptTimeval(fd,
			syscall.SOL_SOCKET,
			syscall.SO_SNDTIMEO,
			&tv); err != nil {
		return
	}

	if err = syscall.Bind(fd, sockAddr); nil != err {
		return
	}
	if readBuffer > 0 {
		err = setReadBuffer(fd, readBuffer)
	}
	runtime.KeepAlive(fd)
	return
}
