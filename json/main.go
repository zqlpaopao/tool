package main

import (
	"encoding/json"
	"fmt"
	"github.com/zqlpaopao/tool/json/src"
	"reflect"
)
//https://mp.weixin.qq.com/s?__biz=MzI1MzYzMTI2Ng==&mid=2247486755&idx=1&sn=89d655556de92d51ef55ae400b854a14&scene=21#wechat_redirect
func main() {
	number()
	var request = `{"id":7044144249855934983}`
	r ,e := src.MarshalWithInterface(request)
	fmt.Println(r ,e)

}

//当不确定的数据类型转换为interface的时候，int会转化为float，导致int数据丢失精度
func number(){
	var request = `{"id":7044144249855934983,"name":"demo"}`

	var test interface{}
	err := json.Unmarshal([]byte(request), &test)
	if err != nil {
		fmt.Println("error:", err)
	}

	obj := test.(map[string]interface{})

	dealStr, err := json.Marshal(test)
	if err != nil {
		fmt.Println("error:", err)
	}

	id := obj["id"]

	// 反序列化之后重新序列化打印
	fmt.Println(string(dealStr))
	fmt.Printf("%+v\n", reflect.TypeOf(id).Name())
	fmt.Printf("%+v\n", id.(float64))
	fmt.Printf("原始的数据%s， 转化后的数据%s",request,dealStr)
}