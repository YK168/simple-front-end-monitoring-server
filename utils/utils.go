package utils

import (
	"bytes"
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

func GetQueryContent(c *gin.Context) string {
	data, err := c.GetRawData()
	if err != nil {
		log.Println("获取请求body失败:", err.Error())
	}
	// 将读取出来的内容重新放回流中
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	return c.Request.Method + " " + c.Request.URL.String() + " " + string(data)
}

func GetBorder(s []int) (int, int) {
	// 截取s切片头尾空元素
	// r = len(s)，防止全0数组时，切片错误
	l, r := 0, len(s)
	for l < r-1 && (s[l] == 0 || s[r-1] == 0) {
		if s[l] == 0 {
			l++
		}
		if s[r-1] == 0 {
			r--
		}
	}
	return l, r
}
