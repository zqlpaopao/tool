package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	config "github.com/zqlpaopao/tool/config/src"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	mySQLManager struct {
		sync.Once
		mu      sync.RWMutex
		clients map[string]*sql.DB
	}
)

var (
	Ctx = new(mySQLManager)
	maxIdleConn, maxOpenConn        = "5", "20"
)

func init() {
	Ctx.Once.Do(func() {
		Ctx.clients = make(map[string]*sql.DB)
	})
}

//-- ----------------------------
//--> @Description  读取配置
//--> @Param 配置信息
//--> @return 错误
//-- ----------------------------
func (manager *mySQLManager) readOption(name string) (map[string]string, error) {
	mysqlViper, err := config.Ctx.GetMysqlConf()
	if err != nil {
		return nil, err
	}
	dbConfig := mysqlViper.GetStringMapString(name)
	if dbConfig["source"] == "" {
		return nil, fmt.Errorf("%q mysql address is empty", name)
	}

	return checkArgs(dbConfig), nil
}

//-- ----------------------------
//--> @Description  校验配置
//--> @Param 配置信息
//--> @return 配置信息
//-- ----------------------------
func checkArgs(dbConfig map[string]string) map[string]string {
	var (
		maxIdleConnectionInt, maxOpenConnectionInt int
		err                                        error
	)
	if maxIdleConnections, exit := dbConfig["max_idle_connections"]; exit {
		if maxIdleConnectionInt, err = strconv.Atoi(maxIdleConnections); err != nil {
			dbConfig["max_idle_connections"] = maxIdleConn
		}
	}
	if maxOpenConnections, exit := dbConfig["max_open_connections"]; exit {
		if maxOpenConnectionInt, err = strconv.Atoi(maxOpenConnections); err != nil {
			dbConfig["max_open_connections"] = maxOpenConn
		}
	}
	if maxIdleConnectionInt > maxOpenConnectionInt {
		dbConfig["max_idle_connections"] = dbConfig["max_open_connections"]
	}
	return dbConfig
}

//-- ----------------------------
//--> @Description  创建新连接
//--> @Param 连接实例
//--> @return 错误
//-- ----------------------------
func (manager *mySQLManager) createClient(name string) (*sql.DB, error) {
	dataConfig, err := manager.readOption(name)
	if err != nil {
		return nil, err
	}
	dataSourceName := dataConfig["source"]
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	maxIdleCons := 10
	if maxIdleConnections, exit := dataConfig["max_idle_connections"]; exit {
		if maxIdleConnections, err := strconv.Atoi(maxIdleConnections); err == nil {
			maxIdleCons = maxIdleConnections
		}
	}
	db.SetMaxIdleConns(maxIdleCons)
	maxOpenCons := 10
	if maxOpenConnections, exit := dataConfig["max_open_connections"]; exit {
		if maxOpenConnections, err := strconv.Atoi(maxOpenConnections); err == nil {
			maxOpenCons = maxOpenConnections
		}
	}
	db.SetMaxOpenConns(maxOpenCons)
	connMaxLifeTime := 10 * time.Minute
	if connectionsMaxLifeTime, exit := dataConfig["connections_max_life_time"]; exit {
		if connectionsMaxLifeTime, err := strconv.Atoi(connectionsMaxLifeTime); err == nil {
			connMaxLifeTime = time.Duration(connectionsMaxLifeTime) * time.Second
		}
	}
	db.SetConnMaxLifetime(connMaxLifeTime)
	return db, err
}

//GetClient -- ----------------------------
//--> @Description  获取连接
//--> @Param 链接实例名称
//--> @return 错误
//-- ----------------------------
func (manager *mySQLManager) GetClient(name string) (*sql.DB, error) {
	manager.mu.RLock()
	client, exist := manager.clients[name]
	manager.mu.RUnlock()
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
	return client, nil
}

//Release -- ----------------------------
//--> @Description  释放连接
//--> @Param 链接实例名称
//--> @return 无
//-- ----------------------------
func (manager *mySQLManager) Release(name string) {
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
func (manager *mySQLManager) ReleaseAll() {
	manager.mu.Lock()
	for name, client := range manager.clients {
		_ = client.Close()
		delete(manager.clients, name)
	}
	manager.mu.Unlock()
}

//IsMySqlNilError -- ----------------------------
//--> @Description  是否为空错误
//--> @Param mysql错误
//--> @return 错误是否为nil错误
//-- ----------------------------
func IsMySqlNilError(err error) bool {
	return err == sql.ErrNoRows
}

//IsMysqlDuplicateUniqueError -- ----------------------------
//--> @Description  是否为唯一索引重复错误
//--> @Param mysql错误
//--> @return 错误是否为唯一索引重复错误
//-- ----------------------------
func IsMysqlDuplicateUniqueError(err error) bool {
	if strings.Contains(err.Error(), "1062") && strings.Contains(err.Error(), "Duplicate") {
		return true
	}
	return false
}
