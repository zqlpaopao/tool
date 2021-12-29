package main

import (
	"github.com/zqlpaopao/tool/zap-log/src"
)

/*
	提供日志分割和日志保存周期控制
 */

func main(){
	type str struct{
		name string
		age int
		sex []int
	}
	s := str{
		name: "name",
		age:  18,
		sex:  []int{1,2,3,4},
	}
	//debug info 是一个级别 warn和errorshi 是一个级别，不同级别可分别记录
	src.InitLoggerHandler(src.NewLogConfig(
		src.InitInfoPathFileName("./demo.log"),
		src.InitWarnPathFileName("./demo.log"),
		src.InitWithMaxAge(0),//日志最长保存时间，乘以小时 默认禁用
		src.InitWithRotationCount(0),//保存的最大文件数 //默认禁用
		src.InitWithRotationTime(0),//最大旋转时间 默认值1小时
		src.InitWithIp(1)))



	src.Info("Info",s).Msg("Info")
	src.Warn("Warn",s).Msg("Warn")
	src.Error("Error",s).Msg("Error")
	src.Debug("Debug",s).Msg("Debug")

	src.Warn("Warn",s).Msg("")
	//{"level":"INFO","time":"2021-12-17 18:08:19","file":"zap-log/main_test.go:27","msg":"Info","params":["{name:name age:18 sex:[1 2 3 4]}"],"errMsg":"Info"}
	//{"level":"WARN","time":"2021-12-17 18:08:19","file":"zap-log/main_test.go:28","msg":"Warn","params":["{name:name age:18 sex:[1 2 3 4]}"],"errMsg":"Warn"}
	//{"level":"ERROR","time":"2021-12-17 18:08:19","file":"zap-log/main_test.go:29","msg":"Error","params":["{name:name age:18 sex:[1 2 3 4]}"],"errMsg":"Error"}
	//{"level":"DEBUG","time":"2021-12-17 18:08:19","file":"zap-log/main_test.go:30","msg":"Debug","params":["{name:name age:18 sex:[1 2 3 4]}"],"errMsg":"Debug"}
	//{"level":"WARN","time":"2021-12-17 18:08:19","file":"zap-log/main_test.go:32","msg":"Warn","params":["{name:name age:18 sex:[1 2 3 4]}"],"errMsg":""}

}