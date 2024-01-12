package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("fileName")
		if err != nil {
			c.String(http.StatusBadRequest, "file uploading failed")
		}
		dst := "/home/jocehyj/goWorkspace/src/Learning/GinLearning/07_fileUpload/"
		c.SaveUploadedFile(file, dst+file.Filename) // 文件上传方法
		c.String(http.StatusOK, fmt.Sprintf("%s uploading succeed", file.Filename))
	})
	router.Run(":9090")
}
