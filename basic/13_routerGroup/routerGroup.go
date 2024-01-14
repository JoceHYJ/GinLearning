package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResGroup struct {
	Data string
	Path string
}

// 项目中有多个路由，通过路由组进行分组/分类管理
func main() {
	router := gin.Default()
	// 路由分组1
	v1 := router.Group("/v1") // 1级
	{
		r := v1.Group("/user")        // 2级
		r.GET("/login", login)        // 响应请求: /v1/user/Login
		r2 := r.Group("/showInfo")    // 3级
		r2.GET("/abstract", abstract) // 响应请求: /v1/user/showInfo/abstract
		r2.GET("/detail", detail)     // 响应请求: /v1/user/showInfo/detail
	}

	// 路由分组2
	v2 := router.Group("/v2") //1级
	{
		v2.GET("/other", other) // 响应请求: /v2/other
	}
	router.Run(":8080")
}

func other(c *gin.Context) {
	c.JSON(http.StatusOK, ResGroup{"Other", c.Request.URL.Path})
}

func detail(c *gin.Context) {
	c.JSON(http.StatusOK, ResGroup{"detail", c.Request.URL.Path})
}

func abstract(c *gin.Context) {
	c.JSON(http.StatusOK, ResGroup{"abstract", c.Request.URL.Path})
}

func login(c *gin.Context) {
	c.JSON(http.StatusOK, ResGroup{"login", c.Request.URL.Path})
}
