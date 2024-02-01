package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"net/http"
	"time"
)

// xorm 依赖: go get github.com/go-xorm/xorm
// 通过 xorm 进行数据库的 CRUD 操作

var x *xorm.Engine
var xormResponse XormResponse

// Stu 定义结构体(xorm 支持双向映射)：没有表会进行创建
type Stu struct {
	Id      int64     `xorm:"pk autoincr" json:"id"`
	StuNum  string    `xorm:"unique" json:"stu_num"`
	Name    string    `json:"name"`
	Age     int       `json:"age"`
	Created time.Time `xorm:"created" json:"created"`
	Updated time.Time `xorm:"updated" json:"updated"`
}

// XormResponse 应答 Client 请求
type XormResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func init() {
	sqlStr := "root:010729@tcp(127.0.0.1:3306)/xorm?charset=utf8mb4&parseTime=true&loc=Local" // xorm: 数据库名称
	var err error
	// 1、创建数据库引擎
	x, err = xorm.NewEngine("mysql", sqlStr)
	if err != nil {
		fmt.Println("数据库连接失败:", err)
	}
	// 2、创建或同步表 Stu
	err = x.Sync(new(Stu))
	if err != nil {
		fmt.Println("数据库同步失败:", err)
	}
}

func main() {
	r := gin.Default()
	// 数据库的 CRUD ---> gin 的 POST GET PUT DELETE 方法
	r.POST("xorm/insert", xormInsertData)
	r.GET("xorm/get", xormGetData)
	r.Run(":8080")
}

// xormGetData 查询操作(单条操作)
func xormGetData(c *gin.Context) {
	stuNum := c.Query("stu_num")
	var stus []Stu
	err := x.Where("stu_num=?", stuNum).Find(&stus)
	if err != nil {
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "查询错误"
		xormResponse.Data = "error"
		c.JSON(http.StatusOK, xormResponse)
		return
	}
	xormResponse.Code = http.StatusOK
	xormResponse.Message = "查询成功"
	xormResponse.Data = stus
	c.JSON(http.StatusOK, xormResponse)
}

// xormInsertData 插入操作
func xormInsertData(c *gin.Context) {
	var s Stu
	err := c.Bind(&s)
	if err != nil {
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "参数错误"
		xormResponse.Data = "error"
		c.JSON(http.StatusOK, xormResponse)
		return
	}
	// affected：受影响记录行数
	affected, err := x.Insert(s)
	if err != nil || affected <= 0 {
		fmt.Printf("insert failed, err:%v\n", err)
		xormResponse.Code = http.StatusBadRequest
		xormResponse.Message = "写入失败"
		xormResponse.Data = err
		c.JSON(http.StatusOK, xormResponse)
		return
	}
	xormResponse.Code = http.StatusOK
	xormResponse.Message = "写入成功"
	xormResponse.Data = "OK"
	c.JSON(http.StatusOK, xormResponse)
	fmt.Println(affected) // 打印结果
}
