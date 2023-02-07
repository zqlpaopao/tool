package src

import (
	"net"
	"sync"
)

const (
	refreshTag = 1 << iota
)
//This is a global variable. The write operation is much smaller than the read operation.
//The read is concurrent
type ipInfo struct {
	ipInfos *sync.Map
	sync.Once
}

var ipInfos = new(ipInfo)

func init(){
	ipInfos.Once.Do(initIpInfo)
}

//init
func initIpInfo(){
	ipInfos = &ipInfo{
		ipInfos: &sync.Map{},
	}
	ipInfos.getLocalIp()
}

//GetLocalIp get local ip
func (i *ipInfo)getLocalIp(){
	var (
		addr      []net.Interface
		byName    *net.Interface
		addresses []net.Addr
		err       error
	)
	if addr, err = net.Interfaces(); nil != err {
		return
	}
	for _, v := range addr {
		if byName, err = net.InterfaceByName(v.Name); err != nil {
			continue
		}
		if addresses, err = byName.Addrs(); nil != err {
			continue
		}
		for _, v1 := range addresses {
			var ip net.IP
			switch v := v1.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			if byName.Name != "en0" && byName.Name != "eth0" {
				continue
			}
			i.ipInfos.Store("eth0",ip.String())
		}
	}
}

//GetEth0 direct acquisition
func GetEth0()string{
	if v,ok := ipInfos.ipInfos.Load("eth0");ok{
		if o,ok1 := v.(string);ok1{
			return o
		}
		return ""
	}
	if v,ok := ipInfos.ipInfos.Load("en0");ok{
		if o,ok1 := v.(string);ok1{
			return o
		}
		return ""
	}
	return ""
}

//RefreshAndGetIp Get first, not refreshing
//Or force refresh acquisition
func RefreshAndGetIp(tag int)string{
	if v := GetEth0();v != "" && tag != refreshTag{
		return v
	}
	ipInfos.ipInfos = &sync.Map{}
	ipInfos.getLocalIp()
	return GetEth0()
}

//RefreshAndGetAllIp get all
func RefreshAndGetAllIp(tag int)(ipInfo map[string]string){
	ipInfo = make(map[string]string)
	if tag != refreshTag{
		goto END
	}
	ipInfos.ipInfos = &sync.Map{}
	ipInfos.getLocalIp()

END:
	ipInfos.ipInfos.Range(func(key, value interface{}) bool {
		if _ ,ok := key.(string);!ok{
			return true
		}
		if _ ,ok := value.(string);!ok{
			return true
		}
		ipInfo[key.(string)] = value.(string)
		return true
	})
	return
}
