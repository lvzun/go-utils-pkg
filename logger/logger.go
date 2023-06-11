package logger

import (
	"fmt"
	"github.com/lvzun/go-utils-pkg/netUtils"
	"os"
	"path"
	"path/filepath"

	"github.com/gemnasium/logrus-graylog-hook/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitGrayLog 初始化Graylog
func InitGrayLog(addr, module string) {
	logrus.SetReportCaller(false)

	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.TraceLevel)
	hook := graylog.NewGraylogHook(addr, map[string]interface{}{"module": module})
	hook.Level = logrus.TraceLevel
	logrus.AddHook(hook)
}

// InitFileLog 初始化本地文件日志
func InitFileLog(filename string) {
	logrus.SetReportCaller(false)

	currDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	logDir := path.Join(currDir, "logs")
	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		_ = os.Mkdir(logDir, os.ModePerm)
	}
	logfilepath := path.Join(logDir, filename)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetOutput(&lumberjack.Logger{
		Filename:   logfilepath,
		MaxSize:    10, // Max megabytes before log is rotated
		MaxBackups: 10, // Max number of old log files to keep
		MaxAge:     15, // Max number of days to retain log files
		Compress:   false,
	})
	logrus.SetLevel(logrus.TraceLevel)
}

func InitLogWithFieldCustomHook(grayArr string, fields map[string]interface{}, console bool, customHook logrus.Hook) {
	// graylog
	logrus.SetReportCaller(false)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.TraceLevel)

	if customHook != nil {
		logrus.AddHook(customHook)
	}

	if len(grayArr) > 2 {
		hook := graylog.NewGraylogHook(grayArr, fields)
		hook.Level = logrus.DebugLevel
		ip, _ := netUtils.CurrentIP()
		if ip != "" {
			hook.Host = fmt.Sprintf("%s/%d", ip, os.Getpid())
		}
		logrus.AddHook(hook)
	}

	// local
	currDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	logDir := path.Join(currDir, "logs")
	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		os.Mkdir(logDir, os.ModePerm)
	}
	module := fields["facility"]
	filename := fmt.Sprintf("%s.log", module.(string))
	logfilepath := path.Join(logDir, filename)
	if console {
		logrus.SetOutput(os.Stdout)
	} else {
		logrus.SetOutput(&lumberjack.Logger{
			Filename:   logfilepath,
			MaxSize:    10, // Max megabytes before log is rotated
			MaxBackups: 30, // Max number of old log files to keep
			MaxAge:     30, // Max number of days to retain log files
			Compress:   false,
		})
	}
}
func InitLogWithField(grayArr string, fields map[string]interface{}, console bool) {
	InitLogWithFieldCustomHook(grayArr, fields, console, nil)
}

// InitLog 初始化本地及服务端日志
func InitLog(addr, module string, console bool) {
	m := map[string]interface{}{"facility": module}
	InitLogWithField(addr, m, console)
}
