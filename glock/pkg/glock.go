package pkg

import (
	"context"
	"fmt"
	redisScript "github.com/zqlpaopao/tool/redis/pkg"
	"github.com/zqlpaopao/tool/retry/pkg"
	"runtime/debug"
	"sync"
	"time"
)

type gLock struct {
	isMaster bool
	opt      *option
	err      error
	reTry    *pkg.RetryManager
	lock     sync.RWMutex
}

//NewGlock get lock object
func NewGlock(f ...Option) Glock {
	return &gLock{
		opt:      NewOptions(f...),
		isMaster: false,
		reTry:    pkg.NewRetryManager(pkg.WithRetryInterval(time.Millisecond * 500)),
		lock:     sync.RWMutex{},
	}
}

//Lock Start locking. If you turn on the short sign, you will always try to lock
func (g *gLock) Lock(arg ...interface{}) Glock {
	g.checkRedis()
	if g.err != nil{return g}
	g.setNx(arg...)
	go g.renewalOften()
	g.joinMemberGroup()
	go g.seizeLock(arg...)
	return g
}

//Trying to get master role
func (g *gLock) setNx(arg ...interface{}) {
	var (
		isMaster bool
		err      error
	)
	ctx, cancel := context.WithTimeout(context.Background(), DefaultRedisTimeOut)
	defer cancel()
	if isMaster, err = g.opt.redisClient.SetNX(ctx, Lock, Master, time.Duration(g.opt.expire)*time.Second).Result(); err != nil {
		return
	}
	g.setMaster(isMaster)
	g.callbackFunc(arg...)
}

//check redisClient
func (g *gLock)checkRedis(){
	_, g.err = g.opt.redisClient.Ping(context.TODO()).Result()
}

//Join the competition group
func (g *gLock) joinMemberGroup() {
	g.reTry.DoSync(func() bool {
		role := Slave
		if g.IsMaster() {
			role = Master
		}
		if _, err := g.opt.redisClient.Eval(context.Background(), redisScript.HSetANdExpire, []string{memberGroup}, 1, int(g.opt.expire), g.opt.key, role).Result(); !redisScript.IsRedisNilError(err) {
			return false
		}
		return true
	})
}

//callbackFunc Successful and failed callback functions
func (g *gLock) callbackFunc(arg ...interface{}) {
	if g.IsMaster() && g.opt.lockSuccessFunc != nil {
		g.opt.lockSuccessFunc(arg...)
	}
	if !g.IsMaster() && g.opt.lockFailFunc != nil {
		g.opt.lockFailFunc(arg...)
	}
}

//renewalOften Obtain the lock and cycle the renewal operation
func (g *gLock) renewalOften() {
	savePanic()()
	if !g.IsMaster() {
		return
	}
	renewal := g.opt.RenewalOften(g.opt.expire)
	timer := time.NewTimer(time.Second * time.Duration(renewal-1))
	script := g.makScript()
	for {
		select {
		case <-timer.C:
			g.doRenewal(script, renewal)
			g.joinMemberGroup()
			timer.Reset(time.Second * time.Duration(renewal-1))
		case <-g.opt.renewalTag:
			goto END
		default:
			if !g.IsMaster() {
				goto END
			}
			time.Sleep(time.Second)
		}
	}
END:
	timer.Stop()
}

//makScript Preloaded Lua script
func (g *gLock) makScript() string {
	var str string
	var err error
	if str, err = g.opt.redisClient.ScriptLoad(context.TODO(), redisScript.SetExpireByTTl).Result(); nil != err {
		return ""
	}
	return str
}

//doRenewal Execute the renewal of lua and check the effective time of the master
func (g *gLock) doRenewal(script string, renewal uint) {
	if script != "" {
		g.reTry.DoSync(func() bool {
			if _, err := g.opt.redisClient.EvalSha(context.Background(), script, []string{Lock}, Master, int(renewal)).Result(); !redisScript.IsRedisNilError(err) {
				return false
			}
			return true
		})
		return
	}
	g.reTry.DoSync(func() bool {
		if _, err := g.opt.redisClient.Eval(context.Background(), script, []string{Lock}, Master, int(renewal)).Result(); !redisScript.IsRedisNilError(err) {
			return false
		}
		return true
	})
}

//seizeLock If you don't get the lock and want to get the lock, keep trying to get the lock
func (g *gLock) seizeLock(arg ...interface{}) {
	savePanic()()
	if !g.opt.seizeTag {
		return
	}
	timer := time.NewTimer(g.opt.seizeCycle)
	for {
		select {
		case <-timer.C:
			timer.Reset(g.opt.seizeCycle)
			if g.IsMaster() {
				continue
			}
			g.setNx(arg...)
			go g.renewalOften()
			g.joinMemberGroup()
		case <-g.opt.seizeClose:
			goto END
		default:
			time.Sleep(time.Second)
		}
	}
END:
	timer.Stop()

}

//UnLock free lock
func (g *gLock) UnLock() Glock {
	g.opt.seizeClose <- struct{}{}
	if g.IsMaster() {
		g.opt.renewalTag <- struct{}{}
	}
	//close(g.opt.seizeClose)
	//close(g.opt.renewalTag)
	g.reTry.DoSync(func() bool {
		if !g.IsMaster() {
			return true
		}
		if _, err := g.opt.redisClient.Del(context.Background(), Lock).Result(); err != nil {
			return false
		}
		if _, err := g.opt.redisClient.HDel(context.Background(), memberGroup,g.opt.key).Result(); err != nil {
			return false
		}
		g.setMaster(false)
		return true
	})
	return g
}

//IsMaster Is it a master
func (g *gLock) IsMaster() (b bool) {
	g.lock.Lock()
	b = g.isMaster
	g.lock.Unlock()
	return
}

//set master
func (g *gLock) setMaster(b bool) {
	g.lock.Lock()
	g.isMaster = b
	g.lock.Unlock()
}

//Error get error
func (g *gLock) Error() error {
	return g.err
}

//GetMembers get all members
func (g *gLock) GetMembers() (mem map[string]string, err error) {
	return g.opt.redisClient.HGetAll(context.Background(), memberGroup).Result()
}

//tidy panic
func savePanic() func() {
	return func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			fmt.Println(string(debug.Stack()))
		}
	}
}
