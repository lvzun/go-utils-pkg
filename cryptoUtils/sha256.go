package cryptoUtils

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Sha256(data []byte) string {
	return fmt.Sprintf("%02x", sha256.Sum256(data))
}

func HmacSHA256(data, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1(data []byte) [20]byte {
	return sha1.Sum(data)
}
