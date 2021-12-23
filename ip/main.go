package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/ip/src"
)

func main(){
	fmt.Println(src.GetEth0())
	fmt.Println(src.RefreshAndGetAllIp(2))
	fmt.Println(src.RefreshAndGetIp(2))
}