package main

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	src2 "github.com/zqlpaopao/tool/ssh-tool/src"
)

func main() {
	testGorm()
	//testGorm()
}

//-- ----------------------------
//--> @Description 正常的mysql连接
//--> @Param
//--> @return
//-- ----------------------------
func testDb() {
	var (
		db  *sql.DB
		err error
	)

	if db, err = src2.NewSSHMysqlClient(&src2.Config{
		Addr:   "xx.xx.xx.xx:22",
		User:   "xxx",
		Passwd: "@xxxxx",
	}, &src2.MysqlConfig{
		UserName: "xxx",
		PassWd:   "^xxxxxxx",
		IpPort:   "xx.xx.xxx.xxx:3306",
		Dbname:   "xxxxxxxxx",
	}); nil != err {
		panic(err)
	}
	fmt.Println(1)
	db = db

}
//-- ----------------------------
//--> @Description GORM连接
//--> @Param
//--> @return
//-- ----------------------------
func testGorm() {

	var (
		db            *gorm.DB
		err           error
	)
	if db, err = src2.NewSSHGormClient(&src2.Config{
		Addr:   "xx.xx.xx.xx:22",
		User:   "xxx",
		Passwd: "@xxxxx",
	}, &src2.MysqlConfig{
		UserName: "xxx",
		PassWd:   "^xxxxxxx",
		IpPort:   "xx.xx.xxx.xxx:3306",
		Dbname:   "xxxxxxxxx",
	}); nil != err {
		panic(err)
	}
	db = db

}
