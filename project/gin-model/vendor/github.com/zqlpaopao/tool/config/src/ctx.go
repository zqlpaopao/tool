package config

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/zqlpaopao/tool/format/src"
	"os"
	"sync"
)

const (
	defaultLogPath = "./log/%Y%m%d-api.log"
)

type (
	confManager struct {
		sync.Once
		mu     sync.RWMutex
		path   string
		vipers map[string]*viper.Viper
	}
)

var (
	Ctx = new(confManager)
)

func init() {
	Ctx.Once.Do(func() {
		Ctx.vipers = make(map[string]*viper.Viper)
	})
}

func (c *confManager) Init(envItem string) error {
	confPath := os.Getenv(envItem)
	return c.InitByPath(confPath)
}

func (c *confManager) InitByPath(confPath string) error {
	if confPath == "" {
		return errors.New("ConfigManager.Init: config path is empty")
	}
	c.path = confPath
	return nil
}

func (c *confManager) getViper(configName string) (*viper.Viper, error) {
	c.mu.RLock()
	viperOb, exists := c.vipers[configName]
	c.mu.RUnlock()
	if exists {
		return viperOb, nil
	}

	viperOb = viper.New()
	viperOb.SetConfigType("yaml")
	viperOb.AddConfigPath(c.path)
	viperOb.SetConfigName(configName)
	err := viperOb.ReadInConfig()
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	c.vipers[configName] = viperOb
	c.mu.Unlock()
	return viperOb, nil
}

func (c *confManager) GetKafkaConf() (*viper.Viper, error) {
	return c.getViper("kafka")
}

func (c *confManager) GetRocketConf() (*viper.Viper, error) {
	return c.getViper("rocket")
}

func (c *confManager) GetOpenSearchConf() (*viper.Viper, error) {
	return c.getViper("opensearch")
}

func (c *confManager) GetRedisConf() (*viper.Viper, error) {
	return c.getViper("redis")
}

func (c *confManager) GetZkConf() (*viper.Viper, error) {
	return c.getViper("zk")
}

func (c *confManager) GetMysqlConf() (*viper.Viper, error) {
	return c.getViper("mysql")
}

func (c *confManager) GetLogsConf() (*viper.Viper, error) {
	return c.getViper("logs")
}

func (c *confManager) GetRpcConf() (*viper.Viper, error) {
	return c.getViper("rpc")
}

func (c *confManager) GetAlarmConf() (*viper.Viper, error) {
	return c.getViper("alarm")
}

func (c *confManager) GetEnvConf() (*viper.Viper, error) {
	return c.getViper("env")
}

func (c *confManager) GetEtcdConf() (*viper.Viper, error) {
	return c.getViper("etcd")
}

func (c *confManager) GetMnsConf() (*viper.Viper, error) {
	return c.getViper("mns")
}

func (c *confManager) GetCronConf() (*viper.Viper, error) {
	return c.getViper("cron")
}

func (c *confManager) GetLogLayout(project string, file string) (logPath string, err error) {
	viperOb, err := c.GetLogsConf()
	if err != nil {
		return
	}
	if file == "" {
		file = "system"
	}
	logPath = defaultLogPath
	logPathFromConfig := viperOb.GetString(project + "." + file)
	if logPathFromConfig != "" {
		logPath = logPathFromConfig
	} else {
		src.PrintRed("warning: log_path is empty, use default log_path (" + logPath + ")")
	}
	return
}
