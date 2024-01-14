package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 登陆验证 Middleware 完成登陆验证
func main() {
	r := gin.Default()
	r.Use(AuthMiddleWare())
	r.GET("/login", func(c *gin.Context) {
		// 获取用户名称，是由 BasicAuth 中间件设置的
		user := c.MustGet(gin.AuthUserKey).(string)
		c.JSON(http.StatusOK, "登陆成功！"+"欢迎您: "+user)
	})
	r.Run(":8080")
}

func AuthMiddleWare() gin.HandlerFunc {
	// 静态添加用户名和密码
	accounts := gin.Accounts{
		"tomato": "pw123456",
		"sprite": "pw012345",
	}
	// 动态添加用户名和密码
	accounts["go"] = "1234555"
	accounts["gin"] = "12222000"

	// 将用户添加到登陆中间件中
	auth := gin.BasicAuth(accounts)
	return auth
}
