package main

import (
	"fmt"
	config "github.com/zqlpaopao/tool/config/src"
)

/*
CONF_DIR=/Users/zql/Desktop/test/tool/config/testconf go run main.go

 */

func main(){

	var (
		logPth string
		err error
	)
	//初始化环境变量
	if err := config.Ctx.Init("CONF_DIR");nil != err{
		panic(err)
	}

	if logPth,err = config.Ctx.GetLogLayout("log-web","path");nil != err{
		panic(err)
	}
	fmt.Println(logPth)
}
