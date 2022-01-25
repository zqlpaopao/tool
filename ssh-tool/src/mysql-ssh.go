package src

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	sqlGorm "gorm.io/driver/mysql"
	V2 "gorm.io/gorm"
	"os"
)

type MysqlConfig struct {
	UserName, PassWd, IpPort, Dbname string
}

//MysqlClient -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (m *MysqlConfig) MysqlClient(client Client) (db *sql.DB, err error) {
	mysql.RegisterDialContext("mysql+tcp", (&ViaSSHDialer{client: client.client}).Dial)
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@mysql+tcp(%s)/%s", m.UserName, m.PassWd, m.IpPort, m.Dbname))
	return
}

//Query -- ----------------------------
//--> @Description 测试数据
//--> @Param
//--> @return
//-- ----------------------------
func Query(db *sql.DB, table string) (count int, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = db.Query("SELECT count(*) from " + table); err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&count); nil != err {
			return
		}
	}

	if err = rows.Close(); err != nil {
		return
	}
	err = db.Close()
	return
}

//GormClient -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (m *MysqlConfig) GormClient(client Client) (db *gorm.DB, err error) {
	mysql.RegisterDialContext("mysql+tcp", (&ViaSSHDialer{client: client.client}).Dial)
	return gorm.Open("mysql", fmt.Sprintf("%s:%s@mysql+tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", m.UserName, m.PassWd, m.IpPort, m.Dbname))
}

//GormClientV2 -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func (m *MysqlConfig) GormClientV2(client Client) (db *V2.DB, err error) {
	mysql.RegisterDialContext("mysql+tcp", (&ViaSSHDialer{client: client.client}).Dial)
	return V2.Open(sqlGorm.Open(fmt.Sprintf("%s:%s@mysql+tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", m.UserName, m.PassWd, m.IpPort, m.Dbname)), &V2.Config{})
}

//GormQuery -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
func GormQuery(db *gorm.DB, table string) (count int, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = db.Exec("SELECT count(*) from " + table).Rows(); nil != err {
		return
	}

	for rows.Next() {
		if err = rows.Scan(&count); nil != err {
			return
		}
	}

	if err = rows.Close(); err != nil {
		return
	}
	err = db.Close()
	return
}

//NewSSHMysqlClient -- ----------------------------
//--> @Description ssh机器和mysql地址信息
//--> @Param
//--> @return
//-- ----------------------------
func NewSSHMysqlClient(sshConfig *Config, mysqlConf *MysqlConfig) (dbClient *sql.DB, err error) {
	var (
		client *Client
		//session *ssh.Session
	)

	//get sshClient
	if client, err = DialWithPasswd(sshConfig); nil != err {
		fmt.Println(err)
		os.Exit(3)
		return
	}

	//get db
	dbClient, err = mysqlConf.MysqlClient(*client)
	fmt.Println(1, err)
	return
}

//NewSSHGormClient -- ----------------------------
//--> @Description ssh机器和mysql地址信息
//--> @Param
//--> @return
//-- ----------------------------
func NewSSHGormClient(sshConfig *Config, mysqlConf *MysqlConfig) (dbClient *gorm.DB, err error) {
	var (
		client *Client
	)

	//get sshClient
	if client, err = DialWithPasswd(sshConfig); nil != err {
		return
	}

	//get db
	dbClient, err = mysqlConf.GormClient(*client)
	return
}
