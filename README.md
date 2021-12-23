# ssh-tool
```
if db, err = src2.NewSSHMysqlClient(&src2.Config{
		Addr:   "xx.xx.xx.xx:22",
		User:   "xxx",
		Passwd: "xx.xx.xx",
	}, &src2.MysqlConfig{
		UserName: "xxx",
		PassWd:   "^xxxxx",
		IpPort:   "xx.xx.xx.xx:3306",
		Dbname:   "xxxxxx",
	}); nil != err {
		panic(err)
	}
```


# zap-log
1. 支持日志分割 可具体到天
2. 支持日志的最大数量，默认禁用
3. 支持日志的最大存活周期，默认禁用
4. 支持文件的旋转周期，默认1小时
5. 支持记录ip信息，多节点部署更容易定位问题

```
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
	//debug info 是一个级别 warn和errorshi 是一个级别，不同级别可分别记录
	src.InitLoggerHandler(&src.LogConfig{
		InfoPathFileName: "./demo.log",
		WarnPathFileName: "./demo.log",
		//WithRotationTime: //最大旋转时间 默认值1小时
		//WithMaxAge: //日志最长保存时间，乘以小时 默认禁用
		//WithRotationCount: //保存的最大文件数 //默认禁用
	})
	src.Info("Info",s).Msg("Info")
	src.Warn("Warn",s).Msg("Warn")
	src.Error("Error",s).Msg("Error")
	src.Debug("Debug",s).Msg("Debug")

	src.Warn("Warn",s).Msg("")

}
```
```
{"level":"INFO","time":"2021-12-23 09:59:11","file":"zap-log/main.go:33","msg":"Info","localIp":"10.254.45.120/22","params":["{name:name age:18 sex:[1 2 3 4]}"],"errMsg":"Info"}
{"level":"WARN","time":"2021-12-23 09:59:11","file":"zap-log/main.go:34","msg":"Warn","localIp":"10.254.45.120/22","params":["{name:name age:18 sex:[1 2 3 4]}"],"errMsg":"Warn"}
{"level":"ERROR","time":"2021-12-23 09:59:11","file":"zap-log/main.go:35","msg":"Error","localIp":"10.254.45.120/22","params":["{name:name age:18 sex:[1 2 3 4]}"],"errMsg":"Error"}
{"level":"DEBUG","time":"2021-12-23 09:59:11","file":"zap-log/main.go:36","msg":"Debug","localIp":"10.254.45.120/22","params":["{name:name age:18 sex:[1 2 3 4]}"],"errMsg":"Debug"}
{"level":"WARN","time":"2021-12-23 09:59:11","file":"zap-log/main.go:38","msg":"Warn","localIp":"10.254.45.120/22","params":["{name:name age:18 sex:[1 2 3 4]}"],"errMsg":""}

```

## string-byte
string-byte 包支持string到[]byte的转换，不会有err信息，nil的[]byte转换为""字符串
会比string([]byte())和 []byte(string) 快
