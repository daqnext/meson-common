package httputils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/daqnext/meson-common/common/accountmgr"
	"github.com/daqnext/meson-common/common/logger"
	"github.com/daqnext/meson-common/common/resp"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
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

func SendGetRequest(url string, param req.Param, authorizationToken string, timeoutSecond int) (*req.Resp, error) {
	var authHeader req.Header = nil
	if authorizationToken != "" {
		authHeader = req.Header{
			"Accept":        "application/json",
			"Authorization": "Basic " + accountmgr.Token,
		}
	}

	r := req.New()
	if timeoutSecond > 0 {
		r.SetTimeout(time.Duration(timeoutSecond) * time.Second)
	}

	response, err := r.Get(url, param, authHeader)
	if err != nil {
		logger.Error("request error", "err", err)
		return nil, err
	}
	return response, nil
}

func SendPostRequest(url string, param req.Param, body interface{}, authorizationToken string, timeoutSecond int) (*req.Resp, error) {
	var authHeader req.Header = nil
	if authorizationToken != "" {
		authHeader = req.Header{
			"Accept":        "application/json",
			"Authorization": "Basic " + accountmgr.Token,
		}
	}

	r := req.New()
	if timeoutSecond > 0 {
		r.SetTimeout(time.Duration(timeoutSecond) * time.Second)
	}

	response, err := r.Post(url, param, authHeader, req.BodyJSON(body))
	if err != nil {
		logger.Error("request error", "err", err)
		return nil, err
	}
	return response, nil
}

func GenResponseStruct(v interface{}) *resp.RespBody {
	return &resp.RespBody{
		Data: v,
	}
}

func HandleResponse(response *req.Resp, v *resp.RespBody) (httpStatusCode int, err error) {
	httpResponse := response.Response()
	httpStatusCode = httpResponse.StatusCode

	err = response.ToJSON(v)
	if err != nil {
		logger.Error("response.ToJSON error", "err", err)
		return httpStatusCode, err
	}
	return httpStatusCode, nil
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
