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
	ErrNotHave         = errors.New("not have the task name")
	ErrTaskClosed      = errors.New("this task has been closed")
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

// SubmitModel Submit task any
type SubmitModel[T any] struct {
	TaskName string
	Data     []T
}

// SubmitItem Submit task parameters
type SubmitItem struct {
	Params []interface{}
}

// check Parameter detection
func (s *SubmitModel[T]) check() error {
	if s.TaskName == "" {
		return ERRTaskNameIsEmpty
	}
	return nil
}
