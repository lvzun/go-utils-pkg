package notify

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type WXWorkerReqModel struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content             string   `json:"content"`
		MentionedList       []string `json:"mentioned_list"`
		MentionedMobileList []string `json:"mentioned_mobile_list"`
	} `json:"text"`
}

var cli = resty.New().SetTransport(&http.Transport{
	IdleConnTimeout:     60 * time.Second,
	MaxIdleConns:        10,
	MaxIdleConnsPerHost: 10,
	TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
}).SetTimeout(30 * time.Second)

func WxWorkerTextNotifyAll(tag, content, groupUrl string) error {
	return WxWorkerNotify(tag, "", content, "", groupUrl)
}
func WxWorkerNotify(tag, messageType, content, people, groupUrl string) error {
	wxWorkerReq := WXWorkerReqModel{}

	wxWorkerReq.Msgtype = "text"
	if messageType != "" {
		wxWorkerReq.Msgtype = messageType
	}

	mention := []string{"@all"}
	if people != "" {
		mention = append(mention, people)
	}

	wxWorkerReq.Text.Content = content
	wxWorkerReq.Text.MentionedList = mention
	bytes, _ := json.Marshal(wxWorkerReq)

	resp, err := cli.R().SetHeader("Content-Type", "application/json").SetBody(bytes).Post(groupUrl)
	if err != nil {
		logrus.WithField("tag", tag).Errorf("req %s,data:%s err:%s", groupUrl, string(bytes), err)
		return err
	}
	body := resp.Body()
	logrus.WithField("tag", tag).Infof("发送企业微信通知 url:%s 参数:%s,结果:%s", groupUrl, string(bytes), body)
	if len(body) == 0 {
		return fmt.Errorf("返回结果为空")
	}
	return nil
}
