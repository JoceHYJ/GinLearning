package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"unicode/utf8"
)

type UserInfo struct {
	Id   string `validate:"uuid" json:"id"`
	Name string `validate:"checkName" json:"name"`
	Age  uint   `validate:"min=0,max=130" json:"age"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("checkName", checkNameFunc)
}

// 自定义参数校验规则
func checkNameFunc(f validator.FieldLevel) bool {
	count := utf8.RuneCountInString(f.Field().String())
	if count >= 2 && count <= 12 {
		return true
	}
	return false
}

func main() {
	r := gin.Default()
	var user UserInfo

	u1 := uuid.New()
	fmt.Println("uuid is :", u1.String())

	r.POST("/validate", func(c *gin.Context) {
		err := c.Bind(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, "请求参数错误")
			return
		}

		// 校验
		err = validate.Struct(user)
		if err != nil {
			for _, e := range err.(validator.ValidationErrors) {
				fmt.Println("错误的字段:", e.Field())
				fmt.Println("错误的值:", e.Value())
				fmt.Println("错误的 tag:", e.Tag())
			}
			c.JSON(http.StatusBadRequest, "数据校验失败")
			return
		}
		c.JSON(http.StatusOK, "数据校验成功")
	})
	r.Run(":8080")
}
