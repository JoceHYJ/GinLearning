package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pjebs/restgate"
	"net/http"
)

func main() {
	r := gin.Default()
	// 通过中间件添加安全认证
	r.Use(authMiddleware())
	r.GET("/auth1", func(c *gin.Context) {
		resData := struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			Data any    `json:"data"`
		}{http.StatusOK, "验证通过", "OK"}
		c.JSON(http.StatusOK, resData)
	})
	r.Run(":8080")
}

// 实现安全认证中间件
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		gate := restgate.New(
			"X-Auth-Key",
			"X-Auth-Secret",
			restgate.Static,
			restgate.Config{
				Key:                []string{"admin", "gin"},
				Secret:             []string{"adminpw", "gin_ok"},
				HTTPSProtectionOff: true, // 默认HTTPS -> 使用 HTTP 设为 true
			})
		// 标志位
		nextCalled := false
		nextAdapter := func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			c.Next()
		}
		gate.ServeHTTP(c.Writer, c.Request, nextAdapter)
		if nextCalled == false {
			c.AbortWithStatus(http.StatusUnauthorized)
		} // false 终止
	}
}

// 中间件: go get github.com/pjebs/restgate
// 2024/01/30 17:16:28 WARNING: HTTPS Protection is off. This is potentially insecure! 可以参考 20_https进行修改
