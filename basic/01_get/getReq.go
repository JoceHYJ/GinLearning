package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()    // 路由引擎
	router.GET("/get", getMsg) // 通过路由引擎开启 get 请求，get 方法请求的数据放在 url 中
	//router.Run("127.0.0.1:9090")
	router.Run(":9090") // 开启服务
}

func getMsg(c *gin.Context) {
	name := c.Query("name")
	//c.String(http.StatusOK,"欢迎您: %s", name) // 返回 String 类型的数据
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,     //状态
		"msg":  "return msg",      // 描述信息
		"data": "welcome:" + name, // 数据
	}) // 返回 json 类型的数据
}
