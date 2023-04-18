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
	Name   string
	Role   string
	Age    int
	Salary int
}

func (u User) Self() {

}

type Default struct {
	Name   string
	Role   string
	Age    int
	Salary int
}

func (d Default) Self() {

}

func main() {

	b := pkg.NewBatchHook[User]()

	task := pkg.InitTaskModel[User]{
		TaskName: "test1",
		Opt: []pkg.Option[User]{
			pkg.WithWaitTime[User](3 * time.Second),
			pkg.WithLoopTime[User](5 * time.Second),
			pkg.WithChanSize[User](1000),
			pkg.WithDoingSize[User](10),
			pkg.WithHandleGoNum[User](3),
			pkg.WithHookFunc[User](func(i []User) bool {
				fmt.Println("=============", len(i))
				fmt.Println("=============", i)

				for _, vl := range i {
					fmt.Println(vl.Age, vl.Role, vl.Salary, vl.Name)
				}
				return true
			}),
			pkg.WithEndHook[User](func(b bool, i ...User) {
				fmt.Println("--------------")
				fmt.Println(b)
				fmt.Println(len(i))
				fmt.Println(i)
			}),
			pkg.WithSavePanic[User](func() {
				if err := recover(); err != nil {
					fmt.Println(err)
					fmt.Println(string(debug.Stack()))
					os.Exit(8)
				}
			}),
		},
	}

	if err := b.InitTask(task); nil != err {
		fmt.Println(err, "============")
		os.Exit(1)
	}

	b.Run()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		i1 := 0
		for {
			if i1 > 20 {
				goto END
			}
			for i := 0; i < 5; i++ {
				i1++
				user := User{Name: "Jinzhu", Age: 0, Role: "Admin", Salary: 200000}

				user.Age = user.Age + i
				user.Name = user.Name + strconv.Itoa(i1)
				user.Role = user.Role + strconv.Itoa(i1)
				b.Submit(&user)
				//item.Item.Params = append(item.Item.Params,user)
				time.Sleep(time.Second)

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

	b.Release()

	b.Wait()
	//_ = b.Wait([]string{"test1", "test2"}...)
	fmt.Println(";;;;;;;;;;;;;;;;;;;;;;;;;")

	fmt.Println("结束")

}
