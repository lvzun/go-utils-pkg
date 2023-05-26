package validateUtils

import (
	"fmt"
	"github.com/lvzun/github.com/lvzun/go-utils-pkg/cryptoUtils"
	"github.com/sirupsen/logrus"

	"regexp"
	"sort"
	"strings"
)

func IsEmpty(data string) bool {
	return VerifyEmpty(data)
}
func VerifyEmpty(data string) bool {
	if len(data) == 0 {
		return true
	}
	if data == "" {
		return true
	}
	return false
}
func CalcValidateValue(idNo string) string {
	idByteArr := []byte(idNo)
	iS := 0
	iW := [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	szVerCode := [11]string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}
	for i := 0; i < 17; i++ {
		iS += int(idByteArr[i]-48) * iW[i]
	}
	return szVerCode[iS%11]
}
func VerifyIdNo(idNo string) bool {
	if len(idNo) != 18 {
		return false
	}
	v := CalcValidateValue(idNo)
	fmt.Printf("last v:%s", v)
	if strings.HasSuffix(idNo, v) {
		return true
	} else {
		return false
	}
}

func VerifyIpv4(ipv4 string) bool {
	if len(ipv4) == 0 {
		return false
	}
	reg := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	return reg.Match([]byte(ipv4))
}

func VerifyDigit(digit string) bool {
	if len(digit) == 0 {
		return false
	}
	reg := regexp.MustCompile(`^\d+(\.\d+)?$`)
	return reg.Match([]byte(digit))
}

func VerifyPhone(phone string) bool {
	if len(phone) != 11 {
		return false
	}

	reg := regexp.MustCompile(`^\d{11}$`)
	return reg.Match([]byte(phone))
}

func VerifySign(data, key, sign string) bool {
	if len(sign) == 0 {
		logrus.Error("VerifySign 失败,clientSign为空")
		return false
	}

	serverSign := cryptoUtils.ToMd5(data + key)
	if strings.EqualFold(serverSign, strings.ToLower(sign)) {
		return true
	} else {
		logrus.Errorf("VerifySign 失败,data:%s,serverSign:%s,clientSign:%s", data+key, serverSign, sign)
		return false
	}

}

func CalcSign(param map[string][]string, key string) (string, string) {
	signData := ""
	if len(param) > 0 {
		var keys []string
		for k := range param {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			if len(param[k]) > 0 && !strings.EqualFold(k, "sign") {
				signData += param[k][0]
			}
		}
	}

	sign := cryptoUtils.ToMd5(signData + cryptoUtils.ToMd5(key))
	return sign, signData
}

func VerifySortSign(param map[string][]string, key, clientSign string) bool {
	sign, signData := CalcSign(param, key)
	if strings.EqualFold(sign, strings.ToLower(clientSign)) {
		return true
	} else {
		logrus.Errorf("VerifySortSign 失败,data:%s,clientSign:%s,serverSign:%s,", signData+cryptoUtils.ToMd5(key), clientSign, sign)
		return false
	}
}

func CheckIPList(ip, configIp string) bool {
	//内网放行
	if strings.EqualFold(ip, "127.0.0.1") || strings.EqualFold(ip, "localhost") || strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "172.16.") || strings.HasPrefix(ip, "192.168.") {
		return true
	}

	if configIp == "" || strings.Contains(configIp, ip) {
		return true
	}
	return false
}
