package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Login struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remark   string `json:"remark"`
}

// 将 Client 提供的 json 数据与 Server 对应的对象(实体/结构体)进行关联
// Gin 的 bind 方法： 将结构体与请求的参数进行绑定(请求参数 json 对应的 key 就是结构体对应的字段)
func main() {
	r := gin.Default()
	r.POST("/login", func(c *gin.Context) {
		var login Login
		err := c.Bind(&login) // Bind
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  "binding failed",
				"data": err.Error(),
			})
			return
		}
		if login.UserName == "user" && login.Password == "123456" {
			c.JSON(http.StatusOK, gin.H{
				"msg":  "Login succeed",
				"data": "OK",
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg":  "Login failed",
				"data": "error",
			})
			return
		}
	})
	r.Run(":8080")
}
