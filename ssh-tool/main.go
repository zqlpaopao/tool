package main

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	src2 "github.com/zqlpaopao/tool/ssh-tool/src"
	"golang.org/x/crypto/ssh"
	"os"
)

func main() {
	testGorm()
	//testGorm()
}

//Terminal 交互式输入
func Terminal(){
	var (
		client *src2.Client
		err error
	)

	//get sshClient
	if client, err = src2.DialWithPasswd(&src2.Config{
		Addr:   "xx.xx.xx.xx:22",
		User:   "root",
		Passwd: "@xxxxxxx",
	}); nil != err {
		panic(err)
	}
	if err = client.Terminal(&src2.TerminalConfig{
		Term :  "xterm-256color",
		Height :14800,
		Weight :14800,
		Modes  :ssh.TerminalModes{
			ssh.ECHO: 1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}}).SetStdio(os.Stdin,os.Stdout, os.Stdin).Start();nil != err{
		panic(err)
	}



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
