package pkg

import "time"

const (
	//RetryCount retry times
	RetryCount = 3
	//RetryInterval retry interval
	RetryInterval = 3 * time.Second
)
