package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zqlpaopao/tool/redis/pkg"
)

func main() {
	//lua := pkg.MakeNewScript(pkg.GetStrSetStr)

	redis := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	res := redis.Eval(context.TODO(), pkg.GetStrSetStr, []string{"test1"}, []interface{}{"test1",600})
	fmt.Println(res.StringSlice())
}
