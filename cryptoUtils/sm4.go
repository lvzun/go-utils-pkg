package cryptoUtils

import (
	"crypto/cipher"
	"github.com/ZZMarquis/gm/sm4"
)

func Sm4EcbPkcs5Encode(data []byte, key []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	data = PKCS5Padding(data, blockSize)
	encrypted := make([]byte, len(data))
	block.Encrypt(encrypted, data)
	return encrypted, nil
}

func Sm4EcbPkcs5Decode(data []byte, key []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(data))
	block.Decrypt(decrypted, data)
	decrypted = PKCS5UnPadding(decrypted)
	return decrypted, nil
}

func Sm4CbcPkcs5Encode(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	data = PKCS5Padding(data, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(data))
	blockMode.CryptBlocks(encrypted, data)
	return encrypted, nil
}

func Sm4CbcPkcs5Decode(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(data))
	blockMode.CryptBlocks(decrypted, data)
	decrypted = PKCS5UnPadding(decrypted)

	return decrypted, nil
}
