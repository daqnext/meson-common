package httputils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func RequestWithTimeOut(method string, url string, payload interface{}, header map[string]string, cTimeout time.Duration, rwTimeout time.Duration) ([]byte, error) {
	var bytesData []byte = nil
	var err error = nil
	if payload != nil {
		bytesData, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	request, _ := http.NewRequest(method, url, bytes.NewBuffer(bytesData))
	if header != nil {
		for k, v := range header {
			request.Header.Add(k, v)
		}
	}

	//http client
	client := &http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(cTimeout, rwTimeout),
		},
	}

	res, err := client.Do(request)
	if err != nil {
		//fmt.Println("Fatal error ", err.Error())
		return nil, err
	}
	if res.Status != "200 OK" {
		return nil, errors.New("Status:" + res.Status)
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}
	return content, nil
}

func Request(method string, url string, payload interface{}, header map[string]string) ([]byte, error) {
	var bytesData []byte = nil
	var err error = nil
	if payload != nil {
		bytesData, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	request, _ := http.NewRequest(method, url, bytes.NewBuffer(bytesData))
	if header != nil {
		for k, v := range header {
			request.Header.Add(k, v)
		}
	}
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		//fmt.Println("Fatal error ", err.Error())
		return nil, err
	}
	if res.Status != "200 OK" {
		return nil, errors.New("Status:" + res.Status)
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}
	return content, nil
}

func ForwardRequest(ctx *gin.Context, scheme string, host string, path string) {
	simpleHostProxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = scheme
			req.URL.Host = host
			req.URL.Path = path
			req.Host = host
		},
	}
	simpleHostProxy.ServeHTTP(ctx.Writer, ctx.Request)
}
