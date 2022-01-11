package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	config "github.com/zqlpaopao/tool/config/src"
	"strconv"
	"sync"
	"time"
)

type (
	myGormSQLManager struct {
		sync.Once
		mu      sync.Mutex
		clients map[string]*gorm.DB
	}
)

var (
	CtxOrm = new(myGormSQLManager)
)

func init() {
	CtxOrm.Once.Do(func() {
		CtxOrm.clients = make(map[string]*gorm.DB)
	})
}

/**
 * 读取配置
 * @params name	连接名称
 * @return	配置信息
 * @return	错误
 */
func (manager *myGormSQLManager) readOption(name string) (map[string]string, error) {
	mysqlViper, err := config.Ctx.GetMysqlConf()
	if err != nil {
		return nil, err
	}
	dbConfig := mysqlViper.GetStringMapString(name)
	if dbConfig["source"] == "" {
		return nil, fmt.Errorf("%q mysql source is empty", name)
	}
	return checkArgs(dbConfig), nil
}

//-- ----------------------------
//--> @Description  创建新连接
//--> @Param 连接实例
//--> @return 错误
//-- ----------------------------
func (manager *myGormSQLManager) createClient(name string) (*gorm.DB, error) {
	var (
		err   error
		dbIns *gorm.DB
	)
	dbConfig, err := manager.readOption(name)
	if err != nil {
		return nil, err
	}
	if dbIns, err = gorm.Open("mysql", dbConfig["source"]); err != nil {
		return nil, err
	}
	// 最大空闲链接
	maxIdleCons := 10
	if maxIdleConnections, exist := dbConfig["max_idle_connections"]; exist {
		if maxIdleConnections, err := strconv.Atoi(maxIdleConnections); err == nil {
			maxIdleCons = maxIdleConnections
		}
	}
	dbIns.DB().SetMaxIdleConns(maxIdleCons)

	// 最大打开连接数
	maxOpenCons := 20
	if maxOpenConnections, exist := dbConfig["max_open_connections"]; exist {
		if maxOpenConnections, err := strconv.Atoi(maxOpenConnections); err == nil {
			maxOpenCons = maxOpenConnections
		}
	}
	dbIns.DB().SetMaxOpenConns(maxOpenCons)

	// 连接最大生命周期
	maxLifeTime := 10 * time.Minute
	if maxLifeTimeCons, exist := dbConfig["connections_max_life_time"]; exist {
		if maxLifeTimeCons, err := strconv.Atoi(maxLifeTimeCons); err == nil {
			maxLifeTime = time.Duration(maxLifeTimeCons) * time.Second
		}
	}
	dbIns.DB().SetConnMaxLifetime(maxLifeTime)
	return dbIns, nil

}

//GetClient -- ----------------------------
//--> @Description  获取连接
//--> @Param 链接实例名称
//--> @return 错误
//-- ----------------------------
func (manager *myGormSQLManager) GetClient(name string) (*gorm.DB, error) {
	manager.mu.Lock()
	client, exist := manager.clients[name]
	manager.mu.Unlock()
	// if client not exist
	if !exist {
		newClient, err := manager.createClient(name)
		if err != nil {
			return nil, err
		}
		manager.mu.Lock()
		if client, exist = manager.clients[name]; !exist {
			manager.clients[name] = newClient
			client = newClient
		}
		manager.mu.Unlock()
		if client != newClient {
			_ = newClient.Close()
		}
	}
	mysqlViper, err := config.Ctx.GetMysqlConf()
	if err == nil {
		if mysqlViper.GetBool("debug") {
			client = client.LogMode(true).Debug()
		}
	}
	return client, nil
}

//Release -- ----------------------------
//--> @Description  释放连接
//--> @Param 链接实例名称
//--> @return 无
//-- ----------------------------
func (manager *myGormSQLManager) Release(name string) {
	manager.mu.Lock()
	if client, ok := manager.clients[name]; ok {
		_ = client.Close()
		delete(manager.clients, name)
	}
	manager.mu.Unlock()
}

//ReleaseAll -- ----------------------------
//--> @Description  释放所有连接
//--> @Param 链接实例名称
//--> @return 无
//-- ----------------------------
func (manager *myGormSQLManager) ReleaseAll() {
	manager.mu.Lock()
	for name, client := range manager.clients {
		_ = client.Close()
		delete(manager.clients, name)
	}
	manager.mu.Unlock()
}

