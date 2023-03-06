package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lllllan-fv/gateway-admin/public/resp"
)

// Recovery Capture all panic and return error message
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				resp.InternalError(c, c.Err())
			}
		}()
		c.Next()
	}
}
