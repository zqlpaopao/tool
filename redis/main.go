package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func main() {
	//lua := pkg.MakeNewScript(pkg.GetStrSetStr)

	redis := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	//str := redis.ScriptLoad(context.TODO(),pkg.IncrNum)
	str := "53b0036f258e3fe5393a8aed723a828b10483685"
	fmt.Println(str)

	res := redis.EvalSha(context.TODO(), str, []string{"testa"}, []interface{}{18,4,400})


	//res := redis.Eval(context.TODO(), pkg.IncrNum, []string{"testa"}, []interface{}{13,4,400})
	fmt.Println(res.StringSlice())

}
