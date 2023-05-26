package cryptoUtils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/sirupsen/logrus"

	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
)

func GenRsaKey(filePath string) error {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create(filePath + "private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create(filePath + "public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

var (
	ErrInputSize  = errors.New("input size too large")
	ErrEncryption = errors.New("encryption error")
)

func RsaEncWithPriKey(priv *rsa.PrivateKey, data []byte) (enc []byte, err error) {

	k := (priv.N.BitLen() + 7) / 8
	tLen := len(data)
	// rfc2313, section 8:
	// The length of the data D shall not be more than k-11 octets
	//if tLen > k-11 {
	//	err = errors.New("input size too large")
	//	return
	//}
	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < k-tLen-1; i++ {
		em[i] = 0xff
	}
	copy(em[k-tLen:k], data)
	c := new(big.Int).SetBytes(em)
	if c.Cmp(priv.N) > 0 {
		err = ErrEncryption
		return
	}
	var m *big.Int
	var ir *big.Int
	if priv.Precomputed.Dp == nil {
		m = new(big.Int).Exp(c, priv.D, priv.N)
	} else {
		// We have the precalculated values needed for the CRT.
		m = new(big.Int).Exp(c, priv.Precomputed.Dp, priv.Primes[0])
		m2 := new(big.Int).Exp(c, priv.Precomputed.Dq, priv.Primes[1])
		m.Sub(m, m2)
		if m.Sign() < 0 {
			m.Add(m, priv.Primes[0])
		}
		m.Mul(m, priv.Precomputed.Qinv)
		m.Mod(m, priv.Primes[0])
		m.Mul(m, priv.Primes[1])
		m.Add(m, m2)

		for i, values := range priv.Precomputed.CRTValues {
			prime := priv.Primes[2+i]
			m2.Exp(c, values.Exp, prime)
			m2.Sub(m2, m)
			m2.Mul(m2, values.Coeff)
			m2.Mod(m2, prime)
			if m2.Sign() < 0 {
				m2.Add(m2, prime)
			}
			m2.Mul(m2, values.R)
			m.Add(m, m2)
		}
	}

	if ir != nil {
		// Unblind.
		m.Mul(m, ir)
		m.Mod(m, priv.N)
	}
	enc = m.Bytes()
	return
}

/**
  @param sourceBytes 原文
  @param publicKey  公钥字符串
  公钥加密
*/
func RsaEncWithPubKey(sourceBytes, pubKey []byte) ([]byte, error) {
	//获取公钥
	//block, _ := pem.Decode([]byte(publicKey))
	//if block == nil {
	//	return nil, errors.New("获取公钥失败")
	//}

	pubInterface, err := x509.ParsePKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	//声明密文,动态数组
	var cipherByte []byte
	//分段加密
	for i := 0; i < len(sourceBytes); i += 245 {
		var slice []byte
		if (i + 245) < len(sourceBytes) {
			slice = sourceBytes[i : i+245]
		} else {
			slice = sourceBytes[i:]
		}
		//Rsa加密,encryptBytes:分段密文
		encryptBytes, err := rsa.EncryptPKCS1v15(rand.Reader, pub, slice)
		if err != nil {
			return nil, err
		}
		//追加分段密文encryptBytes=>cipherByte
		cipherByte = append(cipherByte, encryptBytes...)
	}
	return cipherByte, nil
}

/**
	@param cipherByte 密文
	@param privateKey 私钥字符串
  	私钥解密,返回原文
*/
func RsaDecWithPriKey(cipherByte []byte, privateKey string) ([]byte, error) {
	//获取私钥
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//声明一个动态数组,用来存放解密之后的数据
	var source []byte
	//分段解密
	for i := 0; i < len(cipherByte); i += 256 {
		var slice []byte
		if (i + 256) < len(cipherByte) {
			slice = cipherByte[i : i+256]
		} else {
			slice = cipherByte[i:]
		}
		//rsa解密
		decrypt, err := rsa.DecryptPKCS1v15(rand.Reader, priv, slice)
		if err != nil {
			return nil, err
		}
		//追加解密数据decrypt=>source
		source = append(source, decrypt...)
	}
	return source, nil
}

/**
@param cipherText  待签名字段
RSA私钥签名sha1
*/
func RsaEncWithPKCS8(cipherText []byte, privateKey []byte) ([]byte, error) {
	//获取私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//返回签名结果
	return rsa.EncryptPKCS1v15(rand.Reader, &priv.(*rsa.PrivateKey).PublicKey, cipherText)
}

/**
@param cipherText  待签名字段
RSA私钥签名sha1
*/
func RsaSign(cipherText []byte, privateKey []byte) ([]byte, error) {
	//获取私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//指定HASH类型  SHA1
	h := crypto.Hash.New(crypto.SHA1)
	h.Write(cipherText)
	hashed := h.Sum(nil)
	//返回签名结果
	return rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA1, hashed)
}

/**
@param cipherText  待签名字段
RSA私钥签名sha1
*/
func RsaSignPKCS8(cipherText []byte, privateKey []byte) ([]byte, error) {
	//获取私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		logrus.Errorf("ParsePKCS1PrivateKey:%v", err)
		return nil, err
	}
	//指定HASH类型  SHA1
	h := crypto.Hash.New(crypto.SHA1)
	h.Write(cipherText)
	hashed := h.Sum(nil)
	//返回签名结果
	return rsa.SignPKCS1v15(rand.Reader, priv.(*rsa.PrivateKey), crypto.SHA1, hashed)
}

/**
@param origData 待签名字段
@param sign     签名
公钥验签
*/
func RsaVerifySign(origData []byte, sign []byte, publicKey string) error {
	//获取公钥
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	pub := pubInterface.(*rsa.PublicKey)
	//指定HASH类型  SHA1
	h := crypto.Hash.New(crypto.SHA1)
	h.Write(origData)
	hashed := h.Sum(nil)
	//返回验签结果
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA1, hashed, sign)
	return err
}

var (
	errPublicModulus       = errors.New("crypto/rsa: missing public modulus")
	errPublicExponentSmall = errors.New("crypto/rsa: public exponent too small")
	errPublicExponentLarge = errors.New("crypto/rsa: public exponent too large")
)

func encrypt_priv(priv *rsa.PrivateKey, c *big.Int) *big.Int {
	m := new(big.Int).Exp(c, priv.D, priv.N)
	return m
}

func checkPub(pub *rsa.PublicKey) error {
	if pub.N == nil {
		return errPublicModulus
	}
	if pub.E < 2 {
		return errPublicExponentSmall
	}
	if pub.E > 1<<31-1 {
		return errPublicExponentLarge
	}
	return nil
}

func copyWithLeftPad(dest, src []byte) {
	numPaddingBytes := len(dest) - len(src)
	for i := 0; i < numPaddingBytes; i++ {
		dest[i] = 0
	}
	copy(dest[numPaddingBytes:], src)
}

func encryptPKCS1v15_priv( /*rand io.Reader, */ priv *rsa.PrivateKey, msg []byte) ([]byte, error) {
	if err := checkPub(&priv.PublicKey); err != nil {
		return nil, err
	}
	k := (priv.N.BitLen() + 7) / 8
	if len(msg) > k-11 {
		return nil, rsa.ErrMessageTooLong
	}

	em := make([]byte, k)
	em[1] = 1
	for i := 2; i < len(em)-len(msg)-1; i++ {
		em[i] = 0xff
	}
	mm := em[len(em)-len(msg):]

	//em[1] = 2
	//ps, mm := em[2:len(em)-len(msg)-1], em[len(em)-len(msg):]
	//err := nonZeroRandomBytes(ps, rand)
	//if err != nil {
	//	return nil, err
	//}

	em[len(em)-len(msg)-1] = 0
	copy(mm, msg)

	m := new(big.Int).SetBytes(em)
	c := encrypt_priv(priv, m)

	copyWithLeftPad(em, c.Bytes())
	return em, nil
}

func decrypt_pub(c *big.Int, pub *rsa.PublicKey, m *big.Int) *big.Int {
	e := big.NewInt(int64(pub.E))
	c.Exp(m, e, pub.N)
	return c
}

func nonZeroRandomBytes(s []byte, rand io.Reader) (err error) {
	_, err = io.ReadFull(rand, s)
	if err != nil {
		return
	}

	for i := 0; i < len(s); i++ {
		for s[i] == 0 {
			_, err = io.ReadFull(rand, s[i:i+1])
			if err != nil {
				return
			}

			s[i] ^= 0x42
		}
	}

	return
}

func leftPad(input []byte, size int) (out []byte) {
	n := len(input)
	if n > size {
		n = size
	}
	out = make([]byte, size)
	copy(out[len(out)-n:], input)
	return
}

func decryptPKCS1v15_pub_(pub *rsa.PublicKey, ciphertext []byte) (valid int, em []byte, index int, err error) {
	k := (pub.N.BitLen() + 7) / 8
	if k < 11 {
		err = rsa.ErrDecryption
		return
	}

	c := new(big.Int).SetBytes(ciphertext)
	m := decrypt_pub(new(big.Int), pub, c)

	em = leftPad(m.Bytes(), k)
	firstByteIsZero := subtle.ConstantTimeByteEq(em[0], 0)
	secondByteIsTwo := subtle.ConstantTimeByteEq(em[1], 2)

	lookingForIndex := 1

	for i := 2; i < len(em); i++ {
		equals0 := subtle.ConstantTimeByteEq(em[i], 0)
		index = subtle.ConstantTimeSelect(lookingForIndex&equals0, i, index)
		lookingForIndex = subtle.ConstantTimeSelect(equals0, 0, lookingForIndex)
	}

	validPS := subtle.ConstantTimeLessOrEq(2+8, index)

	valid = firstByteIsZero & secondByteIsTwo & (^lookingForIndex & 1) & validPS
	index = subtle.ConstantTimeSelect(valid, index+1, 0)
	return valid, em, index, nil
}

func decryptPKCS1v15_pub(pub *rsa.PublicKey, ciphertext []byte) ([]byte, error) {
	if err := checkPub(pub); err != nil {
		return nil, err
	}
	valid, out, index, err := decryptPKCS1v15_pub_(pub, ciphertext)
	if err != nil {
		return nil, err
	}
	if valid == 0 {
		return nil, rsa.ErrDecryption
	}
	return out[index:], nil
}

func gen(bit int) {
	priv, err := rsa.GenerateKey(rand.Reader, bit)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fp, err := os.Create(fmt.Sprintf("priv%d.txt", bit))
	defer fp.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	pem.Encode(fp, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})

	fp, err = os.Create(fmt.Sprintf("pub%d.txt", bit))
	defer fp.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pubASN1, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	pem.Encode(fp, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
}

func doit() {
	gen(2048)
	//gen(896);
	buf, err := ioutil.ReadFile("priv2048.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	block, _ := pem.Decode(buf)

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	plain := "1234567889"
	encPub, err := rsa.EncryptPKCS1v15(rand.Reader, &priv.PublicKey, []byte(plain))
	fmt.Println("encPub:", hex.EncodeToString(encPub))
	encPriv, err := encryptPKCS1v15_priv(priv, []byte(plain))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("enc_priv: " + hex.EncodeToString(encPriv))

	decPub, err := decryptPKCS1v15_pub(&priv.PublicKey, encPriv)
	if err != nil {
		fmt.Println("decryptPKCS1v15_pub:" + err.Error())
		return
	}

	fmt.Println("dec_pub: " + hex.Dump(decPub))
}
