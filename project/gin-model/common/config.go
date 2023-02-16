package common

var (
	EnvConf = &Config{
		Env: EnvConfig{
			Web: WebConfig{
				NoAuthUrl:    []string{},
				NoAuthUrlMap: map[string]struct{}{},
			},
			PProf: Pprof{},
			Mode:  "",
		},
	}
	LogConf = &LogConfig{}
	//MysqlInfo = &MysqlEnvInfo{}
	Urls = &ReqAddress{}
)
