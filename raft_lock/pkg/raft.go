package pkg

import (
	"errors"
	"time"
)

type Raft struct {
	opt *RedisOption
	err error
	rdb RDbClient
}

func (r *Raft) Lock(s string, duration time.Duration) (int64, error) {
	return r.rdb.Lock(s, duration)
}

func (r *Raft) UnLock(lockName string) (int64, error) {
	return r.rdb.UnLock(lockName)
}

func (r *Raft) Renewal(s string, duration time.Duration) (int64, error) {
	return r.rdb.Renewal(s, duration)
}

// NewRaft -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func NewRaft(opt *RedisOption) *Raft {
	return &Raft{
		opt: opt,
		rdb: NewRdbCli(opt),
	}
}

// Check -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (r *Raft) Check() error {
	if len(r.opt.addr) < 1 {
		return errors.New("addr err")
	}
	if r.opt.groupName == "" {
		return errors.New("group err")
	}
	return nil
}

// Ping -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (r *Raft) Ping() error {
	return r.rdb.Ping()
}

// SetError -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (r *Raft) SetError(err error) {
	r.err = err
}

// Error -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (r *Raft) Error() error {
	return r.err
}
