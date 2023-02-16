package start

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	config "github.com/zqlpaopao/tool/config/src"
	"github.com/zqlpaopao/tool/project/gin-model/common"
)

// InitUrlsConfig -- --------------------------
// --> @Describe 初始化请求地址
// --> @params
// --> @return
// -- ------------------------------------
func InitUrlsConfig() {
	var (
		env *viper.Viper
		err error
	)
	if env, err = config.Ctx.GetUrlConf(); nil != err {
		panic(err)
	}
	if err = env.Unmarshal(common.Urls, func(decoderConfig *mapstructure.DecoderConfig) {
		decoderConfig.TagName = "json"
	}); nil != err {
		panic(err)
	}
}
