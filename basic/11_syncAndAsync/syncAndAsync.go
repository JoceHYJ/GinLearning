package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// Async: 在 Milddleware 或处理程序中启动新的 Goroutines (需要使用 context 的副本)
func main() {
	r := gin.Default()
	// 同步
	r.GET("/sync", func(c *gin.Context) {
		sync(c)
		c.JSON(200, ">>>主程序(goroutine)同步已执行<<<")
	})
	// 异步
	r.GET("/async", func(c *gin.Context) {
		for i := 0; i < 6; i++ {
			cCp := c.Copy()
			go async(cCp, i)
		}
		c.JSON(http.StatusOK, ">>>>主程序(goroutine)异步已执行<<<<")
	})
	r.Run(":8080")
}

func async(cp *gin.Context, i int) {
	fmt.Println("第" + strconv.Itoa(i) + "个goroutine开始执行:" + cp.Request.URL.Path)
	time.Sleep(time.Second * 3)
	fmt.Println("第" + strconv.Itoa(i) + "个goroutine执行结束")
}

func sync(c *gin.Context) {
	println("开始执行同步任务:" + c.Request.URL.Path)
	time.Sleep(time.Second * 3)
	println("同步任务执行完成!")
}
