package main

import (
	"database/sql"
	"fmt"
	config "github.com/zqlpaopao/tool/config/src"
	mysql "github.com/zqlpaopao/tool/mysql/src"
)

const project = "project"

func main(){



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