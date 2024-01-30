package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // 需要添加数据库的驱动
	"github.com/pjebs/restgate"
	"net/http"
)

// restgate动态配置：不将接口对应的用户名和密码配置到代码中，而是配置在数据库中
func main() {
	r := gin.Default()
	r.Use(authMiddleware2())
	r.GET("/auth2", func(c *gin.Context) {
		resData := struct {
			Code int    `json:"code"`
			Data string `json:"data"`
			Msg  string `json:"msg"`
		}{http.StatusOK, "OK", "验证成功"}
		c.JSON(http.StatusOK, resData)
	})
	r.Run(":8080")
}

var db *sql.DB

func init() {
	db, _ = sqlDB()
}

func sqlDB() (*sql.DB, error) {
	DB_TYPE := "mysql"
	DB_HOST := "localhost"
	DB_NAME := "api_secure"
	DB_USER := "root"
	DB_PASSWORD := "010729"
	DB_PORT := "3306"
	openString := DB_USER + ":" + DB_PASSWORD + "@tcp(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME
	db, err := sql.Open(DB_TYPE, openString)
	return db, err
}

func authMiddleware2() gin.HandlerFunc {
	return func(c *gin.Context) {
		gate := restgate.New(
			"X-Auth-key",
			"X-Auth-Secret",
			restgate.Database,
			restgate.Config{
				DB:                 db,
				TableName:          "users",
				Key:                []string{"keys"},
				Secret:             []string{"secrets"},
				HTTPSProtectionOff: true,
			})
		nextCalled := false
		nextAdapter := func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			c.Next()
		}
		gate.ServeHTTP(c.Writer, c.Request, nextAdapter)
		if nextCalled == false {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
