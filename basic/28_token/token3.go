package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

type EcdsaUser struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

type EcdsaClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

var (
	err3          error
	ecdPrivateKey *ecdsa.PrivateKey
	ecdPublicKey  *ecdsa.PublicKey
)

func init() {
	ecdPrivateKey, ecdPublicKey, err3 = getEcdsaKey(2)
	if err3 != nil {
		panic(err3)
		return
	}
}

// 通过椭圆曲线数字签名（ECDSA）算法生成椭圆曲线密钥对, 对给定的消息进行签名和验证
func main() {
	r := gin.Default()
	// 生成Token
	r.POST("/getTokenECDSA", func(c *gin.Context) {
		u := EcdsaUser{}
		err := c.Bind(&u)
		if err != nil {
			c.JSON(http.StatusBadRequest, "参数错误")
			return
		}
		token, err := ecdsaReleaseToken(u)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "生成Token失败")
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "授权成功",
			"data": token,
		})
	})
	//  验证Token
	r.POST("/checkTokenECDSA", ecdsaTokenMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "验证Token成功",
		})
	})
	r.Run(":8080")
}

// ecdsaTokenMiddleware 中间件验证Token
func ecdsaTokenMiddleware(c *gin.Context) {
	auth := "tomato"
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, auth+":") {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "无效的 Token"})
		c.Abort()
		return
	}
	index := strings.Index(tokenString, auth+":")
	tokenString = tokenString[index+len(auth)+1:]
	claims, err := ecdsaParseToken(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}
	claimsValue := claims.(jwt.MapClaims)
	if claimsValue["user_id"] == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "id 不存在")
		return
	}
	u := EcdsaUser{}
	c.Bind(&u)
	if u.ID != claimsValue["user_id"] {
		c.JSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "msg": "用户不存在"})
		c.Abort()
		return
	}
	c.Next()
}

// ecdsaParseToken 解析Token
func ecdsaParseToken(tokenString string) (any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v", token.Method)
		}
		return ecdPublicKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// ecdsaReleaseToken 分发 Token
func ecdsaReleaseToken(u EcdsaUser) (any, error) {
	claims := &EcdsaClaims{
		UserId: u.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(), // 设置过期时间
			Issuer:    "tomato",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signedString, err := token.SignedString(ecdPrivateKey) // 使用私钥对 Token 进行签名
	return signedString, err
}

// getEcdsaKey 生成椭圆曲线密钥对
func getEcdsaKey(keyType int) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	var err error
	var prk *ecdsa.PrivateKey
	var pub *ecdsa.PublicKey
	var curve elliptic.Curve // 椭圆曲线
	switch keyType {
	case 1:
		curve = elliptic.P224()
	case 2:
		curve = elliptic.P256()
	case 3:
		curve = elliptic.P384()
	case 4:
		curve = elliptic.P521()
	default:
		errors.New("输入的签名 key 类型错误！key 取值：\n 1:椭圆曲线224\n 2:椭圆曲线256\n 3:椭圆曲线384\n 4:椭圆曲线521\n")
		return nil, nil, err
	}
	prk, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	pub = &prk.PublicKey
	return prk, pub, nil
}
