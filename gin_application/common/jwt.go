package common

import (
	"GinLearning/gin_application/model"
	"github.com/golang-jwt/jwt"
	"time"
)

var jwtKey = []byte("my_secret_key") // 证书签名密钥

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

// ReleaseToken 分发 Token
func ReleaseToken(user model.User) (string, error) {
	// 有效时间为当前日期加七天
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	// 创建一个Claims
	claims := &Claims{
		UserId: user.ID,
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

// ParseToken 解析 Token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
