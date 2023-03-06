package resp

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
)

const (
	SuccessCode           = 200
	InvalidParamCode      = 400
	InvalidPermissionCode = 403
	InternalErrorCode     = 500
	NotLoginCode          = 805
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(200, &Response{Code: SuccessCode, Msg: "Success", Data: data})
}

func Error(c *gin.Context, code int, err error, msg ...string) {
	c.JSON(200, &Response{Code: code, Msg: err.Error(), Data: gin.H{}})
	log.Printf("[ERROR] method:%v, url:%v, err:%v(%v)", c.Request.Method, c.Request.URL.Path, err, msg)
	c.Abort()
}

func InvalidParam(c *gin.Context) {
	Error(c, InvalidParamCode, errors.New("invalid param"))
}

func InvalidPermission(c *gin.Context) {
	Error(c, InvalidParamCode, errors.New("invalid permission"))
}

func NotLogin(c *gin.Context) {
	Error(c, NotLoginCode, errors.New("user not login"))
}

func InternalError(c *gin.Context, err error) {
	Error(c, InternalErrorCode, err)
}

// DefaultInternalError Show "internal errors" instead of specific errors
func DefaultInternalError(c *gin.Context, err error) {
	Error(c, InternalErrorCode, errors.New("internal error"), err.Error())
}
