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