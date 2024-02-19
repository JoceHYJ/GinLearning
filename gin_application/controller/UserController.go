package controller

import (
	"GinLearning/gin_application/common"
	"GinLearning/gin_application/model"
	"GinLearning/gin_application/response"
	"GinLearning/gin_application/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

// Register 注册
func Register(c *gin.Context) {
	var requestUser model.User
	c.Bind(&requestUser)
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password

	// 数据验证
	if len(telephone) != 11 {
		// 422 Unprocessable Entity 无法处理的请求实体
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		fmt.Println(telephone, len(telephone))
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	// 如果名称没有传入，给一个 10 位的随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}
	// 判断手机号是否存在
	if isTelephoneExist(common.DB, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}
	// 创建用户
	// 返回密码的 hash 值 (对用户密码进行二次处理，防止数据库泄露后用户密码泄露)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "加密错误")
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hashPassword),
	}
	common.DB.Create(newUser)
	// 分发 token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "系统异常"})
	}
	response.Success(c, gin.H{"token": token}, "注册成功")
}

// Login 登录
func Login(c *gin.Context) {
	var requestUser model.User
	c.Bind(&requestUser)
	//name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	// 数据验证
	if len(telephone) != 11 {
		// 422 Unprocessable Entity 无法处理的请求实体
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		fmt.Println(telephone, len(telephone))
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	// 依据用户输入的手机号查询注册数据记录
	var user model.User
	common.DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": http.StatusUnprocessableEntity, "msg": "用户不存在"})
		return
	}
	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": "密码错误"})
	}
	// 分发 token
	token, err := common.ReleaseToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "系统异常"})
	}
	response.Success(c, gin.H{"token": token}, "登录成功")
}

// Info
func Info(c *gin.Context) {
	user, _ := c.Get("user")
	response.Success(c, gin.H{"user": response.ToUserDto(user.(model.User))}, "响应成功")
}

// isTelephoneExist 判断手机号是否存在
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
