package pkg

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	//"github.com/zqlpaopao/tool/batch_select/pkg"
	"testing"
	"time"
)

/*
$ go test -bench=.  --benchmem
goos: darwin
goarch: arm64
pkg: github.com/xx/tool/batch_select/pkg
BenchmarkNewBatchSelect-10      1000000000               0.5782 ns/op          0 B/op          0 allocs/op
PASS
ok      github.com/xx/tool/batch_select/pkg      17.588s

*/

func BenchmarkNewBatchSelect(t *testing.B) {
	type Info struct {
		Id         int64     `json:"id"`
		CreateTime time.Time `json:"create_time"`
	}

	var (
		db  *gorm.DB
		err error
	)

	if db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local&parseTime=true&allowNativePasswords=true",
		"root", "meimima123", "127.0.0.1", "3306", "joypaw_base")), &gorm.Config{}); nil != err {
		panic(err)
	}

	b := NewBatchSelect[Info]()

	if err = b.InitTask(InitTaskModel[Info]{
		TaskName: "sss",
		Opt: []OptionInter[Info]{
			WithHandleGoNum[Info](100), //处理的协程数量
			//WithDebug(true),                 //是否开启debug
			WithLimit[Info](10000),                //limit 的个数 默认10000
			WithOrderColumn[Info]("dev_ip_int32"), //要进行取舍的列
			WithTable[Info]("device"),
			WithSqlWhere[Info]("dev_ip_int32 > 0", nil),  //where 条件
			WithResChanSize[Info](10000),                 //接受数据的chan大小
			WithMysqlSqlCli[Info](db),                    //接受数据的chan大小
			WithSelectFiled[Info]("dev_ip_int32,dev_ip"), //接受数据的chan大小
			WithCallFunc[Info](func(res *[]Info) {
				println(res)

			}), //接受数据

		},
	}); nil != err {
		panic(err)
	}

	if err = b.Run("xxx"); nil != err {
		panic(err)
	}

	b.Wait()
}
