package pkg

import "github.com/go-redis/redis/v8"



type redisManager struct {
	opt *redis.Options
	redisClient *redis.Client
}

//NewRedisManager make redis object
func NewRedisManager(opt *redis.Options)*redisManager{
	return &redisManager{
		opt: opt,
		redisClient: &redis.Client{},
	}
}

//GetClient 获取redis连接
func(r *redisManager)GetClient(){
	r.redisClient = redis.NewClient(r.opt)
}

//ReleaseAll 释放链接
func (r *redisManager) ReleaseAll() {
	_ = r.redisClient.Close()
}

func MakeRedisZ(Score float64, Member interface{}) redis.Z {
	return redis.Z{Score: Score, Member: Member}
}

//MakeEmptyRedisZSlice 生成空有序集合Z切片
func MakeEmptyRedisZSlice() (redisZSlice []redis.Z) {
	return
}

//IsRedisNilError 是否为空错误 redis错误
func IsRedisNilError(err error) bool {
	return err == redis.Nil
}

// MakeRedisZRangeBy 生成有序集合ZRangeBy信息
func MakeRedisZRangeBy(min, max string, offset, count int64) redis.ZRangeBy {
	return redis.ZRangeBy{min, max, offset, count}
}

//MakeNewScript 生成新的Redis脚本
func MakeNewScript(src string) *redis.Script {
	return redis.NewScript(src)
}

