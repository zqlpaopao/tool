package pkg

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/orcaman/concurrent-map/v2"
	"github.com/redis/go-redis/v9"

	"io"
	"strings"
	"time"
)

type RdbCli struct {
	single    *redis.Client
	cluster   *redis.ClusterClient
	isCluster bool
	retry     int
	lockNum   int
	sha       *cmap.ConcurrentMap[string, string]
	groupName string
}

// Ping -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (r *RdbCli) Ping() error {
	if r.isCluster {
		return r.cluster.Ping(context.Background()).Err()
	}
	return r.single.Ping(context.Background()).Err()
}

// Lock -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (r *RdbCli) Lock(lockName string, duration time.Duration) (res int64, err error) {
	var (
		retry = r.retry
		ok    bool
		sha   string
		t     = time.Now()
	)
	if sha, ok = r.sha.Get(RDBLock); !ok {
		if sha, err = r.makeSha(RDBLock, LockCmd); err != nil {
			return
		}
	}

	for retry > 0 {
		if res, err = r.LoadScript(sha, []string{
			r.groupName,
			fmt.Sprintf("%v%v", DefaultSameSlot, lockName)},
			[]interface{}{
				t.Add(duration).Unix(),
				r.lockNum,
				t.Unix()}); res == Success {
			return
		}
		if err != nil && strings.Contains(err.Error(), "NOSCRIPT") {
			_, _ = r.makeSha(RDBLock, LockCmd)
		}
		retry--
	}

	return
}

// Renewal -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (r *RdbCli) Renewal(lockName string, duration time.Duration) (res int64, err error) {
	var (
		retry = r.retry
		ok    bool
		sha   string
		t     = time.Now()
	)
	if sha, ok = r.sha.Get(RDBRenewal); !ok {
		if sha, err = r.makeSha(RDBRenewal, RenewalCmd); err != nil {
			return
		}
	}

	fmt.Println(r.groupName, fmt.Sprintf("%v%v", DefaultSameSlot, lockName))
	for retry > 0 {
		if res, err = r.LoadScript(sha, []string{
			r.groupName,
			fmt.Sprintf("%v%v", DefaultSameSlot, lockName)},
			[]interface{}{
				t.Add(duration).Unix()}); res == Success {
			return
		}
		fmt.Println("re", res, err)
		if err != nil && strings.Contains(err.Error(), "NOSCRIPT") {
			_, _ = r.makeSha(RDBRenewal, RenewalCmd)
		}
		retry--
	}
	return

}

// GetLockInfo -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (r *RdbCli) GetLockInfo() (res map[string]string, err error) {
	if r.isCluster {
		return r.cluster.HGetAll(context.Background(), r.groupName).Result()
	} else {
		return r.single.HGetAll(context.Background(), r.groupName).Result()
	}
}

// UnLock -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (r *RdbCli) UnLock(lockName string) (res int64, err error) {
	var (
		retry = r.retry
		ok    bool
		sha   string
	)
	if sha, ok = r.sha.Get(RDBUnLock); !ok {
		if sha, err = r.makeSha(RDBUnLock, DelMemberCmd); err != nil {
			return
		}
	}

	for retry > 0 {
		if res, err = r.LoadScript(sha, []string{
			r.groupName,
			fmt.Sprintf("%v%v", DefaultSameSlot, lockName)},
			[]interface{}{}); res == Success {
			return
		}
		if err != nil && strings.Contains(err.Error(), "NOSCRIPT") {
			_, _ = r.makeSha(RDBUnLock, DelMemberCmd)
		}
		retry--
	}
	return
}

// LoadScript -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (r *RdbCli) LoadScript(sha string, keys []string, args []interface{}) (res int64, err error) {
	var (
		ctx    = context.Background()
		result interface{}
		ok     bool
	)
	if r.isCluster {
		result, err = r.cluster.EvalSha(ctx, sha, keys, args).Result()
	} else {
		result, err = r.single.EvalSha(ctx, sha, keys, args).Result()
	}
	if err != nil {
		return
	}
	if res, ok = result.(int64); !ok {
		err = errors.New("load script fail")
	}
	return
}

// NewRdbCli -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func NewRdbCli(redisOpt *RedisOption) *RdbCli {
	cmap.SHARD_COUNT = 3
	mp := cmap.New[string]()
	r := RdbCli{
		single:    nil,
		cluster:   nil,
		isCluster: redisOpt.isCluster,
		retry:     redisOpt.nodeNum + 1,
		lockNum:   redisOpt.lockNum,
		sha:       &mp,
		groupName: redisOpt.groupName,
	}
	if redisOpt.isCluster {
		r.cluster = NewCluster(redisOpt)
	} else {
		r.single = NewSingle(redisOpt)
	}
	return &r
}

// NewCmdSha -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (r *RdbCli) NewCmdSha() {
	var err error
	if _, err = r.makeSha(RDBLock, LockCmd); err != nil {
		return
	}
	if _, err = r.makeSha(RDBUnLock, DelMemberCmd); err != nil {
		return
	}
	if _, err = r.makeSha(RDBRenewal, RenewalCmd); err != nil {
		return
	}

}

// makeSha -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func (r *RdbCli) makeSha(key, cmd string) (sha string, err error) {
	var (
		cxt = context.Background()
	)
	sha = Sha(cmd)
	if r.isCluster {
		err = r.cluster.ScriptLoad(cxt, cmd).Err()
	} else {
		err = r.single.ScriptLoad(cxt, cmd).Err()
	}
	if err != nil {
		return
	}
	r.sha.Set(key, sha)
	return
}

// Sha -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func Sha(src string) string {
	h := sha1.New()
	_, _ = io.WriteString(h, src)
	return hex.EncodeToString(h.Sum(nil))
}

// NewSingle -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func NewSingle(opt *RedisOption) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         opt.addr[0],
		Password:     opt.password,
		DB:           opt.db,
		DialTimeout:  opt.readTimeout,
		ReadTimeout:  opt.readTimeout,
		WriteTimeout: opt.writeTimeout,
		PoolSize:     opt.poolSize,
	})
}

// NewCluster -- --------------------------
// --> @Describe
// --> @params
// --> @return
// -- ------------------------------------
func NewCluster(opt *RedisOption) *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        opt.addr,
		Password:     opt.password,
		DialTimeout:  opt.readTimeout,
		ReadTimeout:  opt.readTimeout,
		WriteTimeout: opt.writeTimeout,
		PoolSize:     opt.poolSize,
	})
}
