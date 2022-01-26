package start

import (
	config "github.com/zqlpaopao/tool/config/src"
	log "github.com/zqlpaopao/tool/zap-log/src"
)

func init(){
	if err := config.Ctx.Init("CONF_PATH");nil != err{
		panic(err)
	}
	InitLog()
}


//InitLog 初始化日志模块
func InitLog(){
	var (
		logPath string
		err error
	)
	if logPath,err = config.Ctx.GetLogLayout("log","path");nil != err{
		panic(err)
	}
	//初始化日志
	log.InitLoggerHandler(log.NewLogConfig(log.InitInfoPathFileName(logPath), log.InitWarnPathFileName(logPath), log.InitWithMaxAge(0),log.InitWithRotationCount(0),log.InitWithRotationTime(0),log.InitWithIp(1)))
}