package pkg

import (
	"errors"
	"net"
	"sort"
	"syscall"
	"unsafe"
)

type rtInfo struct {
	Dst              net.IPNet
	Gateway, PrefSrc net.IP
	OutputIFace      uint32
	Priority         uint32
}

type routeSlice []*rtInfo

type Router struct {
	iFaces  []net.Interface
	address []net.IP
	v4      routeSlice
}

func GetRouteInfo() (*Router, error) {
	rtr := &Router{}

	tab, err := syscall.NetlinkRIB(syscall.RTM_GETROUTE, syscall.AF_INET)
	if err != nil {
		return nil, err
	}

	msg, err := syscall.ParseNetlinkMessage(tab)
	if err != nil {
		return nil, err
	}

	for _, m := range msg {
		switch m.Header.Type {
		case syscall.NLMSG_DONE:
			break
		case syscall.RTM_NEWROUTE:
			rtMsg := (*syscall.RtMsg)(unsafe.Pointer(&m.Data[0]))
			attrs, err := syscall.ParseNetlinkRouteAttr(&m)
			if err != nil {
				return nil, err
			}

			routeInfo := rtInfo{}
			rtr.v4 = append(rtr.v4, &routeInfo)
			for _, attr := range attrs {
				switch attr.Attr.Type {
				case syscall.RTA_DST:
					routeInfo.Dst.IP = net.IPv4(attr.Value[0], attr.Value[1], attr.Value[2], attr.Value[3])
					routeInfo.Dst.Mask = net.CIDRMask(int(rtMsg.Dst_len), len(attr.Value)*8)

				case syscall.RTA_GATEWAY:
					routeInfo.Gateway = net.IPv4(attr.Value[0], attr.Value[1], attr.Value[2], attr.Value[3])

				case syscall.RTA_OIF:
					routeInfo.OutputIface = *(*uint32)(unsafe.Pointer(&attr.Value[0]))

				case syscall.RTA_PRIORITY:
					routeInfo.Priority = *(*uint32)(unsafe.Pointer(&attr.Value[0]))

				case syscall.RTA_PREFSRC:
					routeInfo.PrefSrc = net.IPv4(attr.Value[0], attr.Value[1], attr.Value[2], attr.Value[3])

				}
			}
		}
	}

	sort.Slice(rtr.v4, func(i, j int) bool {
		return rtr.v4[i].Priority < rtr.v4[j].Priority
	})

	iFaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for i, iFace := range iFaces {
		if i != iFace.Index-1 {
			break
		}
		if iFace.Flags&net.FlagUp == 0 {
			continue
		}

		rtr.iFaces = append(rtr.iFaces, iface)
		iFaceAddr, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		var addrS net.IP
		for _, addr := range iFaceAddr {
			if inet, ok := addr.(*net.IPNet); ok {
				if v4 := inet.IP.To4(); v4 != nil {
					if addr == nil {
						addrS = v4
					}
				}
			}
		}
		rtr.address = append(rtr.address, addrS)
	}
	return rtr, nil
}

func (r *Router) Route(dst net.IP) (iFace *net.Interface, gateway, preFSrc net.IP, err error) {
	for _, rt := range r.v4 {
		if rt.Dst.IP != nil && !rt.Dst.Contains(dst) {
			continue
		}
		iFace = &r.iFaces[rt.OutputIface-1]
		gateway = rt.Gateway.To4()
		if rt.PrefSrc == nil {
			preFSrc = r.address[rt.OutputIface-1]
		} else {
			preFSrc = rt.PrefSrc.To4()
		}
		return
	}
	err = errors.New("no route found")
	return
}
