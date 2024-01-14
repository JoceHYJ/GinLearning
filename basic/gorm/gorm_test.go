package gorm

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"testing"
	"time"
)

type Product struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Number         string    `gorm:"unique" json:"number"`
	Category       string    `gorm:"type:varchar(256);not null" json:"category"`
	Name           string    `gorm:"type:varchar(20);not null" json:"name"`
	MadeIn         string    `gorm:"type:varchar(128);not null" json:"made_in"`
	ProductionTime time.Time `json:"production_time"`
}

type GormResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    any    `json:"data"`
}

var gormDB *gorm.DB

var gormResponse GormResponse

func init() {
	var err error
	sqlStr := "root:010729@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=true&loc=Local"
	gormDB, err = gorm.Open(mysql.Open(sqlStr), &gorm.Config{}) // 配置相中预设了连接池
	if err != nil {
		fmt.Println("数据库连接出现了问题：", err)
		return
	}
	err = gormDB.AutoMigrate(Product{})
	if err != nil {
		fmt.Println("创建数据库表格失败：", err)
	}

}

func TestGorm(t *testing.T) {
	r := gin.Default()
	r.POST("gorm/insert", gormInsertData)
	r.GET("gorm/get", gormGetData)
	r.GET("gorm/mulget", gormGetMulData)
	r.PUT("gorm/update", gormUpdateData)
	r.DELETE("gorm/delete", gormDeleteData)
	_ = r.Run(":8080")
}

func gormDeleteData(c *gin.Context) {
	number := c.Query("number")
	var count int64
	gormDB.Model(&Product{}).Where("number=?", number).Count(&count)
	if count <= 0 {
		handleGormError(c, errors.New("数据不存在"))
		return
	}
	tx := gormDB.Where("number=?", number).Delete(&Product{})
	if tx.RowsAffected > 0 {
		gormResponse.Code = http.StatusOK
		gormResponse.Message = "删除成功"
		gormResponse.Data = "OK"
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	handleGormError(c, tx.Error)
}

func gormUpdateData(c *gin.Context) {
	var p Product
	err := c.Bind(&p)
	if err != nil {
		handleGormError(c, err)
		return
	}
	var count int64
	gormDB.Model(&Product{}).Where("number=?", p.Number).Count(&count)
	if count <= 0 {
		handleGormError(c, errors.New("数据不存在"))
		return
	}
	tx := gormDB.Model(&Product{}).Where("number=?", p.Number).Updates(&p)
	if tx.RowsAffected > 0 {
		gormResponse.Code = http.StatusOK
		gormResponse.Message = "更新成功"
		gormResponse.Data = "OK"
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	handleGormError(c, tx.Error)
}

func gormGetMulData(c *gin.Context) {
	category := c.Query("category")
	products := make([]Product, 10)
	tx := gormDB.Where("category =?", category).Find(&products).Limit(10)
	if tx.Error != nil {
		handleGormError(c, tx.Error)
		return
	}
	gormResponse.Code = http.StatusOK
	gormResponse.Message = "读取成功"
	gormResponse.Data = products
	c.JSON(http.StatusOK, gormResponse)
}

func gormGetData(c *gin.Context) {
	number := c.Query("number")
	product := Product{}
	tx := gormDB.Where("number=?", number).First(&product)
	if tx.Error != nil {
		handleGormError(c, tx.Error)
		return
	}
	gormResponse.Code = http.StatusOK
	gormResponse.Message = "读取成功"
	gormResponse.Data = product
	c.JSON(http.StatusOK, gormResponse)
}

func gormInsertData(c *gin.Context) {
	var p Product
	err := c.Bind(&p)
	if err != nil {
		handleGormError(c, err)
		return
	}

	err = gormDB.Table("product").AutoMigrate(&p)
	if err != nil {
		panic("failed to auto migrate")
	}

	tx := gormDB.Create(&p)

	if tx.RowsAffected > 0 {
		gormResponse.Code = http.StatusOK
		gormResponse.Message = "写入成功"
		gormResponse.Data = "OK"
		c.JSON(http.StatusOK, gormResponse)
		return
	}
	fmt.Printf("insert failed, err:%v\n", err)
	handleGormError(c, tx.Error)
}

func handleGormError(c *gin.Context, err interface{}) {
	gormResponse.Code = http.StatusBadRequest
	gormResponse.Message = "出现异常错误"
	gormResponse.Data = err
	c.JSON(http.StatusOK, gormResponse)
}
