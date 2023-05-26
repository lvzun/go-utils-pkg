package cryptoUtils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
)

func ToMd5(str string) string {
	return fmt.Sprintf("%02x", md5.Sum([]byte(str)))
}
func HmacSHA1(data, key string) []byte {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(data))
	return h.Sum(nil)
}
