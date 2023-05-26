package redisUtils

import (
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	Rdb *redis.ClusterClient
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
