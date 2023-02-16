package start

import (
	"context"
	"errors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	config "github.com/zqlpaopao/tool/config/src"
	format "github.com/zqlpaopao/tool/format/src"
	"github.com/zqlpaopao/tool/project/gin-model/common"
	"github.com/zqlpaopao/tool/project/gin-model/module/web/middleware"
	versionInfo "github.com/zqlpaopao/tool/version-num-manager/src"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	InitEnvConfig()
}

// InitEnvConfig 初始化web配置文件
func InitEnvConfig() (err error) {
	var (
		env *viper.Viper
	)
	if env, err = config.Ctx.GetEnvConf(); nil != err {
		return
	}
	if err = env.Unmarshal(common.EnvConf); nil != err {
		return
	}
	tidyAuthUrl(common.EnvConf)
	return err
}

func tidyAuthUrl(config *common.Config) {
	for _, v := range config.Env.Web.NoAuthUrl {
		config.Env.Web.NoAuthUrlMap[v] = struct{}{}
	}
}

// InitWeb 初始化web服务
func InitWeb() {
	if err := checkArgs(); nil != err {
		panic(err)
	}
	InitGin()
}

// checkArgs 检测web配置参数
func checkArgs() error {
	if common.EnvConf.Env.Web.Host == "" || common.EnvConf.Env.Web.Port == "" {
		return common.MsgErrWebListen
	}
	return nil
}

func InitGin() {
	g := gin.Default()
	g.Use(middleware.Cors, middleware.MiddleLog, gin.Recovery(), middleware.InitContext())
	loadRouter(g)
	//gin.ReleaseMode
	gin.SetMode(common.EnvConf.Env.Mode)
	_ = g.SetTrustedProxies(nil)
	openPProf(g)
	openVersionInfo()
	startListen(g)

}

// 开启pprof
func openPProf(g *gin.Engine) {
	if common.EnvConf.Env.PProf.OpenTag {
		pprof.Register(g) // 性能
	}
}

// 开启版本信息
func openVersionInfo() {
	if err := versionInfo.NewVersionNumManager(
		versionInfo.WithNotAuth(false),
		versionInfo.WithBranch(true),
		versionInfo.WithPrint(true),
		versionInfo.WithTag("Version Info "),
	).Do().Error(); err != nil {
		format.PrintRed(err.Error())
	}
}

// 设置指定的use the X-Forwarded-For
// https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
// IPv4 地址、IPv4 CIDR、IPv6 地址或 IPv6 CIDR
// g.SetTrustedProxies([]string{"192.168.1.2"})
// startListen 启动服务
func startListen(r *gin.Engine) {
	srv := &http.Server{
		Addr:    common.EnvConf.Env.Web.Host + ":" + common.EnvConf.Env.Web.Port,
		Handler: r,
	}
	go func() {
		format.PrintGreen("Listen on  " + common.EnvConf.Env.Web.Host + ":" + common.EnvConf.Env.Web.Port)
		if err := srv.ListenAndServe(); nil != err && errors.Is(err, http.ErrServerClosed) {
			format.PrintRed("webRun Error: " + err.Error())
			if common.EnvConf.Env.Mode == gin.DebugMode {
				os.Exit(3)
			}
		}
	}()
	Shutdown(srv)
}

// Shutdown -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func Shutdown(srv *http.Server) {
	// wait for interrupt signal to gracefully shut down the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	format.PrintRed("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		format.PrintRed("Server Shutdown:" + err.Error())
		if common.EnvConf.Env.Mode == gin.DebugMode {
			os.Exit(3)
		}
	}
	format.PrintRed("Server exiting")

	select {
	case <-ctx.Done():
		format.PrintRed("Server exited ...")
		os.Exit(1)
	}
}
