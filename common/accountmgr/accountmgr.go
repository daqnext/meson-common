package accountmgr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/daqnext/meson-common/common/logger"
	"github.com/daqnext/meson-common/common/resp"
)

var Token string

func SLogin(url string, username string, password string) {
	postData := make(map[string]string)
	postData["username"] = username
	postData["password"] = password
	bytesData, _ := json.Marshal(postData)

	res, err := http.Post(
		url,
		"application/json;charset=utf-8",
		bytes.NewBuffer(bytesData),
	)
	if err != nil {
		fmt.Println("Login failed Fatal error ", err.Error())
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Fatal("Login failed Fatal error ", err.Error())
	}
	logger.Debug("response form server", "response string", string(content))
	var respBody resp.RespBody
	if err := json.Unmarshal(content, &respBody); err != nil {
		logger.Error("response from terminal unmarshal error", "err", err)
		logger.Fatal("response from terminal err")
		return
	}

	switch respBody.Status {
	case 0:
		Token = respBody.Data.(string)
		logger.Debug("login success! ", "token", Token)
		logger.Info("login success! Terminal start...")
	case 2101: //username not exist
		logger.Fatal("username not exist,please provide a correct username")
	case 2004: //username or password error
		logger.Fatal("username or password error,please provide a correct username and password")
	default:
		logger.Fatal("server error")
	}
}
