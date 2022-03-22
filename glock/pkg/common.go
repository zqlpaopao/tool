package pkg

import (
	"io/ioutil"
	"strings"
	"time"
)

const (
	DefaultExpireTIme   = 3               //s
	DefaultSeizeTIme    = 1 * time.Second //s
	DefaultRedisTimeOut = 3 * time.Second //s
	memberGroup         = "member:group:"
	Lock                = "member:master"
	Master              = "master"
	Slave               = "slave"
)

//GetHostName get host name
func GetHostName() (serverHostName string, err error) {
	contents, err := ioutil.ReadFile("/etc/hostname")
	if err == nil {
		//因为contents是[]byte类型，直接转换成string类型后会多一行空格,需要使用strings.Replace替换换行符
		serverHostName = strings.Replace(string(contents), "\n", "", 1)
	}
	return serverHostName, err
}
