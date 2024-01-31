package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func ToMd5(str string) string {
	return fmt.Sprintf("%02x", md5.Sum([]byte(str)))
}
func HmacSHA1(data, key string) []byte {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(data))
	return h.Sum(nil)
}
func RandomInt(end int) int {
	return rand.New(rand.NewSource(int64(time.Now().Nanosecond()))).Intn(end)
}

func CheckWorkTime(work_time string, t time.Time) bool {
	if work_time == "" {
		return true
	}
	index := strings.Index(work_time, "-")

	if checkWorkStartTime(work_time[0:index], t) && checkWorkEndTime(work_time[index+1:], t) {
		return true
	} else {
		return false
	}
}

func checkWorkEndTime(work_end string, now time.Time) bool {
	if work_end == "" {
		return true
	}

	index := strings.Index(work_end, ":")
	confHour, err := strconv.Atoi(work_end[0:index])
	if err != nil {
		logrus.Errorf("confHour:%d", confHour)
		return true
	}
	confMinute, err := strconv.Atoi(work_end[index+1:])
	if err != nil {
		logrus.Errorf("confMinute:%d", confHour)
		return true
	}

	if now.Hour() > confHour {
		//logrus.Error("false 当前时间%d 配置时间%d", now.Hour(), confHour)
		return false
	} else if now.Hour() == confHour {
		if now.Minute() > confMinute {
			//logrus.Error("false 当前时间%d:%d 配置时间%d:%d", now.Hour(), now.Minute(), confHour, confMinute)
			return false
		}
	}
	return true
}

func checkWorkStartTime(work_start string, now time.Time) bool {
	if work_start == "" {
		return true
	}

	index := strings.Index(work_start, ":")
	confHour, err := strconv.Atoi(work_start[0:index])
	if err != nil {
		logrus.Errorf("confHour:%d", confHour)
		return true
	}
	confMinute, err := strconv.Atoi(work_start[index+1:])
	if err != nil {
		logrus.Errorf("confMinute:%d", confHour)
		return true
	}

	if now.Hour() < confHour {
		//logrus.Error("false 当前时间%d 配置时间%d", now.Hour(), confHour)
		return false
	} else if now.Hour() == confHour {
		if now.Minute() < confMinute {
			//logrus.Error("false 当前时间%d:%d 配置时间%d:%d", now.Hour(), now.Minute(), confHour, confMinute)
			return false
		}
	}
	return true
}

// 返回 在slice 中的索引值，没有找到的话返回-1
func IndexOfSlice(item string, slice []string) int {
	for index, value := range slice {
		if item == value {
			return index
		}
	}
	return -1
}

func CheckErr(err error, message string) {
	if err != nil {
		logrus.Error(message+":", err)
	}
}

func GetLocalIp() string {
	addrSlice, err := net.InterfaceAddrs()
	if nil != err {
		logrus.Error("Get local IP addr failed!!!")
		return "localhost"
	}
	for _, addr := range addrSlice {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if nil != ipnet.IP.To4() {
				return ipnet.IP.String()
			}
		}
	}
	return "localhost"
}
func LocalMac() string {
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Error : " + err.Error())
	}
	for _, inter := range interfaces {
		mac := inter.HardwareAddr //获取本机MAC地址
		return mac.String()
	}
	return ""
}

func CpuNum() int {
	return runtime.NumCPU()
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
func StrFirstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122 {
		strArry[0] -= 32
	}
	return string(strArry)
}
func StrFirstToLower(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 65 && strArry[0] <= 90 {
		strArry[0] += 32
	}
	return string(strArry)
}
func Timestamp() int64 {
	return time.Now().UnixNano() / 1e6
}

func OutOfInfo(data string, startCount, endCount int) string {

	if len(data) < startCount || len(data) < endCount {
		return data
	} else {
		return data[:startCount] + "***" + data[len(data)-endCount:]
	}
}

func IsJson(data string) bool {
	if len(data) == 0 {
		return false
	}
	data = strings.TrimSpace(data)
	if strings.HasPrefix(data, "{") || strings.HasPrefix(data, "[") {
		return true
	} else {
		return false
	}
}

func GenerateRandomString(length int) string {
	// 定义包含所有可能字符的字符串
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 生成随机字符串
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}

	return string(result)
}
