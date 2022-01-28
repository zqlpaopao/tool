package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/version-num-manager/src"
)

//支持分支_提交者_提交时间_提交内容_hashTag 的信息组合
//2022-01-11 10:53:29 Version: main_zql_2022-01-04 20:34:45_testMsg_85e24f6
func main(){
	getVersion()
}

func getVersion(){
	src.OpenTag = true
	err := src.NewVersionNumManager(
		src.WithNotAuth(false),
		src.WithBranch(true),
		src.WithPrint(true),
		src.WithTag("Version: "),

	).Do().Error()
	if err != nil{
		return
	}

	fmt.Println(src.GetVersion())
}