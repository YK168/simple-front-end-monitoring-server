package utils

import (
	"bytes"
	"io/ioutil"
	"log"
	"math/rand"
	"simple_front_end_monitoring_server/model"
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

func GetBorder[T int | float32](s []T) (int, int) {
	// 截取s切片头尾空元素
	// r = len(s)，防止全0数组时，切片错误
	l, r := 0, len(s)
	for l < r && (s[l] == 0 || s[r-1] == 0) {
		if s[l] == 0 {
			l++
		}
		if r != l && s[r-1] == 0 {
			r--
		}
	}
	return l, r
}

func Get2Border(s1, s2 []int) (int, int) {
	if len(s1) != len(s2) {
		log.Println("Get2Border: 长度不一致")
		return 0, 0
	}
	l, r := 0, len(s1)
	for l < r && ((s1[l] == 0 && s2[l] == 0) || (s1[r-1] == 0 && s2[r-1] == 0)) {
		if s1[l] == 0 && s2[l] == 0 {
			l++
		}
		if r != l && s1[r-1] == 0 && s2[r-1] == 0 {
			r--
		}
	}
	return l, r
}

// 生成一些测试数据
func GenerateTestData(number int, projectKey string, startT, endT int64) {
	msg := "Uncaught ReferenceError: t is not defined"
	filename := "webpack-internal:///./node_modules/test.js"
	position := "378:13"
	url := "http://localhost:8000/"
	title := "测试数据"
	cookies := []string{"userCookie1", "userCookie2", "userCookie3", ""}
	gap := endT - startT
	// 写1秒一条
	for i := 0; i < number; i++ {
		log.Printf("正在插入第%d条数据\n", i)
		ch := byte('a' + rand.Intn(26))
		tUrl := url + string(ch)
		jsErr := model.JSError{
			Title:      title,
			ProjectKey: projectKey,
			Message:    msg,
			URL:        tUrl,
			Position:   position,
			FileName:   filename,
			TimeStamp:  startT + rand.Int63n(gap),
			ErrType:    []string{"jsError", "promiseError"}[rand.Intn(2)],
			Cookie:     cookies[rand.Intn(len(cookies))],
		}
		model.DB.Create(&jsErr)

		model.DB.Model(&model.APIError{}).Create(&model.APIError{
			Title:      title,
			ProjectKey: projectKey,
			URL:        tUrl,
			TimeStamp:  startT + rand.Int63n(gap),
			Cookie:     cookies[rand.Intn(len(cookies))],
			Params:     "POST请求参数",
			Response:   "响应内容",
			Pathname: []string{
				"https://www.baidu.com/colors/test",
				"https://www.taobao.cloud/colors/test",
				"https://www.jd.cloud/colors/test",
				"https://www.tengxun.cloud/colors/test",
				"https://www.xiaomi.cloud/colors/test",
			}[rand.Intn(5)],
			Status:    []string{"200", "400", "500"}[rand.Intn(3)],
			Duration:  rand.Intn(50) + 1,
			EventType: []string{"error", "load"}[rand.Intn(2)],
			Kind:      "stability",
			ReqType:   "xhr",
		})
		model.DB.Model(&model.SourceError{}).Create(&model.SourceError{
			Title:      title,
			URL:        tUrl,
			Cookie:     cookies[rand.Intn(len(cookies))],
			FileName:   filename,
			TimeStamp:  startT + rand.Int63n(gap),
			ProjectKey: projectKey,
			ErrType:    "resourceError",
			TagName:    "IMG",
		})
		model.DB.Model(&model.Performance{}).Create(&model.Performance{
			Title:        title,
			TimeStamp:    startT + rand.Int63n(gap),
			ProjectKey:   projectKey,
			URL:          tUrl,
			Cookie:       cookies[rand.Intn(len(cookies))],
			AnalysisTime: rand.Float32() + float32(rand.Intn(10)),
			AppcacheTime: rand.Float32() + float32(rand.Intn(10)),
			BlankTime:    rand.Float32() + float32(rand.Intn(10)),
			DnsTime:      rand.Float32() + float32(rand.Intn(10)),
			DomReadyTime: rand.Float32() + float32(rand.Intn(10)),
			LoadPageTime: rand.Float32() + float32(rand.Intn(10)),
			RedirectTime: rand.Float32() + float32(rand.Intn(10)),
			ReqTime:      rand.Float32() + float32(rand.Intn(10)),
			TcpTime:      rand.Float32() + float32(rand.Intn(10)),
			TtfbTime:     rand.Float32() + float32(rand.Intn(10)),
			UnloadTim:    rand.Float32() + float32(rand.Intn(10)),
		})
		model.DB.Model(&model.Access{}).Create(&model.Access{
			Title:      title,
			URL:        tUrl,
			Cookie:     cookies[rand.Intn(len(cookies))],
			TimeStamp:  startT + rand.Int63n(gap),
			ProjectKey: projectKey,
			ErrType:    "pv",
		})
	}
}
