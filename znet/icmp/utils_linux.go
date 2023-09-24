//go:build linux
// +build linux

package icmp

import (
	"golang.org/x/net/bpf"
	"golang.org/x/sys/unix"
	"syscall"
	"unsafe"
)

// makeRecConn4 -- ------------------------------
// --> @Describe make the receive conn v4 client
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeRecFDV4() (fd int, err error) {
	var assembled []bpf.RawInstruction
	if fd,
		err =
		MakeRevFd(
			true,
			p.option.fdReadWriteTimeOutNesc,
			p.option.readBuffer,
			&syscall.SockaddrInet4{
				Port: 0,
				Addr: [4]byte{},
			}); err != nil {
		return
	}

	if !p.option.bpf ||
		len(p.option.filter) < 1 {
		return
	}
	if assembled,
		err =
		bpf.Assemble(p.option.filter); err != nil {
		return
	}
	return fd,
		SetBPF(fd, assembled)
}

// makeRecConn4 -- ------------------------------
// --> @Describe make the receive conn v4 client
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeRecFDV6() (fd int, err error) {
	var assembled []bpf.RawInstruction
	if fd,
		err =
		MakeRevFd(
			false,
			p.option.fdReadWriteTimeOutNesc,
			p.option.readBuffer,
			&syscall.SockaddrInet6{
				Port:   0,
				ZoneId: 0,
				Addr:   [16]byte{},
			}); err != nil {

		return
	}
	if !p.option.bpf ||
		len(p.option.filter) < 1 {
		return
	}
	if assembled,
		err =
		bpf.Assemble(p.option.filter); err != nil {
		return
	}

	return fd,
		SetBPF(fd, assembled)
}

// SetBPF  -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func SetBPF(fd int,
	assembled []bpf.RawInstruction) (
	err error) {
	var errno syscall.Errno
	b := (*[unix.SizeofSockFprog]byte)(unsafe.Pointer(&unix.SockFprog{
		Len:    uint16(len(assembled)),
		Filter: (*unix.SockFilter)(unsafe.Pointer(&assembled[0])),
	}))[:unix.SizeofSockFprog]
	_, _, errno = syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(fd),
		uintptr(unix.SOL_SOCKET),
		uintptr(unix.SO_ATTACH_FILTER),
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(len(b)),
		0)
	if errno != 0 {
		return errnoErr(errno)
	}
	return
}

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case syscall.EAGAIN:
		return syscall.EAGAIN
	case syscall.EINVAL:
		return syscall.EINVAL
	case syscall.ENOENT:
		return syscall.ENOENT
	}
	return e
}
