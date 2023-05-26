package httpsqs

import (
	"encoding/json"
	"github.com/lvzun/go-utils-pkg/httpUtils"
	"github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

type QueueInfo struct {
	Name     string `json:"name"`
	Maxqueue int    `json:"maxqueue"`
	Putpos   int    `json:"putpos"`
	Putlap   int    `json:"putlap"`
	Getpos   int    `json:"getpos"`
	Getlap   int    `json:"getlap"`
	Unread   int    `json:"unread"`
}

type Queue struct {
	Url      string
	HttpUtil *httpUtils.HttpUtil
}

func New(url string) *Queue {
	q := new(Queue)
	q.Url = url
	q.HttpUtil = httpUtils.New(10, 1, "")
	return q
}

func (q *Queue) InQueue(queue_name, data string) bool {
	queueUrl := q.Url
	queueUrl = strings.Replace(queueUrl, "{0}", queue_name, -1)
	queueUrl = strings.Replace(queueUrl, "{1}", "put", -1)
	queueUrl += "&data=" + url.QueryEscape(data)
	for i := 0; i < 3; i++ {
		res, err := q.HttpUtil.Get(queueUrl)
		if err != nil {
			logrus.Errorf("与队列出错，err：%v", err)
			continue
		}
		if string(res) != "HTTPSQS_PUT_OK" {
			logrus.Error("写队列出错，结果：" + string(res) + " 数据：" + queueUrl)
			continue
		} else {
			return true
		}
	}
	return false

}
func (q *Queue) OutQueue(queue_name string) string {
	queueUrl := q.Url
	queueUrl = strings.Replace(queueUrl, "{0}", queue_name, -1)
	queueUrl = strings.Replace(queueUrl, "{1}", "get", -1)
	res, _ := q.HttpUtil.Get(queueUrl)

	if len(res) == 0 || strings.EqualFold(string(res), "HTTPSQS_GET_END") {
		return ""
	} else {
		data, _ := url.QueryUnescape(string(res))
		return data
	}
}

func (q *Queue) QueueInfo(queue_name string) *QueueInfo {
	resp := new(QueueInfo)
	queueUrl := q.Url
	queueUrl = strings.Replace(queueUrl, "{0}", queue_name, -1)
	queueUrl = strings.Replace(queueUrl, "{1}", "status_json", -1)
	for i := 0; i < 3; i++ {
		res, err := q.HttpUtil.Get(queueUrl)
		if len(res) == 0 || err != nil {
			logrus.Errorf("第%d次，查询队列信息出错，数据：%s,err:%v", i, queueUrl, err)
			continue
		}
		err = json.Unmarshal([]byte(res), &resp)
		if err != nil {
			logrus.Errorf("第%d次，查询数据出错，结果：%s 数据：%s，err：%v", i, string(res), queueUrl, err)
			continue
		}
		return resp
	}
	return resp
}
