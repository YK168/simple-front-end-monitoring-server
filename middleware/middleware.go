package middleware

import (
	"net/http"
	"simple_front_end_monitoring_server/utils"
	"strconv"
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
		c.JSON(status, utils.Response{
			Status: status,
			Msg:    msg,
		})
		c.Abort()
		return
	}
	c.Next()
}

// 解析资源GET请求时参数和格式是否正确
func ParseURL(c *gin.Context) {
	projectKey := c.Query("projectKey")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	if projectKey == "" || startTime == "" || endTime == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "projectKy或startTime或endTime参数为空",
		})
		c.Abort()
		return
	}
	_, err1 := strconv.ParseInt(startTime, 10, 64)
	_, err2 := strconv.ParseInt(endTime, 10, 64)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "startTime或endTime参数格式不正确，非时间戳格式",
		})
		c.Abort()
		return
	}
	c.Next()
}
