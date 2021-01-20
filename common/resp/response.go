package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RespBody struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}

func ErrorResp(c *gin.Context, err *httpErr) {
	c.JSON(http.StatusOK, gin.H{
		"status": err.Code(),
		"msg":    err.Error(),
	})
}

func SuccessResp(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": success,
		"data":   data,
	})
}
