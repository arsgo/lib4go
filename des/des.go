package des

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
)

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
	dataBuffer = PKCS5Padding(dataBuffer, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, keyBuffer)
	crypted := make([]byte, len(dataBuffer))
	blockMode.CryptBlocks(crypted, dataBuffer)
	r = hex.EncodeToString(crypted)
	return
}

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
	buffer = ZeroUnPadding(buffer)
	r = hex.EncodeToString(buffer)
	return
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
