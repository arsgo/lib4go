package net

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

//HTTPClient HTTP客户端
type HTTPClient struct {
	client *http.Client
}
type HTTPClientRequest struct{
	headers map[string]string
	client *http.Client
	method string
	url string 
	params string
}
//NewHTTPClientCert 根据pem证书初始化httpClient
func NewHTTPClientCert(certFile string, keyFile string, caFile string) (client *HTTPClient, err error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}
	caData, err := ioutil.ReadFile(caFile)
	if err != nil {
		return
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)
	ssl := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	ssl.Rand = rand.Reader
	client = &HTTPClient{}
	client.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: ssl,
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
func (c *HTTPClient) NewRequest(method string, url string)(*HTTPClientRequest){
	request:=&HTTPClientRequest{}
	request.client=c.client
	request.headers=make(map[string]string)
	request.method=strings.ToUpper(method)
	request.params=""
	request.url=url
	return request
}
func (c *HTTPClientRequest) SetData(params string){
	c.params=params
}
func (c *HTTPClientRequest) SetHeader(key string,value string){
	c.headers[key]=value
}

//Request 发送http请求, method:http请求方法包括:get,post,delete,put等 url: 请求的HTTP地址,不包括参数,params:请求参数,
//header,http请求头多个用/n分隔,每个键值之前用=号连接
func (c *HTTPClientRequest) Request() (content string, status int, err error) {
	req, err := http.NewRequest(c.method,c.url, strings.NewReader(c.params))
	if err != nil {
		return
	}
	for i,v:=range c.headers{
		req.Header.Set(i,v)
	}
	resp, err := c.client.Do(req)
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
