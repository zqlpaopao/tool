package pkg

import (
	"errors"
	"time"
)

const (

	// ReadTimeout the maximum time for a single read after the connection is established
	ReadTimeout time.Duration = 1 << 4

	//WriteTimeout the maximum time for a single write after the connection is established
	WriteTimeout time.Duration = 1 << 4

	//DialTimeout  maximum connection establishment time for connections
	DialTimeout time.Duration = 1 << 4

	//MaxConnNum pool default value
	MaxConnNum = 2 << 2

	//MaxIdleNum Maximum idle connection
	MaxIdleNum

	//CheckNum  check client de lifetime
	CheckNum

	//MaxLifeTime Maximum effective lifecycle of the connection
	MaxLifeTime time.Duration = 60 * time.Second

	//MaxIdleTime Maximum effective lifecycle of the connection
	MaxIdleTime = 60 * time.Second

	//HeartbeatTime Heartbeat detection time
	HeartbeatTime = 10 * time.Second

	//IsUsing Is cli in use
	IsUsing = 1
	//IsFree cli is free
	IsFree = 2
)

type NetWork string

const (
	Tcp = "tcp"
	Udp = "udp"
)

var (
	ErrAddrEmpty       = errors.New("addr is empty")
	ErrSshClientConfig = errors.New("ssh config parameter error")

	//pool

	ErrPoolTag   = errors.New("pool tag is empty")
	ErrNotExist  = errors.New("cli is not  exist")
	ErrItemIsNil = errors.New("item is nil")
)
