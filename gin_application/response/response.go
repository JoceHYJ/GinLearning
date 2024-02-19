package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 封装响应体
func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	ctx.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

// Success 响应成功
func Success(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, http.StatusOK, data, msg)
}

// Fail 响应失败
func Fail(ctx *gin.Context, data gin.H, msg string) {
	Response(ctx, http.StatusOK, http.StatusBadRequest, data, msg)
}
