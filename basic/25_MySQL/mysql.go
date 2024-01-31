package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

// Client 提交的数据需要在 Server 永久保存
// mysql 驱动: go get github.com/go-sql-driver/mysql

var sqlDb *sql.DB
var sqlResponse SqlResponse

func init() {
	// 1、打开数据库
	sqlStr := "root:010729@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=true&loc=Local"
	var err error
	sqlDb, err = sql.Open("mysql", sqlStr)
	if err != nil {
		fmt.Println("数据库打开出现了问题:", err)
		return
	}
	// 2、测试与数据库建立的连接(校验连接是否正确)
	err = sqlDb.Ping()
	if err != nil {
		fmt.Println("数据库连接出现了问题:", err)
		return
	}
}

// SqlUser Client 提交的数据
type SqlUser struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

// SqlResponse 应答体(响应 Client 的请求)
type SqlResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func main() {
	r := gin.Default()
	// 数据库的 CRUD ---> Gin 的 POST GET PUT DELETE 方法
	r.POST("sql/insert", insertData)
	r.GET("sql/get", getData)
	r.GET("sql/mulget", getMulData)
	r.PUT("sql/update", updateData)
	r.DELETE("sql/delete", deleteData)
	r.Run(":8080")
}

// deleteData 删除操作
func deleteData(c *gin.Context) {
	name := c.Query("name")
	var count int
	// 1、先查询
	sqlStr := "select count(*) from user where name=?"
	err := sqlDb.QueryRow(sqlStr, name).Scan(&count)
	if count <= 0 || err != nil {
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "删除数据不存在"
		sqlResponse.Data = "error"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	// 2、再删除
	delStr := "delete from user where name=?"
	ret, err := sqlDb.Exec(delStr, name)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "删除失败"
		sqlResponse.Data = "error"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	sqlResponse.Code = http.StatusOK
	sqlResponse.Message = "删除成功"
	sqlResponse.Data = "OK"
	fmt.Println(ret.LastInsertId()) // 打印结果
}

// updateData 修改操作
func updateData(c *gin.Context) {
	var u SqlUser
	err := c.BindJSON(&u)
	if err != nil {
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "参数错误"
		sqlResponse.Data = "error"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	sqlStr := "update user set age=?, address=? where name=?"
	ret, err := sqlDb.Exec(sqlStr, u.Age, u.Address, u.Name)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "更新失败"
		sqlResponse.Data = "error"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	sqlResponse.Code = http.StatusOK
	sqlResponse.Message = "更新成功"
	sqlResponse.Data = "OK"
	c.JSON(http.StatusOK, sqlResponse)
	fmt.Println(ret.LastInsertId()) // 打印结果
}

// getMulData 查询操作(多条记录)
func getMulData(c *gin.Context) {
	address := c.Query("address")
	sqlStr := "select name, age from user where address=?"
	rows, err := sqlDb.Query(sqlStr, address)
	if err != nil {
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "查询错误"
		sqlResponse.Data = "error"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	defer rows.Close()
	// 创建切片并通过 for 循环把查询到的数据写入切片中
	resUser := make([]SqlUser, 0)
	for rows.Next() {
		var userTemp SqlUser
		rows.Scan(&userTemp.Name, &userTemp.Age)
		userTemp.Address = address
		resUser = append(resUser, userTemp)
	}
	// 将数据返回 Client
	sqlResponse.Code = http.StatusOK
	sqlResponse.Message = "读取成功"
	sqlResponse.Data = resUser
	c.JSON(http.StatusOK, sqlResponse)
}

// getData 查询操作(单条记录)
func getData(c *gin.Context) {
	name := c.Query("name")
	sqlStr := "select age, address from user where name = ?"
	var u SqlUser
	err := sqlDb.QueryRow(sqlStr, name).Scan(&u.Age, &u.Address)
	if err != nil {
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "查询错误"
		sqlResponse.Data = "error"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	u.Name = name
	sqlResponse.Code = http.StatusOK
	sqlResponse.Message = "读取成功"
	sqlResponse.Data = u
	c.JSON(http.StatusOK, sqlResponse)
}

// insertData 插入操作
func insertData(c *gin.Context) {
	var u SqlUser
	err := c.BindJSON(&u)
	if err != nil {
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "参数错误"
		sqlResponse.Data = "error"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	sqlStr := "insert into user(name, age, address)values (?,?,?)"
	ret, err := sqlDb.Exec(sqlStr, u.Name, u.Age, u.Address)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		sqlResponse.Code = http.StatusBadRequest
		sqlResponse.Message = "写入失败"
		sqlResponse.Data = "error"
		c.JSON(http.StatusOK, sqlResponse)
		return
	}
	sqlResponse.Code = http.StatusOK
	sqlResponse.Message = "写入成功"
	sqlResponse.Data = "OK"
	c.JSON(http.StatusOK, sqlResponse)
	fmt.Println(ret.LastInsertId()) // 打印结果
}
