package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"time"
)

// Logrus 中的复杂功能，如: 设置保存最大时间、设置切割时间间隔等
// 文件切割依赖: go get github.com/lestrrat-go/file-rotatelogs
// hook 机制依赖: go get github.com/rifflock/lfshook

var (
	logFilePath = "./basic/24_logrus2/log/"
	logFileName = "system.log"
)

func main() {
	r := gin.Default()
	r.Use(logMiddleware())
	r.GET("/logrus2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "响应成功",
			"data": "OK",
		})
	})
	r.Run(":8080")
}

func logMiddleware() gin.HandlerFunc {
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	// 写入文件
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
	}
	// 实例化
	logger := logrus.New()
	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	// 设置输出
	logger.Out = file
	// 设置 rotatelogs 实现 log 文件分割
	logWriter, err := rotatelogs.New(
		// 分割后的文件名称
		fileName+".%Y%m%d.log",
		// 生成软连接指向最新的日志文件
		rotatelogs.WithLinkName(fileName),
		// 设置最大保存时间
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// 设置日志切割时间间隔(1天)
		rotatelogs.WithRotationTime(1*time.Hour),
	)
	// hook 机制的设置
	writerMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	// 给 logrus 添加 hook
	logger.AddHook(lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	}))
	return func(c *gin.Context) {
		c.Next()
		// 请求方式
		method := c.Request.Method
		// 请求路由
		reqUrl := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求 ip
		clientIP := c.ClientIP()
		logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"client_ip":   clientIP,
			"req_method":  method,
			"req_url":     reqUrl,
		}).Info()
	}
}
