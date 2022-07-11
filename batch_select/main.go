package main

import (
	"encoding/json"
	"fmt"
	"github.com/zqlpaopao/tool/batch_select/pkg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type D1 struct {
	Id  int64     `json:"id"`
	CreateTime  time.Time `json:"create_time"`
}

type D2 struct {
	Id         uint32    `json:"id"`
	CreateTime time.Time `json:"create_time"`
}
type D3 struct {
	ItemId  int64     `json:"item_id"`
	CreateTime time.Time `json:"create_time"`
}

func main() {
	var (
		db  *gorm.DB
		err error
	)

	if db, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local&parseTime=true&allowNativePasswords=true",
		"root", "xxxx", "127.0.0.1", "3306", "xxxxx")), &gorm.Config{}); nil != err {
		panic(err)
	}

	b := pkg.NewBatchSelect()

	infos := []D1{}
	infos1 := []D2{}
	infos2 := []D3{}
	if err = b.InitTask([]pkg.InitTaskModel{
		{
			TaskName: "d1",
			Opt:      []pkg.Option{
				pkg.WithHandleGoNum(100),//处理的协程数量
				pkg.WithDebug(true),//是否开启debug
				pkg.WithLimit(10000),//limit 的个数 默认10000
				pkg.WithOrderColumn("id"),//要进行取舍的列
				pkg.WithTable("d1"),
				pkg.WithSqlWhere("id > 0"),//where 条件
				pkg.WithResChanSize(10000),//接受数据的chan大小
				pkg.WithMysqlSqlCli(db),//接受数据的chan大小
				pkg.WithOrderId(true),//是否走主键id
				pkg.WithSelectFiled("d1,d2"),//接受数据的chan大小
				pkg.WithCallFunc(func(res *[]map[string]interface{}) {
					fmt.Println("res",len(*res))
					//fmt.Println("res",*res)
					var (
						bs []byte
						info []D1
					)
					if bs,err = json.Marshal(res);nil != err{
						panic(err)
					}

					if err = json.Unmarshal(bs,&info);nil!= err{
						panic(err)
					}
					fmt.Println("len(infos1)",len(info))

					//fmt.Println(info)
					infos = append(infos,info...)
					fmt.Println("len(infos)",len(infos))


				}),//接受数据

			},
		},
		{
			TaskName: "d2",
			Opt:      []pkg.Option{
				pkg.WithHandleGoNum(100),//处理的协程数量
				pkg.WithDebug(true),//是否开启debug
				pkg.WithLimit(10000),//limit 的个数 默认10000
				pkg.WithOrderColumn("id"),//要进行取舍的列
				pkg.WithTable("d2"),
				pkg.WithSqlWhere("id > 0"),//where 条件
				pkg.WithResChanSize(10000),//接受数据的chan大小
				pkg.WithMysqlSqlCli(db),//接受数据的chan大小
				pkg.WithSelectFiled("id,d2"),//接受数据的chan大小
				pkg.WithOrderId(true),//是否走主键id
				pkg.WithCallFunc(func(res *[]map[string]interface{}) {
					fmt.Println("res-dd",len(*res))
					//fmt.Println("res-cmdb",*res)
					var (
						bs []byte
						info []D2
					)
					if bs,err = json.Marshal(res);nil != err{
						panic(err)
					}

					if err = json.Unmarshal(bs,&info);nil!= err{
						panic(err)
					}
					fmt.Println("len(infos1)-dd",len(info))

					//fmt.Println(info)
					infos1 = append(infos1,info...)
					fmt.Println("len(infos)-dd",len(infos))


				}),//接受数据

			},
		},

		{
			TaskName: "d3",
			Opt:      []pkg.Option{
				pkg.WithHandleGoNum(100),//处理的协程数量
				pkg.WithDebug(true),//是否开启debug
				pkg.WithLimit(10000),//limit 的个数 默认10000
				pkg.WithOrderColumn("item_id"),//要进行取舍的列
				pkg.WithTable("d3"),
				pkg.WithSqlWhere("item_id > 0"),//where 条件
				pkg.WithResChanSize(10000),//接受数据的chan大小
				pkg.WithMysqlSqlCli(db),//接受数据的chan大小
				pkg.WithSelectFiled("item_id,item_id1"),//接受数据的chan大小
				pkg.WithOrderId(true),//是否走主键id
				pkg.WithCallFunc(func(res *[]map[string]interface{}) {
					fmt.Println("res-item_id",len(*res))
					//fmt.Println("res-cmdb",*res)
					var (
						bs []byte
						info []D3
					)
					if bs,err = json.Marshal(res);nil != err{
						panic(err)
					}

					if err = json.Unmarshal(bs,&info);nil!= err{
						panic(err)
					}
					fmt.Println("len(infos1)-item_id",len(info))

					//fmt.Println(info)
					infos2 = append(infos2,info...)
					fmt.Println("len(infos)-item_id",len(infos))


				}),//接受数据

			},
		},
	}...); nil != err {
		panic(err)
	}

	//if err = b.Run([]string{"d3"}...);nil != err{
	if err = b.Run([]string{"d3","d1","d2"}...);nil != err{
		panic(err)
	}

	b.Wait()
	fmt.Println(len(infos))
	fmt.Println(len(infos1))
	fmt.Println(len(infos2))


}
