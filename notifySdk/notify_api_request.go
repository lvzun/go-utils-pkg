package notifySdk

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/lvzun/go-utils-pkg/cryptoUtils"
	"net/http"
	"time"
)

var (
	restyClient = resty.New().SetTransport(&http.Transport{
		IdleConnTimeout:     90 * time.Second,
		MaxIdleConns:        512,
		MaxIdleConnsPerHost: 512,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}).SetTimeout(time.Duration(15) * time.Second)
)
var GlobalInstance *notifyApi

type notifyApi struct {
	reqUrl string
	appId  string
	appKey string
}

func InitGlobal(notifyApiPath, key, id string) *notifyApi {
	GlobalInstance = &notifyApi{}
	GlobalInstance.appKey = key
	GlobalInstance.appId = id
	GlobalInstance.reqUrl = notifyApiPath
	return GlobalInstance
}

func NewInstance(notifyApiPath, key, id string) *notifyApi {
	n := &notifyApi{}
	n.appKey = key
	n.appId = id
	n.reqUrl = notifyApiPath
	return n
}

func Send(sendReq *SendMessageRequestParams) error {
	if GlobalInstance != nil {
		return GlobalInstance.Send(sendReq)
	} else {
		return errors.New("notify api global instance is nil")
	}
}

func (n *notifyApi) Send(sendRequest *SendMessageRequestParams) (err error) {

	if sendRequest == nil {
		return errors.New("notify request params SendMessageRequestParams is empty")
	}

	if len(sendRequest.Text) == 0 {
		return errors.New("notify request params text is empty")
	}

	if len(sendRequest.Receiver) == 0 {
		return errors.New("notify request params receiver is empty")
	}

	marshal, err := json.Marshal(sendRequest)
	if err != nil {
		return err
	}
	sign := cryptoUtils.ToMd5(n.appId + string(marshal) + cryptoUtils.ToMd5(n.appKey))
	//logrus.Infof("signData:%s,sign=%s", n.appId+string(marshal)+cryptoUtils.ToMd5(n.appKey), sign)

	reqUrl := fmt.Sprintf("%s?appId=%s&sign=%s", n.reqUrl, n.appId, sign)

	_, err = restyClient.R().SetHeader("Content-Type", "application/json").
		SetBody(string(marshal)).Post(reqUrl)
	return err
}
