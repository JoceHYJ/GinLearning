package middleware

import (
	"GinLearning/gin_application/common"
	"GinLearning/gin_application/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthMiddleware token 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := "tomato"
		// 获取 authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, auth+":") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "前缀错误"})
			c.Abort()
			return
		}
		index := strings.Index(tokenString, auth+":") // 找到 Token 前缀对应的位置
		// 真实 Token
		tokenString = tokenString[index+len(auth)+1:] // 截取真实 Token (开始位置为: 索引开始的位置 + 关键字符长度+ 1(:的长度))
		// 解析 Token
		token, claims, err := common.ParseToken(tokenString)
		fmt.Println(err)
		if err != nil || !token.Valid { // 解析错误或证书无效
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "无效的 Token"})
			c.Abort()
			return
		}
		userID := claims.UserId
		var user model.User
		common.DB.First(&user, userID)
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "用户不存在"})
			c.Abort()
			return
		}
		// 验证通过，将用户信息保存到请求中
		c.Set("user", user)
		c.Next()
	}
}
