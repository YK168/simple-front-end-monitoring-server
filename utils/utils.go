package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT密钥
var JWTsecret = []byte("simple-front-end-monitoring-server")

type Claims struct {
	ID     uint   `json:"id"`
	Number string `json:"number"`
	Passwd string `json:"passwd"`
	jwt.StandardClaims
}

// 生成、签发token
func GenerateToken(id uint, number, passwd string) (string, error) {
	claims := Claims{
		ID:     id,
		Number: number,
		Passwd: passwd,
		StandardClaims: jwt.StandardClaims{
			// 24小时候token过期，需要重新获取
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			// 签发机构
			Issuer: "simple-front-end-monitoring-server",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(JWTsecret)
	return token, err
}

// 验证token
func ParseToken(token string) (*Claims, error) {
	tokenClaim, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return JWTsecret, nil
	})
	if tokenClaim != nil {
		if claims, ok := tokenClaim.Claims.(*Claims); ok && tokenClaim.Valid {
			return claims, nil
		}
	}
	return nil, err
}
