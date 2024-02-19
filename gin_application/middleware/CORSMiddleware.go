package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 设置允许访问源
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 预检结果缓存时间
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		// 允许请求类型
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		// 允许请求头字段
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		// 是否携带cookie
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusOK)
		} else {
			ctx.Next()
		}
	}
}
