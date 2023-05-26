package cache

import (
	"github.com/coocood/freecache"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"strconv"
	"strings"
)

var (
	FreeCache *freecache.Cache
)

func init() {
	cacheSize := 100 * 1024 * 1024
	FreeCache = freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)
}

func GetBool(key string) bool {
	key = strings.ToLower(key)
	value, err := FreeCache.Get([]byte(key))
	if err != nil {
		logrus.Errorf("FreeCache GetBool %s err:%v", key, err)
		return false
	}
	parseBool, err := strconv.ParseBool(string(value))
	if err != nil {
		logrus.Errorf("FreeCache GetBool Convert %s,err:%v", string(value), err)
		return false
	}
	return parseBool
}

func GetInt32(key string, defaultValue int) int {
	key = strings.ToLower(key)
	value, err := FreeCache.Get([]byte(key))
	if err != nil {
		logrus.Errorf("FreeCache GetInt32 %s err:%v", key, err)
		return defaultValue
	}
	parseInt, err := strconv.Atoi(string(value))
	if err != nil {
		logrus.Errorf("FreeCache GetBool Convert %s,err:%v", string(value), err)
		return defaultValue
	}
	return parseInt
}

func GetString(key string, defaultValue string) string {
	key = strings.ToLower(key)
	value, err := FreeCache.Get([]byte(key))
	if err != nil {
		logrus.Errorf("FreeCache GetString %s err:%v", key, err)
		return defaultValue
	}
	return string(value)
}
