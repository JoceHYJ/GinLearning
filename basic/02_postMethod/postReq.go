package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.POST("/post", postMsg)
	router.Run(":9090")
}

func postMsg(c *gin.Context) {
	//name := c.Query("name") // 获取的是 url 中的数据
	name := c.DefaultPostForm("name", "Gin") // 获取 form 中的值，返回一个值
	fmt.Println(name)
	form, b := c.GetPostForm("name") // 获取 form 中的值，返回两个值(该值以及一个布尔值，指示该字段是否存在)
	fmt.Println(form, b)
	c.JSON(http.StatusOK, "welcome:"+name)
}
