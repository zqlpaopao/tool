package start

import (
	config "github.com/zqlpaopao/tool/config/src"
)

func init() {

}

// Init -- --------------------------
// --> @Describe 初始化服务
// --> @params
// --> @return
// -- ------------------------------------
func Init(conf string) (err error) {
	if err := config.Ctx.InitOsArgs(conf); nil != err {
		panic(err)
	}
	InitLog()
	if err = InitEnvConfig(); nil != err {
		return
	}
	//InitMysqlConfig()
	InitUrlsConfig()
	//mysql.InitDB()
	return
}
