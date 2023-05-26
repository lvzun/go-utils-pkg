package cryptoUtils

import (
	"bytes"
	"compress/gzip"
	"github.com/sirupsen/logrus"
	"io"
)

func Gzip(data []byte) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(data); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return b.Bytes()
}

func UnGzip(compressSrc []byte) []byte {
	if compressSrc != nil {
		b := bytes.NewReader(compressSrc)
		var out bytes.Buffer
		r, err := gzip.NewReader(b)
		if err != nil {
			logrus.Error("UnGzip error:", err)
			return nil
		}
		io.Copy(&out, r)
		return out.Bytes()
	}
	return nil
}
