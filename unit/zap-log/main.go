package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/zap-log/src"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

/*
	提供日志分割和日志保存周期控制
*/

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()/4 + 1)
	go func() {
		http.ListenAndServe("127.0.0.1:6001", nil)
	}()
	//go asyncLog()
	for {
		fmt.Println(runtime.NumGoroutine())
		time.Sleep(200 * time.Second)
		os.Exit(3)
	}
}

func init() {
	src.InitLoggerHandler(src.NewLogConfig(
		src.InitInfoPathFileName("./demo_%Y_%m_%d.log"),
		src.InitWarnPathFileName("./demo_%Y_%m_%d.log"),
		src.InitWithMaxAge(0),        //日志最长保存时间，乘以小时 默认禁用
		src.InitWithRotationCount(0), //保存的最大文件数 //默认禁用
		src.InitWithRotationTime(0),  //最大旋转时间 默认值1小时
		src.InitWithIp(1),
		src.InitBufferSize(50),
		src.InitCallFunc(func(level int, s string, back *src.CallBack) {
			fmt.Println(level)
			fmt.Println(s)
			fmt.Println(back)
		})))
	//src.Debug("Debug",s).Msg("Debug")
	//2048 比较合适
	src.NewAsyncLogConfig(src.InitLogAsyncBuffSize(2048), src.InitLogAsyncGoNum(10))
	src.DebugAsync("DebugAsync", "s1").MsgAsync("DebugAsync")
	src.InfoAsync("tag", "s1").MsgAsync("MsgAsync")

}

/*
		20s的时间对比，写入相同内容，异步可以写入
		238M Feb 17 14:23 demo_2022_02_17.log
	 	92M Feb 17 14:26 demo_2022_02_17.log
		差了三倍
*/
func Log() {
	type str struct {
		name string
		age  int
		sex  []int
	}
	s := str{
		name: "name",
		age:  18,
		sex:  []int{1, 2, 3, 4},
	}
	s1 := str{
		name: "name1",
		age:  181,
		sex:  []int{1, 2, 3, 41},
	}

	src.Info("InfoAsync", s).Msg("InfoAsync")

	src.Warn("WarnAsync", s).Msg("WarnAsync")

	src.Debug("DebugAsync", s1).Msg("DebugAsync")

	src.Error("ErrorAsync", s1).Msg("ErrorAsync")

}

func AsyncLog() {
	type str struct {
		name string
		age  int
		sex  []int
	}
	s := str{
		name: "name",
		age:  18,
		sex:  []int{1, 2, 3, 4},
	}
	s1 := str{
		name: "name1",
		age:  181,
		sex:  []int{1, 2, 3, 41},
	}

	//debug info 是一个级别 warn和errorshi 是一个级别，不同级别可分别记录

	src.InfoAsync("InfoAsync", s).MsgAsync("InfoAsync")
	src.WarnAsync("WarnAsync", s).MsgAsync("WarnAsync")
	src.DebugAsync("DebugAsync", s1).MsgAsync("DebugAsync")
	src.ErrorAsync("ErrorAsync", s1).MsgAsync("ErrorAsync")
}
