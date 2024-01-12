package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.POST("/upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, "uploading files failed")
		}
		files := form.File["file_key"] // 上传的所有文件
		dst := "files/"                // 保存文件路径
		// 遍历文件
		for _, file := range files {
			c.SaveUploadedFile(file, dst+file.Filename)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d个文件上传完成", len(files)))
	})
	router.Run(":9090")
}
