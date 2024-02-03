package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
	"strings"
	"time"
	"xorm.io/xorm"
)

// xorm 依赖: go get xorm.io/xorm
// 通过 xorm 进行数据库的 CRUD 操作

var x *xorm.Engine
var xormResponse XormResponse

type XormValue struct {
	Xorm *xorm.Engine
}

type OrmValue[T any] struct {
	//Cfg Config
	Key    string
	sqlStr string
	XVal   XormValue
}

type Repository interface {
	InsertData(ctx *gin.Context)
}

var _ Repository = &OrmValue[any]{}

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

func (orm *OrmValue[T]) GetDB() error {
	var err error
	// 1、创建数据库引擎
	x, err = xorm.NewEngine("mysql", orm.sqlStr)
	if err != nil {
		fmt.Println("数据库连接失败:", err)
	}
	// 2、创建或同步表 Stu
	err = x.Sync(new(T))
	if err != nil {
		fmt.Println("数据库同步失败:", err)
	}
	orm.XVal.Xorm = x

	return err
}

type Config struct {
	DB struct {
		MySQL string `yaml:"mysql"`
	} `yaml:"db"`
	Business struct {
		Key string `yaml:"key"`
	} `yaml:"business"`
}

func (orm *OrmValue[T]) GetCfg() error {
	// 读取YAML文件
	yamlFile, err := os.ReadFile("config/cfg.yaml")
	if err != nil {
		fmt.Printf("无法读取YAML文件：%v", err)
		return err
	}

	// 解析YAML文件
	var cfg Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		fmt.Printf("无法解析YAML文件：%v", err)
		return err
	}

	orm.Key = cfg.Business.Key
	orm.sqlStr = cfg.DB.MySQL
	return nil
}

func main() {
	r := gin.Default()

	db := OrmValue[Stu]{}
	err := db.GetCfg()
	if err != nil {
		panic(err)
	}

	err = db.GetDB()
	if err != nil {
		panic(err)
	}

	//db.Key = "stu_num"

	// 数据库的 CRUD ---> gin 的 POST GET PUT DELETE 方法
	r.POST("xorm/insert", db.InsertData)
	r.GET("xorm/get", db.GetData)
	r.GET("xorm/mulget", XormGetMulData[Stu])
	r.PUT("xorm/update", XormUpdateData)
	r.DELETE("xorm/delete", XormDeleteData[Stu])
	r.Run(":8080")
}

// InsertData 插入操作
func (orm *OrmValue[T]) InsertData(c *gin.Context) {
	var typ T
	err := c.Bind(&typ)
	if err != nil {
		HandleResponse(c, http.StatusBadRequest, "参数错误", "error")
		return
	}
	// affected：受影响记录行数
	affected, err := orm.XVal.Xorm.Insert(typ) // TODO 解耦
	if err != nil || affected <= 0 {
		fmt.Printf("insert failed, err:%v\n", err)
		HandleResponse(c, http.StatusBadRequest, "写入失败", err)
		return
	}
	HandleResponse(c, http.StatusOK, "写入成功", "OK")
	fmt.Println(affected) // 打印结果
}

// XormDeleteData 删除操作
func XormDeleteData[T any](c *gin.Context) {
	stuNum := c.Query("stu_num")
	// 1、先查询
	//var stus []Stu
	var typs []T
	err := x.Where("stu_num=?", stuNum).Find(&typs)
	if err != nil || len(typs) <= 0 {
		HandleResponse(c, http.StatusBadRequest, "数据不存在", "error")
		return
	}
	// 2、再删除
	affected, err := x.Where("stu_num=?", stuNum).Delete(&Stu{})
	if err != nil || affected <= 0 {
		HandleResponse(c, http.StatusBadRequest, "删除失败", "error")
		return
	}
	HandleResponse(c, http.StatusOK, "删除成功", "OK")
	fmt.Println(affected) // 打印结果
}

// XormUpdateData 修改操作
func XormUpdateData(c *gin.Context) {
	var s Stu
	err := c.Bind(&s)
	if err != nil {
		HandleResponse(c, http.StatusBadRequest, "参数错误", "error")
		return
	}
	// 1、先查询
	var stus []Stu
	err = x.Where("stu_num=?", s.StuNum).Find(&stus)
	if err != nil || len(stus) <= 0 {
		HandleResponse(c, http.StatusBadRequest, "数据不存在", "error")
		return
	}
	// 2、再修改
	affected, err := x.Where("stu_num=?", s.StuNum).Update(&Stu{Name: s.Name, Age: s.Age})
	if err != nil || affected <= 0 {
		HandleResponse(c, http.StatusBadRequest, "修改失败", "error")
		return
	}
	HandleResponse(c, http.StatusOK, "修改成功", "OK")
	fmt.Println(affected) // 打印结果
}

// XormGetMulData 查询操作(多条记录)
func XormGetMulData[T any](c *gin.Context) {
	name := c.Query("name")
	var typs []T
	err := x.Where("name=?", name).And("age>20").Limit(10, 0).Asc("age").Find(&typs)
	if err != nil {
		HandleResponse(c, http.StatusBadRequest, "查询错误", "error")
		return
	}
	HandleResponse(c, http.StatusOK, "查询成功", typs)
}

func (orm *OrmValue[T]) QueryData(c *gin.Context) ([]T, error) {
	stuNum := c.Query(orm.Key)
	var types []T
	var sb strings.Builder
	sb.WriteString(orm.Key + "=?")

	types, err := orm.FindVal(sb.String(), stuNum)

	return types, err
}

func (orm *OrmValue[T]) FindVal(query any, args ...any) ([]T, error) {
	var types []T
	err := x.Where(query, args...).Find(&types)

	return types, err
}

// GetData 查询操作(单条记录)
func (orm *OrmValue[T]) GetData(c *gin.Context) {
	typs, err := orm.QueryData(c)
	if err != nil {
		HandleResponse(c, http.StatusBadRequest, "查询错误", "error")
		return
	}
	HandleResponse(c, http.StatusOK, "查询成功", typs)
}

func HandleResponse(c *gin.Context, code int, message string, data interface{}) {
	xormResponse.Code = code
	xormResponse.Message = message
	xormResponse.Data = data
	c.JSON(http.StatusOK, xormResponse)
}
