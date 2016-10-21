package des

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"strings"
)

func Encrypt(input string, skey string) (r string, err error) {
	origData := []byte(input)
	key := []byte(skey)
	block, err := des.NewCipher(key)
	if err != nil {
		return
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	r = strings.ToUpper(hex.EncodeToString(crypted))
	return
}

func Decrypt(input string, skey string) (r string, err error) {
	crypted, err := hex.DecodeString(input)
	if err != nil {
		return
	}
	key := []byte(skey)
	block, err := des.NewCipher(key)
	if err != nil {
		return
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	r = string(origData)
	return
}

//Encrypt3DES
func Encrypt3DES(input string, skey string) (r string, err error) {
	origData := []byte(input)
	key := []byte(skey)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	r = strings.ToUpper(hex.EncodeToString(crypted))
	return
}

//Decrypt3DES TripleDesDecrypt  3DES解密
func Decrypt3DES(input, skey string) (r string, err error) {
	crypted, err := hex.DecodeString(input)
	if err != nil {
		return
	}
	key := []byte(skey)
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return
	}
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	r = string(origData)
	return
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS7Padding(data []byte) []byte {
	blockSize := 16
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)

}

/**
 *  去除PKCS7的补码
 */
func UnPKCS7Padding(data []byte) []byte {
	length := len(data)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}
