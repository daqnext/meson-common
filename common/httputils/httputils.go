package httputils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

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
