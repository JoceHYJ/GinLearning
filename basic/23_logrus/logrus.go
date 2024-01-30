package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var log = logrus.New() // 创建一个 logrus 示例

// initLogrus 初始化函数
func initLogrus() error {
	log.Formatter = &logrus.JSONFormatter{} // 设置为 json 格式的日志
	file, err := os.OpenFile("./basic/23_logrus/gin_log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("创建文件/打开文件失败！")
		return err
	}
	log.Out = file               // 设置 log 的默认文件输出
	gin.SetMode(gin.ReleaseMode) // 发布版本
	gin.DefaultWriter = log.Out  // gin 框架的日志也输出到 log 的默认文件中
	log.Level = logrus.InfoLevel // 设置日志级别
	return nil
}

// 通过 logrus 记录日志
func main() {
	err := initLogrus()
	if err != nil {
		fmt.Println(err)
		return
	}
	r := gin.Default()
	r.GET("/logrus", func(c *gin.Context) {
		log.WithFields(logrus.Fields{
			"url":    c.Request.RequestURI,
			"method": c.Request.Method,
			"params": c.Query("name"),
			"IP":     c.ClientIP(),
		}).Info()
		resData := struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			Data any    `json:"data"`
		}{http.StatusOK, "响应成功", "OK"}
		c.JSON(http.StatusOK, resData)
	})
	r.Run(":8080")
}

// logrus 依赖: go get github.com/sirupsen/logrus
