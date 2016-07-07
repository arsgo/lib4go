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

//Request 发送http请求, method:http请求方法包括:get,post,delete,put等 url: 请求的HTTP地址,不包括参数,params:请求参数,
//header,http请求头多个用/n分隔,每个键值之前用=号连接
func (c *HTTPClient) Request(method string, url string, params string, headers string) (content string, status int, err error) {
	req, err := http.NewRequest(method, url, strings.NewReader(params))
	if err != nil {
		return
	}
	hd := strings.Split(headers, "\n")
	for _, v := range hd {
		h := strings.Split(v, "=")
		if len(h) != 2 {
			continue
		}
		req.Header.Set(h[0], h[1])
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
