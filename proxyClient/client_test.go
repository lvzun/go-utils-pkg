package proxyClient

import (
	"fmt"
	"strings"
	"testing"
)

func GetHostFromProxy1(proxy string) string {
	if len(proxy) == 0 {
		return ""
	}
	index := strings.LastIndex(proxy, "@")

	if index >= 0 {
		proxy = proxy[index+1:]
	}
	index = strings.LastIndex(proxy, ":")

	if index >= 0 {
		proxy = proxy[:index]
	}
	return proxy
}
func TestGetHostFromProxy(t *testing.T) {
	//fmt.Sprintf("test1:%s\n",GetHostFromProxy1("ssds:sdsd@aldkasd@192.168.0.1:2323/"))
	//fmt.Sprintf("test1:%s\n",GetHostFromProxy1("192.168.0.1:2323/"))
	//fmt.Sprintf("test1:%s\n",GetHostFromProxy1("192.168.0.1"))

	client := NewProxyClient()
	client.releaseProxyList["a"] = ProxyInfo{Proxy: "a"}
	client.releaseProxyList["b"] = ProxyInfo{Proxy: "b"}
	client.releaseProxyList["c"] = ProxyInfo{Proxy: "c"}
	tmp := client.releaseProxyList
	for _, p := range tmp {
		delete(client.releaseProxyList, p.Proxy)
		delete(client.releaseProxyList, p.Proxy)

		fmt.Printf("map len:%d", len(client.releaseProxyList))
	}
	fmt.Printf("map len:%d", len(client.releaseProxyList))

}
