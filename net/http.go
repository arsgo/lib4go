package net

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

//HTTPClient HTTP客户端
type HTTPClient struct {
	client *http.Client
}

//NewHTTPClient 构建HTTP客户端，用于发送GET POST等请求
func NewHTTPClient() (client *HTTPClient) {
	client = &HTTPClient{}
	client.client = &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				c, err := net.DialTimeout(netw, addr, 0)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			MaxIdleConnsPerHost:   0,
			ResponseHeaderTimeout: 0,
		},
	}
	return
}

//Get http get请求
func (c *HTTPClient) Get(url string) (content string, status int, err error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	content = string(body)
	status = resp.StatusCode
	return
}

//Post http Post请求
func (c *HTTPClient) Post(url string, params string) (content string, status int, err error) {
	resp, err := c.client.Post(url, "application/x-www-form-urlencoded", strings.NewReader(params))
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	content = string(body)
	status = resp.StatusCode
	return
}
