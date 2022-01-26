package start

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	config "github.com/zqlpaopao/tool/config/src"
	format "github.com/zqlpaopao/tool/format/src"
	"github.com/zqlpaopao/tool/gin-model/common"
	"github.com/zqlpaopao/tool/gin-model/module/web/middleware"
	"os"
)

func init() {
	InitEnvConfig()
}

//InitEnvConfig 初始化web配置文件
func InitEnvConfig() {
	var (
		env *viper.Viper
		err error
	)
	if env, err = config.Ctx.GetEnvConf(); nil != err {
		panic(err)
	}
	if err = env.Unmarshal(common.EnvConf); nil != err {
		panic(err)
	}
	tidyAuthUrl(common.EnvConf)
}

func tidyAuthUrl(config *common.Config) {
	for _, v := range config.Env.Web.NoAuthUrl {
		config.Env.Web.NoAuthUrlMap[v] = struct{}{}
	}
}

//InitWeb 初始化web服务
func InitWeb() {
	if err := checkArgs(); nil != err {
		panic(err)
	}
	InitGin()
}

//checkArgs 检测web配置参数
func checkArgs() error {
	if common.EnvConf.Env.Web.Host == "" || common.EnvConf.Env.Web.Port == "" {
		return common.MsgErrWebListen
	}
	return nil
}

func InitGin() {
	g := gin.Default()
	g.Use(middleware.Cors, middleware.MiddleLog, gin.Recovery())
	loadRouter(g)
	//gin.ReleaseMode
	gin.SetMode(common.EnvConf.Env.Mode)
	_ = g.SetTrustedProxies(nil)
	startListen(g)

}

//设置指定的use the X-Forwarded-For
//https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
// IPv4 地址、IPv4 CIDR、IPv6 地址或 IPv6 CIDR
//g.SetTrustedProxies([]string{"192.168.1.2"})

//startListen 启动服务
func startListen(r *gin.Engine) {
	go func() {
		format.PrintGreen(fmt.Sprintf(" Will listening and serving HTTP on %s:%s", common.EnvConf.Env.Web.Host, common.EnvConf.Env.Web.Port))
		if common.EnvConf.Env.Mode != gin.DebugMode{
			if err := endless.ListenAndServe(common.EnvConf.Env.Web.Host + ":" + common.EnvConf.Env.Web.Port, r);nil != err{
				format.PrintRed("webRun Error: " + err.Error())
				os.Exit(3)
			}
		}
		if err := r.Run(common.EnvConf.Env.Web.Host + ":" + common.EnvConf.Env.Web.Port); nil != err {
			format.PrintRed("webRun Error: " + err.Error())
			os.Exit(3)
		}
	}()
}
