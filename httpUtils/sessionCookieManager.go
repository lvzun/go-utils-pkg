package httpUtils

import (
	"fmt"
	"github.com/lvzun/go-utils-pkg/utils"
	"net/http"
	"strings"
)

type SessionCookieManager struct {
	Key   []string // cookie 的key
	Value []string // cookie 的vlaue
}

func NewCookieManager() *SessionCookieManager {
	return &SessionCookieManager{
		Key:   make([]string, 0),
		Value: make([]string, 0),
	}
}

func (cookieManager *SessionCookieManager) ParseCookieString(cookieString string) {
	if len(strings.TrimSpace(cookieString)) > 0 {
		cookieArray := strings.Split(cookieString, ";")
		if len(cookieArray) > 0 {
			for _, cookie := range cookieArray {
				if strings.Index(cookie, "=") > 0 {
					cookies := strings.Split(cookie, "=")
					cookieManager.Key = append(cookieManager.Key, strings.TrimSpace(cookies[0]))
					cookieManager.Value = append(cookieManager.Value, strings.TrimSpace(cookies[1]))
				}
			}
		}
	}
}

// 有则更新 无则插入
func (scm *SessionCookieManager) Upsert(key, value string) {
	//需要判断是否这个key已经存在了 如果存在就更新
	//slice 是否包含某个字符串 -1 表示没有
	indexOf := utils.IndexOfSlice(key, scm.Key)
	if indexOf == -1 {
		//插入
		scm.Key = append(scm.Key, key)
		scm.Value = append(scm.Value, value)
	} else {
		//这里需要更新
		scm.Value[indexOf] = value
	}
	return
}

// cookie 键值 slice to string
func (scm *SessionCookieManager) String() string {
	cookieString := ""
	keyLength := len(scm.Key)
	valueLength := len(scm.Value)
	if keyLength == valueLength {
		for index := 0; index < keyLength; index++ {
			k := scm.Key[index]
			v := scm.Value[index]
			cookieString = cookieString + k + "=" + v + "; "
		}
	} else {
		fmt.Println("cookie的键值 数量不匹配")
	}
	return cookieString
}

// 根据响应里的cookie 更新cookie 从而达到维持session的目的

func (scm *SessionCookieManager) UpdateFromResponseWithCallback(cookies []*http.Cookie, callback func(key, value string)) {
	for _, c := range cookies {
		scm.Upsert(c.Name, c.Value)
		if callback != nil {
			callback(c.Name, c.Value)
		}
	}
}

// 根据响应里的cookie 更新cookie 从而达到维持session的目的
func (scm *SessionCookieManager) UpdateWithCallback(res *http.Response, callback func(key, value string)) {
	if len(res.Header["Set-Cookie"]) > 0 {
		if len(res.Cookies()) == 0 {
			return
		} else {
			for _, c := range res.Cookies() {
				scm.Upsert(c.Name, c.Value)
				if callback != nil {
					callback(c.Name, c.Value)
				}
			}
		}
	}
}
func (scm *SessionCookieManager) Update(res *http.Response) {
	scm.UpdateWithCallback(res, nil)
}
func (scm *SessionCookieManager) ClearCookies() {
	scm.Key = make([]string, 0)
	scm.Value = make([]string, 0)
}
