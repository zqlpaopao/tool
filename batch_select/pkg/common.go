package pkg

import (
	"errors"
)

const (
	Limit       = 10000
	HandleGoNum = 10
	ChanSize    = 10000
)

var (
	ERRTaskNameIsEmpty          = errors.New("task name is empty")
	ERRMySqlCli                 = errors.New("mysql cli is empty")
	ERROrderColumn              = errors.New("order column is empty")
	ErrNotHave                  = errors.New("not have the task name")
	ErrSelectNotHaveOrderColumn = errors.New("the query column has no columns to sort")
)

//InitTaskModel Initialize task model
type InitTaskModel struct {
	TaskName string
	Opt      []Option
}

type MinMaxInfo struct {
	MinId string
	MaxId string
}

//check Parameter detection
func (s *InitTaskModel) check() error {
	if s.TaskName == "" {
		return ERRTaskNameIsEmpty
	}
	return nil
}
