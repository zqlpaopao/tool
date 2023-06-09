package pkg

import (
	"errors"
	"time"
)

const (
	DefaultIsRunning     = 1
	DefaultIsClose       = 2
	DefaultPoolSize      = 3
	DefaultInitPoolSize  = 3
	DefaultIdleTimeout   = 10 * time.Second
	DefaultCheckInterval = 3 * time.Second
)

var (
	ErrClosed        = errors.New("pool is closed")
	ErrFactoryIssNil = errors.New("factory func is nil")
	ErrCloseIssNil   = errors.New("close func is nil")
	ErrPingIssNil    = errors.New("ping func is nil")
)
