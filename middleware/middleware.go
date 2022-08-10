package middleware

import (
	"net/http"
	"simple_front_end_monitoring_server/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT(c *gin.Context) {
	status := http.StatusOK
	msg := ""
	token := c.GetHeader("x-token")
	if token == "" {
		status = http.StatusUnauthorized
		msg = "未从header中获取到token，请检查请求头是否存在x-token"
	} else {
		claims, err := utils.ParseToken(token)
		if err != nil {
			status = http.StatusForbidden
			msg = "解析token失败"
		} else if time.Now().Unix() > claims.ExpiresAt {
			status = http.StatusForbidden
			msg = "token已过期"
		}
	}
	if status != http.StatusOK {
		c.JSON(status, gin.H{
			"status": status,
			"msg":    msg,
			"data":   nil,
		})
		c.Abort()
		return
	}
	c.Next()
}
