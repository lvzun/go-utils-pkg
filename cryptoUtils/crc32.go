package cryptoUtils

import "hash/crc32"

func Crc32(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}
