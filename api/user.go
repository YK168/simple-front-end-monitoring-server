package api

import (
	"log"
	"net/http"
	"simple_front_end_monitoring_server/service"
	"simple_front_end_monitoring_server/utils"

	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	content := utils.GetQueryContent(c)
	log.Println(content)
	var userRegister service.UserService
	if err := c.ShouldBind(&userRegister); err == nil {
		res := userRegister.Register()
		c.JSON(res.Status, res)
	} else {
		log.Println("用户注册失败 json解码失败", err)
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "用户注册失败 json解码失败",
			Error:  err.Error(),
		})
	}
}

func UserLogin(c *gin.Context) {
	var userLogin service.UserService
	if err := c.ShouldBind(&userLogin); err == nil {
		res := userLogin.Login()
		c.JSON(res.Status, res)
	} else {
		log.Println("用户登录失败 json解码失败", err)
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "用户登录失败 json解码失败",
			Error:  err.Error(),
		})
	}
}
