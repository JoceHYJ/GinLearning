package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

var sessionName string
var sessionValue string

type MyOption struct {
	sessions.Options
}

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte("session_secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/session", func(c *gin.Context) {
		name := c.Query("name")
		if len(name) <= 0 {
			c.JSON(http.StatusBadRequest, "数据错误")
			return
		}
		sessionName = "session_" + name
		sessionValue = "session_value_" + name
		session := sessions.Default(c)
		sessionData := session.Get(sessionName)
		if sessionData != sessionValue {
			// 保存 session 值
			session.Set(sessionName, sessionValue)
			o := MyOption{}
			o.Path = "/"
			o.MaxAge = 3600
			session.Options(o.Options)
			session.Save()
			c.JSON(http.StatusOK, "首次访问，session 已经保存")
		}
		c.JSON(http.StatusOK, "访问成功， 您的 session 是:`"+sessionData.(string))
	})
	r.Run(":8080")
}
