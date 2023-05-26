package httpUtils

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpUtil struct {
	proxy         string
	timeout       int
	PoolSize      int
	client        *http.Client
	CookieManager *SessionCookieManager
	Redirect      bool
}

func NewDefault() *HttpUtil {
	util := new(HttpUtil)
	util.PoolSize = 1
	util.timeout = 15
	util.proxy = ""
	util.CookieManager = NewCookieManager()
	util.Redirect = false
	return util
}

func New(timeout, poolSize int, proxy string) *HttpUtil {
	util := new(HttpUtil)
	util.PoolSize = poolSize
	util.timeout = timeout
	util.proxy = proxy
	util.Redirect = false
	util.CookieManager = NewCookieManager()
	return util
}

func (c *HttpUtil) SetProxy(proxy string) {
	c.proxy = proxy
}

func (c *HttpUtil) NewConnect(timeout int, proxy string) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout: time.Duration(timeout) * time.Second,
			//Deadline: time.Now().Add(time.Duration(timeout)* time.Second),
			KeepAlive: time.Duration(timeout) * time.Second,
			DualStack: true,
		}).DialContext,
		//Dial: func(netw, addr string) (net.Conn, error) {
		//	c, err := net.DialTimeout(netw, addr, time.Duration(timeout)*time.Second) //设置建立连接超时
		//	if err != nil {
		//		return nil, err
		//	}
		//	c.SetDeadline(time.Now().Add(time.Duration(timeout)* time.Second)) //设置发送接收数据超时
		//	return c, nil```
		//},
		MaxIdleConnsPerHost: 10,
		MaxIdleConns:        10,
		IdleConnTimeout:     60 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		//ExpectContinueTimeout: 1 * time.Second,
		//DisableKeepAlives:true,
	}
	if proxy != "" {
		u, err := url.Parse(proxy)
		if err != nil {
			logrus.Error("proxy 解析失败，不使用代理")
		} else {
			tr.Proxy = http.ProxyURL(u)
		}
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(timeout) * time.Second}
	if !c.Redirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return errors.New(req.URL.String())
		}
	}
	return client
}
func (c *HttpUtil) PostForm(url, queryString, form string) (string, error) {
	return c.HttpRequest("POST", url, queryString, "form")
}

func (c *HttpUtil) PostJson(url, queryString string) (string, error) {
	return c.HttpRequest("POST", url, queryString, "json")
}
func (c *HttpUtil) Post(url, queryString string) (string, error) {
	return c.HttpRequest("POST", url, queryString, "")
}
func (c *HttpUtil) Get(url string) (string, error) {
	return c.HttpRequest("GET", url, "", "")
}
func (c *HttpUtil) HttpBaseRequest(req *http.Request) ([]byte, *http.Response, error) {
	if c.client == nil {
		c.client = c.NewConnect(c.timeout, c.proxy)
	}
	res, err := c.client.Do(req)
	defer func() {
		if res != nil && res.Body != nil {
			_ = res.Body.Close()
		}
	}()
	if err != nil {
		logrus.Errorf("Error at:%v \n url:%s", err, req.URL)
		c.client = c.NewConnect(c.timeout, c.proxy)
		return nil, nil, err
	} else {
		b, err := c.responseBinaryReader(res)
		if err != nil {
			return nil, nil, err
		}
		return b, res, nil
	}
}
func (c *HttpUtil) HttpRequest(method, url, queryString, contentType string) (string, error) {

	req, err := http.NewRequest(method, url, strings.NewReader(queryString))
	if err != nil {
		logrus.Errorf("NewRequest:", err)
		return "", err
	}

	//req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")
	//req.Header.Add("accept", "*/*")
	req.Header.Add("accept-encoding", "gzip,deflate")
	if contentType == "form" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	} else if contentType == "json" {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("accept-language", "zh-CN,en-US;q=0.8")
	//req.Header.Add("Connection", "keep-alive")
	res, resp, err := c.HttpBaseRequest(req)
	_ = resp
	return string(res), err
}

func (c *HttpUtil) timeoutWarning(tag, detailed string, start time.Time, timeLimit float64) {
	dis := time.Now().Sub(start).Seconds() * 1000
	if dis > timeLimit {
		fmt.Sprintln("["+tag+"]", detailed, "using", dis, "mi")
	}
}

func (c *HttpUtil) responseBinaryReader(resp *http.Response) ([]byte, error) {
	if resp != nil && resp.Body != nil {
		//Content-Encoding
		body, err := ioutil.ReadAll(resp.Body)
		if len(body) > 2 && body[0] == 0x1f && body[1] == 0x8b {
			reader, _ := gzip.NewReader(bytes.NewReader(body))
			tmp := make([]byte, 0)
			for {
				buf := make([]byte, 1024)
				n, err := reader.Read(buf)
				if n > 0 {
					tmp = append(tmp, buf[0:n]...)
				}
				_ = buf
				if err != nil && err != io.EOF {
					logrus.Error("读取http的gzip数据异常")
					return nil, err
				}
				if err == io.EOF {
					break
				}
			}
			body = tmp
		}
		if err != nil {
			return nil, err
		}
		return body, nil
	} else {
		return nil, errors.New("response is nil")
	}
}

func (c *HttpUtil) Close() {
	if c != nil {
		_ = c.client
		_ = c.CookieManager
		_ = c.proxy
	}

}
