package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type ValUser struct {
	Name    string       `validate:"required" json:"name"`
	Age     uint         `validate:"gte=0,lte=130" json:"age"`
	Email   string       `validate:"required" json:"email"`
	Address []ValAddress `validate:"dive" json:"address"` // 结构体嵌套：dive
}

type ValAddress struct {
	Province string `validate:"required" json:"province"`
	City     string `validate:"required" json:"city"`
	Phone    string `validate:"numeric,len=11" json:"phone"` // 不能有空格
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func main() {
	r := gin.Default()
	var user ValUser
	r.POST("/validate", func(c *gin.Context) {
		//testData(c)
		err := c.Bind(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, "参数错误，绑定失败！")
			return
		}
		// 执行校验
		if validateUser(user) {
			c.JSON(http.StatusOK, "数据校验成功！")
			return
		}
		c.JSON(http.StatusBadRequest, "数据校验失败！")
	})
	r.Run(":8080")
}

func testData(c *gin.Context) {
	address := ValAddress{
		Province: "浙江省",
		City:     "杭州市",
		Phone:    "13366663636",
	}
	user := ValUser{
		Name:    "tomato",
		Age:     20,
		Email:   "gin@163.com",
		Address: []ValAddress{address, address},
	}
	c.JSON(http.StatusOK, user)
}

func validateUser(u ValUser) bool {
	err := validate.Struct(u)
	for _, e := range err.(validator.ValidationErrors) {
		fmt.Println("错误的字段:", e.Field())
		fmt.Println("错误的值", e.Value())
	}
	if err != nil {
		return false
	}
	return true
}
