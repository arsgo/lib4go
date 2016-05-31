package des

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
)

//Encrypt 加密
func Encrypt(input string, key string) (r string, err error) {
	dataBuffer, err := hex.DecodeString(input)
	if err != nil {
		return
	}
	keyBuffer := []byte(key)
	block, err := des.NewCipher(keyBuffer)
	if err != nil {
		return
	}
	dataBuffer = pKCS5Padding(dataBuffer, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, keyBuffer)
	crypted := make([]byte, len(dataBuffer))
	blockMode.CryptBlocks(crypted, dataBuffer)
	r = hex.EncodeToString(crypted)
	return
}

//Decrypt 解密
func Decrypt(input string, key string) (r string, err error) {
	dataBuffer, err := hex.DecodeString(input)
	if err != nil {
		return
	}
	keyBuffer := []byte(key)
	block, err := des.NewCipher(keyBuffer)
	if err != nil {
		return
	}

	blockMode := cipher.NewCBCDecrypter(block, keyBuffer)
	var buffer []byte
	blockMode.CryptBlocks(buffer, dataBuffer)
	buffer = zeroUnPadding(buffer)
	r = hex.EncodeToString(buffer)
	return
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func zeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
