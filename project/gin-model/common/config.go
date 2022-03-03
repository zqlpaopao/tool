package common

var EnvConf = &Config{
	Env: EnvConfig{
		Web:   WebConfig{
			NoAuthUrl:    []string{},
			NoAuthUrlMap: map[string]struct{}{},
		},
		PProf: Pprof{},
		Mode : "",
	},
}

var LogConf = &LogConfig{}