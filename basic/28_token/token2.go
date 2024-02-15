package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
	"time"
)

type RsaUser struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

type RsaClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

var (
	resPrivateKey []byte
	resPublicKey  []byte
	err2_1        error
	err2_2        error
)

func init() {
	// 读取 RSA 密钥文件
	resPrivateKey, err2_1 = os.ReadFile("./basic/28_token/private.pem")
	resPublicKey, err2_2 = os.ReadFile("./basic/28_token/public.pem")
	if err2_1 != nil || err2_2 != nil {
		panic(fmt.Sprintf("读取密钥文件失败, err1: %v, err2: %v", err2_1, err2_2))
		return
	}
}

// 通过 RSA 签名实现 Token
// RSA 密钥生成工具: http://www.metools.info/code/c80.html
func main() {
	r := gin.Default()
	// Token 分发
	r.POST("/getTokenRSA", func(c *gin.Context) {
		u := RsaUser{}
		err := c.Bind(&u)
		if err != nil {
			c.JSON(http.StatusBadRequest, "参数错误")
			return
		}
		token, err := rsaReleaseToken(u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Token 生成失败")
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "Token 分发成功",
			"data": token,
		})
	})
	// Token 验证
	r.POST("/checkTokenRSA", rsaAuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "Token 验证通过",
		})
	})
	r.Run(":8080")
}

// rsaAuthMiddleware 中间件，检查 Token 的有效性
func rsaAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := "tomato"
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, auth+":") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "无效的 Token"})
			c.Abort()
			return
		}
		index := strings.Index(tokenString, auth+":")
		tokenString = tokenString[index+len(auth)+1:]
		claims, err := rsaParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "证书无效"})
			c.Abort()
			return
		}
		claimsValue := claims.(jwt.MapClaims) // 断言为 jwt.MapClaims 类型
		if claimsValue["user_id"] == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "用户不存在"})
			c.Abort()
			return
		}
		u := RsaUser{}
		c.Bind(&u)
		id := claimsValue["user_id"].(string)
		if u.Id != id {
			c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "用户信息不匹配"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// rsaParseToken 解析 Token
func rsaParseToken(tokenString string) (any, error) {
	pem, err := jwt.ParseRSAPublicKeyFromPEM(resPublicKey)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("解析的方法错误")
		} // 获取签名方法
		return pem, err
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// rsaReleaseToken 分发 Token
func rsaReleaseToken(u RsaUser) (any, error) {
	tokenGen, err := rsaJwtTokenGen(u.Id)
	return tokenGen, err
}

// rsaJwtTokenGen 生成 Token
func rsaJwtTokenGen(id string) (any, error) {
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(resPrivateKey)
	if err != nil {
		return nil, err
	}
	claims := RsaClaims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), // 过期时间
			Issuer:    "tomato",                                  // 发布者
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &claims)
	signedString, err := token.SignedString(privateKey)
	return signedString, err
}
