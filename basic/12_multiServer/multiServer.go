package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

// 定义路由组
var g errgroup.Group

// 实现多服务器程序
func main() {
	// Server1
	server01 := &http.Server{
		Addr:         ":8081",
		Handler:      router01(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Server2
	server02 := &http.Server{
		Addr:         ":8082",
		Handler:      router02(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 开启服务
	g.Go(func() error {
		return server01.ListenAndServe()
	})
	g.Go(func() error {
		return server02.ListenAndServe()
	})
	// 阻塞主 Goroutine ---> 处于等待状态 ---> 执行 server01/2 对应的监听
	if err := g.Wait(); err != nil {
		fmt.Println("执行失败:", err)
	}
}

func router02() http.Handler {
	r2 := gin.Default()
	r2.GET("/myServer02", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "msg_server02",
		},
		)
	})
	return r2
}

func router01() http.Handler {
	r1 := gin.Default()
	r1.GET("/myServer01", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "msg_server01",
		},
		)
	})
	return r1
}
