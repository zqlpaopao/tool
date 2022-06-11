package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/batch_hook/pkg"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

type User struct {
	Name        string
	Role        string
	Age         int
	Salary int
}


func main() {

	b := pkg.NewBatchHook()

	task := []pkg.InitTaskModel{
		{
			TaskName: "test1",
			Opt: []pkg.Option{
				pkg.WithWaitTime(3 * time.Second),
				pkg.WithLoopTime(5 * time.Second),
				pkg.WithChanSize(1000),
				pkg.WithDoingSize(10),
				pkg.WithHandleGoNum(3),
				pkg.WithHookFunc(func(i []interface{}) bool {
					fmt.Println("=============", len(i))
					fmt.Println("=============", i)

					for _, vl := range i {
						if v1, ok := vl.(User); ok {
							fmt.Println(v1)
						} else {
							os.Exit(1)
						}
					}
					return true
				}),
				pkg.WithEndHook(func(b bool, i ...interface{}) {
					fmt.Println("--------------")
					fmt.Println(b)
					fmt.Println(len(i))
					fmt.Println(i)
				}),
				pkg.WithSavePanic(func(i interface{}) {
					if err := recover(); err != nil {
						fmt.Println(err)
						fmt.Println(string(debug.Stack()))
						os.Exit(8)
					}
				}),
			},
		},
		{
			TaskName: "test2",
			Opt: []pkg.Option{
				pkg.WithWaitTime(3 * time.Second),
				pkg.WithLoopTime(5 * time.Second),
				pkg.WithChanSize(1000),
				pkg.WithDoingSize(10),
				pkg.WithHandleGoNum(3),
				pkg.WithHookFunc(func(i []interface{}) bool {
					fmt.Println("=============", len(i))
					fmt.Println("=============", i)

					for _, vl := range i {
						if v1, ok := vl.(User); ok {
							fmt.Println(v1)
						} else {
							os.Exit(1)
						}
					}
					return true
				}),
				pkg.WithEndHook(func(b bool, i ...interface{}) {
					fmt.Println("--------------")
					fmt.Println(b)
					fmt.Println(len(i))
					fmt.Println(i)
				}),
				pkg.WithSavePanic(func(i interface{}) {
					if err := recover(); err != nil {
						fmt.Println(err)
						fmt.Println(string(debug.Stack()))
						os.Exit(9)

					}
				}),
			},
		},
	}

	if err := b.InitTask(task...); nil != err {
		fmt.Println(err, "============")
		os.Exit(1)
	}

	if err := b.Run([]string{"test1", "test2"}...); nil != err {
		fmt.Println(err)
		os.Exit(4)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		i1 := 0
		for {
			if i1 > 20 {
				goto END
			}
			item := pkg.SubmitModel{
				TaskName: "test1",

			}
			for i := 0; i < 5; i++ {
				i1++
				user := User{Name: "Jinzhu", Age: 0, Role: "Admin", Salary: 200000}

				user.Age = user.Age + i
				user.Name = user.Name + strconv.Itoa(i1)
				user.Role = user.Role + strconv.Itoa(i1)
				item.Data = append(item.Data, &user)
				//item.Item.Params = append(item.Item.Params,user)
				time.Sleep(time.Second)

			}
			if err := b.Submit(item); nil != err {
				fmt.Println(err)
				os.Exit(2)
			}
		}
	END:
		wg.Done()
	}()

	go func() {
		i1 := 0

		for {
			if i1 > 20 {
				goto END
			}
			item := pkg.SubmitModel{
				TaskName: "test2",
			}
			for i := 0; i < 5; i++ {
				i1++
				user := User{Name: "ZHANGSAN", Age: 0, Role: "ZHANGSAN", Salary: 300000}

				user.Age = user.Age + i
				user.Name = user.Name + strconv.Itoa(i1)
				user.Role = user.Role + strconv.Itoa(i1)
				item.Data = append(item.Data, &user)
				//item.Item.Params = append(item.Item.Params,user)
				time.Sleep(time.Second)

			}
			if err := b.Submit(item); nil != err {
				fmt.Println(err)
				os.Exit(2)
			}
		}
	END:
		wg.Done()

	}()

	go func() {
		for {
			fmt.Println("num-go-", runtime.NumGoroutine())
			time.Sleep(2 * time.Second)
		}
	}()

	time.Sleep(20 * time.Second)
	wg.Wait()
	if err := b.Release([]string{"test1", "test2"}...); nil != err {
		fmt.Println(err)
		os.Exit(21)
	}


	b.WaitAll()
	//_ = b.Wait([]string{"test1", "test2"}...)
	fmt.Println(";;;;;;;;;;;;;;;;;;;;;;;;;")



	fmt.Println("结束")


}
