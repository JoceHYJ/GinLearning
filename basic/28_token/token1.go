package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

type HmacUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

type MyClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

var jwtKey = []byte("a_secret_key") // 证书签名密钥

func main() {
	r := gin.Default()
	// Token 分发
	r.POST("getTokenHMAC", func(c *gin.Context) {
		var u HmacUser
		c.Bind(&u)
		token, err := hmacReleaseToken(u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "分发成功",
			"data": token,
		})
	})
	// Token 认证
	r.POST("/checkTokenHMAC", hmacAuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, "Token认证成功")
	})
	r.Run(":8080")
}

// hmacAuthMiddleware 中间件，检查 Token 的有效性
func hmacAuthMiddleware() gin.HandlerFunc {
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
		token, claims, err := hmacParseToken(tokenString)
		fmt.Println(err)
		if err != nil || !token.Valid { // 解析错误或证书无效
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "无效的 Token"})
			c.Abort()
			return
		}
		var u HmacUser
		c.Bind(&u)
		if u.Id != claims.UserId {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "用户不存在"})
			c.Abort()
			return
		}
		// 验证通过，将用户信息保存到请求中
		c.Next()
	}
}

// hmacParseToken 解析 Token
func hmacParseToken(tokenString string) (*jwt.Token, *MyClaims, error) {
	claims := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}

// hmacReleaseToken 分发 Token
func hmacReleaseToken(u HmacUser) (string, error) {
	// 有效时间为当前日期加七天
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	// 创建一个Claims
	claims := &MyClaims{
		UserId: u.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),     // 签发时间
			Subject:   "user token",          // 主题
			Issuer:    "tomato",              // 签发者
		},
	}
	// 生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
