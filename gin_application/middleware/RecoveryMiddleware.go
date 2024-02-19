package middleware

import (
	"GinLearning/gin_application/response"
	"fmt"
	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware 恢复中间件，用于捕获panic异常
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.Fail(c, nil, fmt.Sprint(err))
				c.Abort()
				return
			}
		}()
	}
}
