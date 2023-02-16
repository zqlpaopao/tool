package start

import (
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	config "github.com/zqlpaopao/tool/config/src"
	"github.com/zqlpaopao/tool/project/gin-model/common"
	log "github.com/zqlpaopao/tool/zap-log/src"
	"os"
)

func init() {
	app := &cli.App{
		Name:  "DCos",
		Usage: "DCos controller",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "config path",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			conf := c.String("conf")
			return Init(conf)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

// InitLog 初始化日志模块
func InitLog() {
	var (
		logCtx *viper.Viper
		err    error
	)
	if logCtx, err = config.Ctx.GetLogsConf(); nil != err {
		panic(err)
	}
	if err = logCtx.Unmarshal(common.LogConf); nil != err {
		panic(err)
	}

	//初始化日志
	log.InitLoggerHandler(
		log.NewLogConfig(
			log.InitInfoPathFileName(common.LogConf.Log.InfoPath),
			log.InitWarnPathFileName(common.LogConf.Log.WarnPath),
			log.InitWithMaxAge(common.LogConf.Log.MaxAge),
			log.InitWithRotationCount(common.LogConf.Log.RotationCount),
			log.InitWithRotationTime(common.LogConf.Log.RotationTime),
			log.InitWithIp(1),
			log.InitBufferSize(common.LogConf.Log.BufferSize)))

	//初始化异步日志
	log.NewAsyncLogConfig(
		log.InitLogAsyncBuffSize(common.LogConf.Log.AsyncBuffSize),
		log.InitLogAsyncGoNum(common.LogConf.Log.AsyncGoNum))
}
