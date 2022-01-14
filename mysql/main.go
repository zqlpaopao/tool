package main

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	config "github.com/zqlpaopao/tool/config/src"
	mysql "github.com/zqlpaopao/tool/mysql/src"
)

const project = "project"

type CoverInfos interface {
	TabName()string
}

type CronHost struct {
	Id     int16  `json:"id" cv:"id_1"`
	Name   string `json:"name"`
	Alias  string `json:"alias"`
	Port   int32  `json:"port"`
	Remark string `json:"remark"`
}



func main(){
	Conversion()



}


func Conversion(){
	c := CronHost{
		Id:     1,
		Name:   "name",
		Alias:  "alias",
		Port:   98,
		Remark: "remark",
	}
	c1,err := mysql.GetSQL(mysql.CoverReqInfo{
		Table: "table",
		StructInfo: &c,
	})
	fmt.Println(err)
	fmt.Println(c1)
}

func GetCLinet(){

	var (
		err error
		dbClient *sql.DB
	)
	//初始化环境变量
	if err = config.Ctx.Init("CONF_DIR");nil != err{
		panic(err)
	}

	if dbClient ,err = mysql.Ctx.GetClient(project);nil != err{
		panic(err)
	}

	fmt.Println(dbClient)
}

func GetGormClient(){

	var (
		err error
		dbClient *gorm.DB
	)
	//初始化环境变量
	if err = config.Ctx.Init("CONF_DIR");nil != err{
		panic(err)
	}

	if dbClient ,err = mysql.CtxOrm.GetClient(project);nil != err{
		panic(err)
	}

	dbClient = dbClient
}