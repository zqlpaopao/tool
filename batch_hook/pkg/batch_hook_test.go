package pkg

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"runtime/debug"
	"strconv"
	"testing"
	"time"
)

var DB *gorm.DB

func init() {
	var err error
	dsn := "xxxx:xxxx@tcp(127.0.0.1:xxxx)/test?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Test1 struct {
	Id          uint64    `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	AddrTag     string    `json:"addr_tag"`
	Role        string    `json:"role"`
	Test1       string    `json:"test1"`
	Test2       string    `json:"test2"`
	Test3       string    `json:"test3"`
	Test4       string    `json:"test4"`
	Test5       string    `json:"test5"`
	Test6       string    `json:"test6"`
	Test7       string    `json:"test7"`
	Test8       string    `json:"test8"`
	Test9       int32     `json:"test9"`
	CreatedTime time.Time `json:"created_time"`
}

func (_ Test1) TableName() string {
	return "test1"
}
func BenchmarkInsert(b *testing.B) {
	t := time.Now()
	for i := 0; i < 10000; i++ {
		item := Test1{
			Name:        "zhangSan-" + strconv.Itoa(i),
			Address:     "China",
			AddrTag:     "china-bj",
			Role:        "admin",
			Test1:       "test1" + strconv.Itoa(i),
			Test2:       "test2" + strconv.Itoa(i),
			Test3:       "test3" + strconv.Itoa(i),
			Test4:       "test4" + strconv.Itoa(i),
			Test5:       "test5" + strconv.Itoa(i),
			Test6:       "test6" + strconv.Itoa(i),
			Test7:       "test7" + strconv.Itoa(i),
			Test8:       "test8" + strconv.Itoa(i),
			Test9:       0,
			CreatedTime: time.Now(),
		}

		if err := DB.
			Model(&Test1{}).
			//Debug().
			Create(&item).
			Error; nil != err {
			fmt.Println("=========", err)
		}
	}
	to := time.Now().Sub(t)
	fmt.Println("BenchmarkInsert", to)

}

func BenchmarkNewBatchHook(b *testing.B) {
	t := time.Now()
	var err error
	df := NewBatchHook[Test1]()
	task := InitTaskModel[Test1]{
		TaskName: "test1",
		Opt: []Option[Test1]{
			WithWaitTime[Test1](2 * time.Second),
			WithLoopTime[Test1](1 * time.Second),
			WithChanSize[Test1](1000),
			WithDoingSize[Test1](3000),
			WithHandleGoNum[Test1](3),
			WithHookFunc[Test1](func(item []Test1) bool {
				//fmt.Println("======len(item)=======", len(item))
				//fmt.Println("=============", i)
				var arpHosts []Test1
				for _, vl := range item {
					fmt.Println(vl)
				}
				//fmt.Println("len(arpHosts)",len(arpHosts))
				//fmt.Println(len(arpHosts)/3000)
				//os.Exit(2)
				//for i := 0;i < len(arpHosts)/3000;i++{
				//	var info = arpHosts[i*3000:i*3000+3000]
				if err = DB.
					Model(&Test1{}).
					//Debug().
					Create(&arpHosts).
					Error; nil != err {
					fmt.Println("=========", err)
					os.Exit(6)
				}

				//}

				return true
			}),
			WithEndHook[Test1](func(b bool, i ...Test1) {
				//fmt.Println("--------------")
				//fmt.Println(b)
				//fmt.Println(len(i))
				//fmt.Println(i)
			}),
			WithSavePanic[Test1](func() {
				if err := recover(); err != nil {
					fmt.Println(err)
					fmt.Println(string(debug.Stack()))
					os.Exit(8)
				}
			}),
		},
	}

	if err = df.InitTask(task); nil != err {
		fmt.Println(err)
		os.Exit(3)
	}

	df.Run()

	for i := 0; i < 10000; i++ {
		items := Test1{
			Name:        "zhang1San-" + strconv.Itoa(i),
			Address:     "China",
			AddrTag:     "china-bj",
			Role:        "admin",
			Test1:       "test1" + strconv.Itoa(i),
			Test2:       "test2" + strconv.Itoa(i),
			Test3:       "test3" + strconv.Itoa(i),
			Test4:       "test4" + strconv.Itoa(i),
			Test5:       "test5" + strconv.Itoa(i),
			Test6:       "test6" + strconv.Itoa(i),
			Test7:       "test7" + strconv.Itoa(i),
			Test8:       "test8" + strconv.Itoa(i),
			Test9:       0,
			CreatedTime: time.Now(),
		}
		df.Submit(&items)
	}

	df.Release()
	df.Wait()

	to := time.Now().Sub(t)
	fmt.Println("BenchmarkNewBatchHook-end", to)
}
