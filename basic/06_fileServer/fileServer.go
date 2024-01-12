package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/file", fileServer)
	router.Run(":9090")
}

func fileServer(c *gin.Context) {
	path := "/home/jocehyj/goWorkspace/src/Learning/GinLearing/06_fileServer/"
	fileName := path + c.Query("name")
	c.File(fileName)
}
