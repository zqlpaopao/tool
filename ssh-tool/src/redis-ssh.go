package src

import (
	"context"
	"github.com/go-redis/redis/v8"
	redisGO "github.com/gomodule/redigo/redis"
	"net"
)

type RedisConfig struct {
	UserName, PassWd, IpPort string
	DbNum                    int
}

//GoRedisClient -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
//Redis(go-redis/redis)
func (r *RedisConfig) GoRedisClient(client Client) (redisCli *redis.Client, err error) {
	redisCli = redis.NewClient(&redis.Options{
		Network: "tcp", // 连接方式，默认使用tcp，可省略
		Addr:    r.IpPort,
		DB:      r.DbNum, // 选择要操作的数据库，默认是0 （redis中select index命令）
		Dialer: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
			return client.client.Dial(network, addr)
		},
	})

	if err = redisCli.Ping(context.TODO()).Err(); nil != err {
		return nil, err
	}
	return
}

//RedisGOClient -- ----------------------------
//--> @Description
//--> @Param
//--> @return
//-- ----------------------------
//https://github.com/gomodule/redigo
func (r *RedisConfig) RedisGOClient(client Client) (redisConn redisGO.Conn, err error) {
	var conn net.Conn
	if conn, err = client.client.Dial("tcp", r.IpPort); nil != err {
		return
	}
	redisConn = redisGO.NewConn(conn, -1, -1)
	return
}
