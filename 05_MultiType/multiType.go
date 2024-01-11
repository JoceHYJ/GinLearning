// Gin 框架的多形式渲染
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	// json 数据格式渲染
	router.GET("/json", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"html": "<b>hello, tomato</b>",
		})
	})

	// 原样输出 html (html 渲染)
	router.GET("/someHTML", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, gin.H{
			"html": "<b>hello, tomato</b>",
		})
	})

	// 输出 xml 形式 ( XML 渲染)
	router.GET("/someXML", func(c *gin.Context) {
		type Message struct {
			Name string
			Msg  string
			Age  int
		}
		info := Message{}
		info.Name = "tomato"
		info.Msg = "happy birthday"
		info.Age = 24
		c.XML(http.StatusOK, info)
	})

	// 输出 yaml 形式 ( YAML 渲染)
	router.GET("/someYAML", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{
			"message": "Happy birthday my tomato",
			"status":  200,
		})
	})

	// 开启服务
	router.Run(":9090")
}
