package common

import "time"

//Config web 相关
type Config struct {
	Env EnvConfig
}

type WebConfig struct {
	Host         string
	Port         string
	NoAuthUrl    []string
	NoAuthUrlMap map[string]struct{}
}

type Pprof struct {
	OpenTag bool
	Host    string
	Port    string
}

type EnvConfig struct {
	VersionInfo bool
	PProf       Pprof
	Mode        string
	Web         WebConfig
}

//LogConfig 相关
type LogConfig struct {
	Log struct{
		InfoPath      string
		WarnPath      string
		MaxAge        int
		RotationCount uint
		RotationTime  int
		WithIp        int
		BufferSize    int
		AsyncBuffSize int
		AsyncGoNum    int
	}
}

//ReqLogInfo 请求日志记录
type ReqLogInfo struct {
	StartTime string
	EndTime   string
	RunTime   time.Duration
	ReqMethod string
	ReqUrl    string
	ClientIP  string
	ReqArgs   string
	RespCode  int
	RespInfo  string
}
