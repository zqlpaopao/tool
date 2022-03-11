package pkg

import "time"

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 5 * time.Second
	idleTimeout  = 60 * time.Second
	poolTimeOut  = 60 * time.Second
	poolSize     = 20
)
