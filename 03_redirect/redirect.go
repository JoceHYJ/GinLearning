package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	// 一般重定向：重定向至外部网络
	router.GET("/redirect1", func(c *gin.Context) {
		url := "https://www.baidu.com"
		c.Redirect(http.StatusMovedPermanently, url)
	})

	// 路由重定向：重定向到具体的路由
	router.GET("/redirect2", func(c *gin.Context) {
		c.Request.URL.Path = "/TestRedirect"
		router.HandleContext(c)
	})
	// 路由：127.0.0.1：9090/TestRedirect
	router.GET("/TestRedirect", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "响应成功",
		})
	})
	router.Run(":9090")
}
