package pkg

import (
	"time"
)

type Cmder interface {
	Check() error
	SetError(err error)
	Error() error
	RDbClient
}

// RDbClient interface
type RDbClient interface {
	Ping() error
	Lock(string, time.Duration) (int64, error)
	UnLock(string2 string) (int64, error)
	Renewal(string, time.Duration) (int64, error)
	GetLockInfo() (map[string]string, error)
}
