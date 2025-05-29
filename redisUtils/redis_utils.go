package redisUtils

import (
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	Rdb       *redis.ClusterClient
	RdbSingle *redis.Client
)

func InitRedisCluster(addrList []string, auth string) {
	Rdb = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrList,
		Password:     auth,
		DialTimeout:  time.Second * 10,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	})
}
func InitRedis(addr string, auth string, db int) {
	RdbSingle = redis.NewClient(&redis.Options{
		Addr:     addr, // Redis 服务器地址
		Password: auth, // 密码，没有则为空
		DB:       db,   // 使用默认 DB
	})
}
