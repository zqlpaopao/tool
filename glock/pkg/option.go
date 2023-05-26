package pkg

import (
	"github.com/go-redis/redis/v8"
	randString "github.com/zqlpaopao/tool/rand-string/pkg"
	"time"
)

type RenewalTypeFunc func(executeCount uint) uint

//LockCallbackFun Retry the function that completed execution
type LockCallbackFun func(...interface{})

type Glock interface {
	Lock(...interface{}) Glock
	UnLock() Glock
	IsMaster() bool
	Error() error
	GetMembers() (map[string]string, error)
}

type Option interface {
	apply(opt *option)
}

type option struct {
	seizeClose      chan struct{}
	renewalTag      chan struct{}
	seizeTag        bool
	seizeCycle      time.Duration
	expire          uint
	redisTimeout    time.Duration
	key             string
	masterKey       string
	RenewalOften    RenewalTypeFunc
	redisClient     *redis.Client
	lockFailFunc    LockCallbackFun
	lockSuccessFunc LockCallbackFun
}

//OpFunc type func
type OpFunc func(*option)

//NewOptions make option
func NewOptions(f ...Option) *option {
	return clone().WithOptions(f...)
}

//apply assignment function entity
func (o OpFunc) apply(opt *option) {
	o(opt)
}

//clone  new object
func clone() *option {
	return &option{
		seizeClose:      make(chan struct{}),
		renewalTag:      make(chan struct{}),
		seizeTag:        false,
		seizeCycle:      DefaultSeizeTIme,
		expire:          DefaultExpireTIme,
		redisTimeout:    DefaultRedisTimeOut,
		key:             randString.RandGenString(randString.RandSourceLetterAndNumber, 8),
		masterKey:       Lock,
		RenewalOften:    DefaultRenewalTime,
		redisClient:     &redis.Client{},
		lockSuccessFunc: nil,
		lockFailFunc:    nil,
	}
}

//WithOptions Execute assignment function entity
func (o *option) WithOptions(f ...Option) *option {
	for _, v := range f {
		v.apply(o)
	}
	return o
}

//WithMasterKey Set master Key .default Lock
func WithMasterKey(key string) OpFunc {
	return func(o *option) {
		o.masterKey = key
	}
}

//WithSeizeTag Set Seize tag
func WithSeizeTag(tag bool) OpFunc {
	return func(o *option) {
		o.seizeTag = tag
	}
}

//WithSeizeCycle Set Seize cycle
func WithSeizeCycle(t time.Duration) OpFunc {
	return func(o *option) {
		o.seizeCycle = t
	}
}

//WithExpireTime Expiration time of key
func WithExpireTime(time uint) OpFunc {
	return func(o *option) {
		o.expire = time
	}
}

//WithLockKey lock key
func WithLockKey(key string) OpFunc {
	return func(o *option) {
		o.key = key
	}
}

//WithRenewalOften Time to renew type the contract
func WithRenewalOften(f RenewalTypeFunc) OpFunc {
	return func(o *option) {
		o.RenewalOften = f
	}
}

//WithRedisTimeout Time to renew type the contract
func WithRedisTimeout(t time.Duration) OpFunc {
	return func(o *option) {
		o.redisTimeout = t
	}
}

//WithRedisClient Time to renew type the contract
func WithRedisClient(cl *redis.Client) OpFunc {
	return func(o *option) {
		o.redisClient = cl
	}
}

//DefaultRenewalTime The default is to renew the contract in half of the current time
//if is zero of 1
func DefaultRenewalTime(time uint) uint {
	t := (time + 2) / 2
	if t > 0 {
		return t
	}
	return 1
}

//CustomRenewalTime Custom renewal time
func CustomRenewalTime(time uint) uint {
	return time
}

//WithLockFailFunc get lock fail callback func
func WithLockFailFunc(lockFail LockCallbackFun) OpFunc {
	return func(o *option) {
		o.lockFailFunc = lockFail
	}
}

//WithLockSuccessFunc get lock success callback func
func WithLockSuccessFunc(lockSuss LockCallbackFun) OpFunc {
	return func(o *option) {
		o.lockSuccessFunc = lockSuss
	}
}
