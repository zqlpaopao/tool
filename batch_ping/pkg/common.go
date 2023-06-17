package pkg

import (
	"errors"
	"time"
)

const (
	DoingSize   = 1000
	HandleGoNum = 10
	ChanSize    = 1000
	WaitTime    = 2 * time.Second
	LoopTime    = 1 * time.Second
)

const (
	// OPENED represents that the pool is opened.
	OPENED = iota
	// CLOSED represents that the pool is closed.
	CLOSED
)

var (
	ERRTaskNameIsEmpty = errors.New("task name is empty")
	ERRHookFuncIsEmpty = errors.New("hookFunc is empty")
)

// InitTaskModel Initialize task any
type InitTaskModel[T any] struct {
	TaskName string
	Opt      []Option[T]
}

// check Parameter detection
func (s *InitTaskModel[T]) check() error {
	if s.TaskName == "" {
		return ERRTaskNameIsEmpty
	}
	return nil
}
