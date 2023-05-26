package netUtils

import (
	"github.com/lvzun/go-utils-pkg/validateUtils"
)

func GetRemoteAddress(headers map[string][]string, remoteIP string) string {
	remoteIpFields := []string{"X-Real-Ip", "Realip", "X-Forwarded-For", "Remote_addr"}
	for _, item := range remoteIpFields {
		value := headers[item]
		if len(value) > 0 && validateUtils.VerifyIpv4(value[0]) {
			return value[0]
		}
	}
	return remoteIP
}
