package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/string-byte/src"
)

func main(){
	b := src.String2Bytes("")
	fmt.Println(b)
	t := src.Bytes2String(nil)
	fmt.Printf("%#v",&t)



}
