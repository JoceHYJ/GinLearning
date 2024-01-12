package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("GetOtherData", func(c *gin.Context) { //客户端发送请求
		//url := "http://www.baidu.com"
		url := "https://c-ssl.duitang.com/uploads/item/201602/28/20160228163528_Vr8tQ.jpeg"
		response, err := http.Get(url) // 服务端发送请求
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable) // 应答 client
			return
		}
		body := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")
		// 将数据写入响应体(返回客户端)
		c.DataFromReader(http.StatusOK, contentLength, contentType, body, nil)
	})
	router.Run(":9090")
}
