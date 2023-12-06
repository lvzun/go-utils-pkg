package cryptoUtils

import (
	"bytes"
	"fmt"
	"github.com/andybalholm/brotli"
	"io/ioutil"
)

// CompressBrotli 压缩数据使用 Brotli 算法
func CompressBrotli(data []byte) ([]byte, error) {
	var compressedBuffer bytes.Buffer

	writer := brotli.NewWriter(&compressedBuffer)
	_, err := writer.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to compress data: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close Brotli writer: %v", err)
	}

	return compressedBuffer.Bytes(), nil
}

// DecompressBrotli 解压缩使用 Brotli 算法压缩的数据
func DecompressBrotli(compressedData []byte) ([]byte, error) {
	reader := brotli.NewReader(bytes.NewReader(compressedData))
	decompressed, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %v", err)
	}

	return decompressed, nil
}
