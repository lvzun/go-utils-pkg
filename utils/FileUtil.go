package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Fatal(err)
	}
	return dir
}
func GetLogPath(curPath string) string {
	path := GetCurrentDirectory() + curPath
	if !FileExists(path) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			logrus.Fatalf("目录%s创建失败", path)
		}
	}
	if strings.EqualFold("windows", runtime.GOOS) {
		path = strings.Replace(path, "\\", "\\\\", -1)
		path = strings.Replace(path, "/", "\\\\", -1)
	} else {
		path = strings.Replace(path, "\\", "/", -1)
	}
	return path
}
func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func EnsureDir(dirName string) error {
	// 使用Stat获取文件或目录的状态
	_, err := os.Stat(dirName)

	// 如果返回的错误是因为文件或目录不存在
	if os.IsNotExist(err) {
		// 使用MkdirAll创建目录，包括任何必需的父目录
		return os.MkdirAll(dirName, 0755) // 0755是常用的目录权限
	}
	return err
}
