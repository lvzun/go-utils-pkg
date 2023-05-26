package proxyClient

import (
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/lvzun/go-utils-pkg/netUtils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ProxyInfo struct {
	Proxy     string
	Supplier  string
	StartTime time.Time
	Count     int
}

var ProxyDispatcherUrl string
var ProxyAuth string

type ProxyClient struct {
	releaseProxyList map[string]ProxyInfo
}

var restyHttp = resty.New().SetTransport(&http.Transport{
	IdleConnTimeout:     90 * time.Second,
	MaxIdleConns:        512,
	MaxIdleConnsPerHost: 512,
	TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
}).SetTimeout(time.Duration(15) * time.Second)

func NewProxyClient() *ProxyClient {
	client := &ProxyClient{}
	client.releaseProxyList = make(map[string]ProxyInfo)
	return client
}

func (client *ProxyClient) GetProxy(id string) ProxyInfo {
	proxy := ""
	supplier := ""

	for i := 0; i < 10; i++ {
		post, err := restyHttp.R().SetBody("").Post("http://" + ProxyDispatcherUrl + "/GetusableIP")
		if err != nil {
			logrus.Infof("第%d次,拉取代理失败,%s", i+1, err.Error())
		}

		tmp := string(post.Body())

		if len(tmp) > 0 {
			tmps := strings.Split(tmp, "|")
			if len(tmps) > 1 {
				proxy = tmps[0]
				if len(ProxyAuth) > 0 {
					proxy = strings.Replace(proxy, "http://", "http://"+ProxyAuth+"@", -1)
				}
				supplier = tmps[1]
				if client.ProxyVerifyUri(id, proxy) {
					break
				}
				time.Sleep(100 * time.Microsecond)
			}
		}
	}
	p := ProxyInfo{proxy, supplier, time.Now(), 0}
	client.releaseProxyList[p.Proxy] = p

	host, port := GetHostAndPortFromProxy(proxy)
	logrus.Infof(fmt.Sprintf("使用代理:%s:%d", host, port))
	return p
}

func (client *ProxyClient) ReleaseAll() {
	if len(client.releaseProxyList) > 0 {
		temp := client.releaseProxyList
		for _, p := range temp {
			client.ReleaseProxy(p)
			client.OfflineProxy(p, "1")
			delete(client.releaseProxyList, p.Proxy)
		}
	}
}

func (client *ProxyClient) ReleaseProxy(proxy ProxyInfo) {
	if len(ProxyAuth) > 0 {
		proxy.Proxy = strings.Replace(proxy.Proxy, "http://"+ProxyAuth+"@", "http://", -1)
	}

	for i := 0; i < 3; i++ {
		post, err := restyHttp.R().SetBody("data=" + proxy.Proxy).Post("http://" + ProxyDispatcherUrl + "/UsedRetrunIP")
		if err != nil {
			logrus.Infof("第%d次,释放代理失败,%s", i+1, err.Error())
			continue
		}
		logrus.Infof("释放代理结果：%s", string(post.Body()))

		break
	}
	delete(client.releaseProxyList, proxy.Proxy)
}

func (client *ProxyClient) ChangeProxy(proxy ProxyInfo) {
	if len(ProxyAuth) > 0 {
		proxy.Proxy = strings.Replace(proxy.Proxy, "http://"+ProxyAuth+"@", "http://", -1)
	}

	for i := 0; i < 3; i++ {
		post, err := restyHttp.R().SetBody("data=" + proxy.Proxy).Post("http://" + ProxyDispatcherUrl + "/NotifyUnused")
		if err != nil {
			logrus.Infof("第%d次,反馈代理异常,%s", i+1, err.Error())
			continue
		}
		logrus.Infof("反馈代理结果：%s", string(post.Body()))
		break
	}
}

func (client *ProxyClient) ChangeProxyByUrl(proxy string) {
	for i := 0; i < 3; i++ {
		post, err := restyHttp.R().SetBody("data=" + proxy).Post("http://" + ProxyDispatcherUrl + "/NotifyUnused")
		if err != nil {
			logrus.Infof("第%d次,反馈代理异常,%s", i+1, err.Error())
			continue
		}
		logrus.Infof("反馈代理结果：%s", string(post.Body()))
		break
	}
}

func (client *ProxyClient) ProxyVerifyUri(id, uri string) bool {
	ipAddress, port := GetHostAndPortFromProxy(uri)
	flag := CheckNetwork(ipAddress, port)
	if flag {
		return true
	} else {
		logrus.Infof("代理[%s:%d]连通检测不通过,开始更换代理,", ipAddress, port)
		client.ChangeProxyByUrl(uri)
		return false
	}
}

// 代理验证是否可用
func CheckNetwork(ipAddress string, port int) bool {
	_, err := netUtils.Ping(ipAddress)
	if err != nil {
		return false
	}
	return netUtils.TelNet(ipAddress, fmt.Sprintf("%d", port))
}

func GetHostAndPortFromProxy(paramProxy string) (string, int) {
	if len(paramProxy) == 0 {
		return "", 0
	}
	var err error
	host := ""
	port := 0
	proxy := paramProxy
	index := strings.LastIndex(proxy, "@")

	if index >= 0 {
		proxy = proxy[index+1:]
	}
	index = strings.LastIndex(proxy, ":")

	if index >= 0 {
		host = proxy[:index]
		port, err = strconv.Atoi(proxy[index+1:])
		if err != nil {
			logrus.Errorf("解析代理端口出错proxy string：%s", paramProxy)
			return "", 0
		}
		return host, port
	} else {
		if strings.HasPrefix(paramProxy, "https://") {
			return proxy, 443
		} else if strings.HasPrefix(paramProxy, "http://") {
			return proxy, 80
		} else {
			return proxy, 0
		}
	}
}

func ChangeIP() bool { //重试三次拨号
	for i := 0; i < 3; i++ {
		post, err := restyHttp.R().SetBody("").Post("http://localhost:20185/ChangeIp")
		if err != nil {
			logrus.Infof("第%d次,ChangeIp异常,%s", i+1, err.Error())
			continue
		}
		result := string(post.Body())
		twip := strings.Split(result, "|")
		if twip[0] == "success" {
			return true
		}
		logrus.Infof("ChangeIp结果：%s", result)
		break
	}
	return false
}

// offValue代理强制下线次数
func (client *ProxyClient) OfflineProxyByUrl(proxy string, offValue string) {
	if len(ProxyAuth) > 0 {
		proxy = strings.Replace(proxy, "http://"+ProxyAuth+"@", "http://", -1)
	}

	for i := 0; i < 3; i++ {
		post, err := restyHttp.R().SetBody("data=" + proxy + "&off_code=" + offValue).Post("http://" + ProxyDispatcherUrl + "/NotifyOffline")
		if err != nil {
			logrus.Infof("第%d次,反馈代理异常,%s", i+1, err.Error())
			continue
		}
		result := string(post.Body())
		logrus.Infof("反馈代理结果：%s", result)
		break
	}
}

// offValue代理强制下线次数
func (client *ProxyClient) OfflineProxy(proxy ProxyInfo, offValue string) {
	if len(ProxyAuth) > 0 {
		proxy.Proxy = strings.Replace(proxy.Proxy, "http://"+ProxyAuth+"@", "http://", -1)
	}

	for i := 0; i < 3; i++ {
		post, err := restyHttp.R().SetBody("data=" + proxy.Proxy + "&off_code=" + offValue).Post("http://" + ProxyDispatcherUrl + "/NotifyOffline")
		if err != nil {
			logrus.Infof("第%d次,反馈代理异常,%s", i+1, err.Error())
			continue
		}
		result := string(post.Body())
		logrus.Infof("反馈代理结果：%s", result)
		break
	}
}

func getPorxyIp(porxyString string) string {
	var proxyString string
	tmpProxyStrings := strings.Split(porxyString, "@")

	if len(tmpProxyStrings) > 1 {
		proxyString = tmpProxyStrings[1]
	} else {
		proxyString = tmpProxyStrings[0]
	}
	return proxyString
}
