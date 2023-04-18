package main

import (
	"fmt"
	"github.com/zqlpaopao/tool/batch_select/pkg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

type D1 struct {
	ID         uint64    `json:"id" gorm:"column:id"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
}

func (m *D1) TableName() string {
	return "xxx"
}

func main() {
	var (
		db   *gorm.DB
		err  error
		info []D1
		lock = &sync.Mutex{}
	)

	if db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local&parseTime=true&allowNativePasswords=true",
		"root", "xxx", "127.0.0.1", "3306", "xxx")), &gorm.Config{}); nil != err {
		panic(err)
	}

	b := pkg.NewBatchSelect[D1]()

	if err = b.InitTask(pkg.InitTaskModel[D1]{
		TaskName: "d1",
		Opt: []pkg.OptionInter[D1]{
			pkg.WithHandleGoNum[D1](100),  //处理的协程数量
			pkg.WithDebug[D1](true),       //是否开启debug
			pkg.WithLimit[D1](10000),      //limit 的个数 默认10000
			pkg.WithOrderColumn[D1]("id"), //where 走的索引
			pkg.WithTable[D1]("xxx"),
			pkg.WithSqlWhere[D1]("id > 0", nil), //where 条件
			pkg.WithResChanSize[D1](10000),      //接受数据的chan大小
			pkg.WithMysqlSqlCli[D1](db),         //接受数据的chan大小
			pkg.WithSelectFiled[D1]("*"),        //接受数据的chan大小
			pkg.WithCallFunc[D1](func(res *[]D1) {
				lock.Lock()
				info = append(info, *res...)
				lock.Unlock()

			}), //接受数据

		},
	}); nil != err {
		panic(err)
	}

	if err = b.Run("d1"); nil != err {
		panic(err)
	}

	b.Wait()
	fmt.Println(len(info))

}
