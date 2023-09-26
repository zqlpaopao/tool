//go:build windows
// +build windows

package user_icmp

// makeRecConn4 -- ------------------------------
// --> @Describe make the receive conn v4 client
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeRecFDV4() (fd int, err error) {
	return MakeRevFd(&syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{},
	},
		true)
}

// makeRecConn4 -- ------------------------------
// --> @Describe make the receive conn v4 client
// --> @params
// --> @return
// -- ------------------------------------
func (p *Pool) makeRecFDV6() (fd int, err error) {
	return MakeRevFd(&syscall.SockaddrInet6{
		Port:   0,
		ZoneId: 0,
		Addr:   [16]byte{},
	},
		false)
}
