package logger

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type logrusFieldsHook struct {
	mu sync.RWMutex
}

// TODO 后期可以做成配置项
var logFieldsKeys = []string{"tag",
	"tcOrderNo",
	"duration",
	"loopTimes",
	"remoteIp",
	"secondfacility",
	"request",
	"proxy",
	"baiduOcrLogId",
	"orderId",
	"userName",
	"requestId",
	"channel",
	"lockMidFailReason",
	"httpStatus",
	"operationType",
	"tips",
	"activeBBid",
	"lockResult",
	"bbid",
	"color",
	"path",
	"did",
	"mqWorder",
	"mqUnread",
}

func (hook *logrusFieldsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *logrusFieldsHook) Fire(e *logrus.Entry) error {
	hook.mu.RLock()
	defer hook.mu.RUnlock()
	if e.Context != nil {
		for _, item := range logFieldsKeys {
			value := e.Context.Value(item)
			if value != nil {
				e.Data[item] = value
			}
		}
	}
	return nil
}
