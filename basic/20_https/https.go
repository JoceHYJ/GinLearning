package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"net/http"
)

type HttpRes struct {
	Code   int    `json:"code"`
	Result string `json:"result"`
}

func main() {
	r := gin.Default()
	r.Use(httpsHandler())
	r.GET("/https_test", func(c *gin.Context) {
		fmt.Println(c.Request.Host)
		c.JSON(http.StatusOK, HttpRes{
			Code:   http.StatusOK,
			Result: "测试成功",
		})
	})

	path := "/home/jocehyj/goWorkspace/src/Learning/GinLearning/basic/CA/"
	r.RunTLS(":8080", path+"ca.crt", path+"ca.key") // 开启 HTTPS 服务
}

func httpsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		secureMiddle := secure.New(secure.Options{ // 创建一个新的安全中间件
			//SSLHost:     true, // 仅允许 HTTPS 请求
			SSLRedirect:        true,    // 将所有 HTTP 请求重定向到 HTTPS
			STSSeconds:         1536000, // 设置 STS 头的 max-age 为 1536000 秒（18 天）
			STSPreload:         true,    // 将网站添加到 STS 预加载列表中
			FrameDeny:          true,    // 防止页面被嵌入到其他网站中
			ContentTypeNosniff: true,    // 阻止浏览器猜测内容类型
			BrowserXssFilter:   true,    // 启用浏览器 XSS 过滤器
		})
		err := secureMiddle.Process(context.Writer, context.Request)
		// 如果不安全，终止
		if err != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, "数据不安全") //中间件进行拦截
			return
		}
		// 如果重定向，终止
		if status := context.Writer.Status(); status > 300 && status < 399 {
			context.Abort()
			return
		}
		// 安全，向下执行
		context.Next()
	}
}

// 中间件：github.com/unrolled/secure
