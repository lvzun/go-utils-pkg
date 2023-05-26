package cryptoUtils

import (
	"github.com/ZZMarquis/gm/sm2"
)

func Sm2Encrypt(plaintext []byte, pubKeyBytes []byte, cipherTextType int) ([]byte, error) {
	pub, err := sm2.RawBytesToPublicKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	cipherText, err := sm2.Encrypt(pub, plaintext, sm2.Sm2CipherTextType(cipherTextType))
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}
