package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// Product ---> 注意结构体名称 Product 而 创建的表的名称为 Products
type Product struct {
	ID             int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Number         string    `gorm:"unique" json:"number"`                       // 商品编号(唯一)
	Category       string    `gorm:"type:varchar(256);not null" json:"category"` // 商品类别
	Name           string    `gorm:"type:varchar(20);not null" json:"name"`      // 商品名称
	MadeIn         string    `gorm:"type:varchar(128);not null" json:"made_in"`  // 生产地
	ProductionTime time.Time `json:"production_time"`                            //  生产时间
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
	sqlStr := "root:010729@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	// 连接数据库
	gormDB, err = gorm.Open(mysql.Open(sqlStr), &gorm.Config{}) // 配置项预设了连接池 ConnPool
	if err != nil {
		fmt.Println("连接数据库出现了问题:", err)
		return
	}
}

// gorm 依赖: go get gorm.io/gorm
// 数据库驱动: go get gorm.io/driver/mysql
func main() {
	r := gin.Default()
	r.POST("gorm/insert", gormInsertData)
	r.GET("gorm/get", gormGetData)
	r.GET("gorm/mulget", gormGetMulData)
	r.PUT("gorm/update", gormUpdateData)
	r.DELETE("gorm/delete", gormDeleteData)
	r.Run(":8080")
}

// gormDeleteData 删除数据
func gormDeleteData(c *gin.Context) {
	number := c.Query("number")
	// 1.先查询
	var count int64
	gormDB.Model(&Product{}).Where("number=?", number).Count(&count)
	if count <= 0 {
		HandleResponse(c, http.StatusBadRequest, "数据不存在", "error")
	}
	// 2.再删除
	tx := gormDB.Where("number=?", number).Delete(&Product{})
	if tx.RowsAffected > 0 {
		HandleResponse(c, http.StatusOK, "删除成功", "OK")
	}
	HandleResponse(c, http.StatusBadRequest, "删除失败", tx)
	fmt.Println(tx) // 打印结果
}

// gormUpdateData 更新数据
func gormUpdateData(c *gin.Context) {
	var p Product
	err := c.Bind(&p)
	if err != nil {
		HandleResponse(c, http.StatusBadRequest, "参数错误", err)
	}
	// 1.先查询
	var count int64
	gormDB.Model(&Product{}).Where("number=?", p.Number).Count(&count)
	//fmt.Println("count:", count)
	if count <= 0 {
		HandleResponse(c, http.StatusBadRequest, "数据不存在", "error")
	}
	// 2.再更新
	tx := gormDB.Model(&Product{}).Where("number=?", p.Number).Updates(&p)
	if tx.RowsAffected > 0 {
		HandleResponse(c, http.StatusOK, "更新成功", "OK")
	}
	fmt.Printf("update failed, err:%v\n", err)
	HandleResponse(c, http.StatusBadRequest, "更新失败", tx)
	fmt.Println(tx) // 打印结果
}

// gormGetMulData 查询多条数据
func gormGetMulData(c *gin.Context) {
	category := c.Query("category")
	products := make([]Product, 10)
	tx := gormDB.Where("category = ?", category).Find(&products).Limit(10)
	if tx.Error != nil {
		HandleResponse(c, http.StatusBadRequest, "查询错误", tx.Error)
	}
	HandleResponse(c, http.StatusOK, "查询成功", products)
}

// gormGetData 查询单条数据
func gormGetData(c *gin.Context) {
	number := c.Query("number")
	product := Product{}
	tx := gormDB.Where("number=?", number).First(&product)
	if tx.Error != nil {
		HandleResponse(c, http.StatusBadRequest, "查询错误", tx.Error)
	}
	HandleResponse(c, http.StatusOK, "查询成功", product)
}

// gormInsertData 插入操作
func gormInsertData(c *gin.Context) {
	var p Product
	err := c.Bind(&p)
	if err != nil {
		HandleResponse(c, http.StatusBadRequest, "参数错误", err)
	}
	tx := gormDB.Create(&p)
	if tx.RowsAffected > 0 {
		HandleResponse(c, http.StatusOK, "写入成功", "OK")
	}
	fmt.Printf("insert failed, err:%v\n", err)
	HandleResponse(c, http.StatusBadRequest, "写入失败", tx)
	fmt.Println(tx) // 打印结果
}

func HandleResponse(c *gin.Context, code int, msg string, data any) {
	gormResponse.Code = code
	gormResponse.Message = msg
	gormResponse.Data = data
	c.JSON(http.StatusOK, gormResponse)
	return
}
