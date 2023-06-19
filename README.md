# ssh-tool
通过代理连接线上mysql、redis等
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
go get 或者引入包

src.InitLoggerHandler(src.NewLogConfig(
		src.InitInfoPathFileName("./demo_%Y_%m_%d.log"),
		src.InitWarnPathFileName("./demo_%Y_%m_%d.log"),
		src.InitWithMaxAge(0),        //日志最长保存时间，乘以小时 默认禁用
		src.InitWithRotationCount(0), //保存的最大文件数 //默认禁用
		src.InitWithRotationTime(0),  //最大旋转时间 默认值1小时
		src.InitWithIp(1),
		src.InitBufferSize(50)))
	//2048 比较合适
	src.NewAsyncLogConfig(src.InitLogAsyncBuffSize(2048), src.InitLogAsyncGoNum(10))
```
基准测试
```
write setup code here...
goos: darwin
goarch: arm64
pkg: github.com/zqlpaopao/tool/zap-log
BenchmarkLog-10                 1000000000               0.3152 ns/op          0 B/op          0 allocs/op
BenchmarkAsyncLog-10            1000000000               0.2413 ns/op          0 B/op          0 allocs/op
PASS
write teardown code here...
ok      github.com/zqlpaopao/tool/zap-log       8.989s
```
1亿次写入，每次快近0.1ns，速度和开启的缓冲池和协程的大小有关

20s的时间写入速度
```
20s的时间对比，写入相同内容，异步可以写入
	238M Feb 17 14:23 demo_2022_02_17.log
 	92M Feb 17 14:26 demo_2022_02_17.log
	差了三倍
```

## string-byte
string-byte 包支持string到[]byte的转换，不会有err信息，nil的[]byte转换为""字符串
会比string([]byte())和 []byte(string) 快

## ip
获取机的ip信息，支持ipv6


## config
获取配置文件内容 viper封装

## json
针对不确定的类型转为interface的时候int 转为float导致精度丢失
```
var request = `{"id":7044144249855934983}`
r ,e := src.MarshalWithInterface(request)
fmt.Println(r ,e)//map[id:7044144249855934983] <nil>


正常的
原始的数据{"id":7044144249855934983,"name":"demo"}， 转化后的数据{"id":7044144249855935000,"name":"demo"}
```

## 并发下载大文件 downloader

## 版本号生成工具
version-num-manager
```
_"github.com/zqlpaopao/tool/version-num-manager/src"
```
样式
```
支持分支_提交者_提交时间_提交内容_hashTag 的信息组合
2022-01-11 10:53:29 Version: main_zql_2022-01-04 20:34:45_testMsg_85e24f6
```
也可以自定义

## mysql支持的工具
1. 支持ssh获取原生mysql客户端
2. 支持ssh获取gorm的客户端
3. 支持ON DUPLICATE KEY UPDATE
4. 支持insert sql
5. 支持update sql
6. 支持struct转where sql
7. 支持map 转where sql
8. 支持gorm的created&query&delete&update&raw操作的前后触发插件

## 函数重试 tetry
```
package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/retry/pkg"
)

func main(){
	pkg.NewRetryManager(
		//pkg.WithRetryCount(5),
		//pkg.WithRetryInterval(4*time.Second),
		pkg.WithDelayType(pkg.WithDefaultDelayType),//指数级别重试
		).RegisterRetryCallback(func(u uint) {
		fmt.Println("这是重试的次数",u)
	}).RegisterCompleteCallback(func(u uint, b bool, i ...interface{}) {
		fmt.Println("这是重试的次数，结果，和传递参数",u,b,i)
	}).DoSync(func() bool {
		fmt.Println("这是重试方法")
		return false
	},"111",[]string{"a"})
}

//默认重试间隔3s ，可以根据自己的需要调整 秒，纳秒，毫秒 pkg.WithRetryInterval(4*time.Second)
//默认重试次数3次 ，可以根据需要调整
//支持 每次重试回调方法RegisterRetryCallback
//支持 重试完成回调方法，结果看b bool RegisterCompleteCallback
```
## redis lua 脚本

原子递增
```
//IncrNum num 增加操作
//KEYS[1] 是key
//ARGV[1] 是总数
//ARGV[2] 要增加的数值
//ARGV[3] 过期时间，默认是-1
//1 是成功 返回增加后的值
//2 是失败 返回现在的值或者错误
```
原子操作
```
//GetStrSetStr string 操作，获取key是否存在，不存在设置key及过期时间
//1代表成功，2代表失败
//没有对设置key的过程结果做处理
```
```
"github.com/zqlpaopao/tool/redis/pkg"
```

# 多节点、多集群kafka写入思路
[github]（https://github.com/smallnest/gofer/tree/master）
