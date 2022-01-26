package common

import "time"

//Config web 相关
type Config struct {
	Env EnvConfig
}

type WebConfig struct {
	Host string
	Port string
	NoAuthUrl []string
	NoAuthUrlMap map[string]struct{}
}

type Pprof struct {
	OpenTag bool
	Host string
	Port string
}

type EnvConfig struct {
	Web WebConfig
	PProf Pprof
	Mode string
}



//请求日志记录
// 结束时间
//endTime := time.Now()
//// 执行时间
//latencyTime := endTime.Sub(startTime)
//// 请求方式
//reqMethod := g.Request.Method
//// 请求路由
//reqUri := g.Request.RequestURI
//// 状态码
//statusCode := g.Writer.Status()
//// 请求IP
//clientIP := g.ClientIP()
////请求参数
//req ,err := requestParams(g)

type ReqLogInfo struct {
	StartTime string
	EndTime string
	RunTime time.Duration
	ReqMethod string
	ReqUrl string
	ClientIP string
	ReqArgs string
	RespCode int
	RespInfo string
}