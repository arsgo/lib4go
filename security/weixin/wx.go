package weixin

import (
	"encoding/json"
)

//Decrypt 解密请求报文
func Decrypt(content string) (r string, err error) {
	response, err := parseEncryptTextRequestBody([]byte(content))
	if err != nil {
		return
	}
	buff, err := json.Marshal(response)
	if err != nil {
		return
	}
	r = string(buff)
	return
}

//Encrypt 加密响应报文
func Encrypt(fromUserName, toUserName, content, nonce, timestamp string) (r string, err error) {
	buffer, err := makeEncryptResponseBody(fromUserName, toUserName, content, nonce, timestamp)
	if err != nil {
		return
	}
	r = string(buffer)
	return
}

//MakeSign 构建签名
func MakeSign(timestamp, nonce, msgEncrypt string) string {
	return makeMsgSignature(timestamp, nonce, msgEncrypt)
}
