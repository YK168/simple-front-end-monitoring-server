package api

import (
	"log"
	"net/http"
	"simple_front_end_monitoring_server/service"
	"simple_front_end_monitoring_server/utils"

	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
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

// 根据用户手机号码和项目名称创建对应项目key
func ProjectCreate(c *gin.Context) {
	// 获取token
	token := c.Request.Header.Get("x-token")
	// 根据token获取用户信息
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Status: http.StatusInternalServerError,
			Msg:    "jwt解析token失败",
			Error:  err.Error(),
		})
		return
	}

	// 创建项目逻辑
	var project = service.ProjectService{Number: claims.Number}
	// 这里的project只绑定了Title
	if err := c.ShouldBind(&project); err == nil {
		// TODO: 验证用户是否存在

		log.Printf("正在生成用户%s的%s项目KEY...", project.Number, project.Title)
		key := utils.MD5(project.Number + project.Title)
		project.ProjectKey = key
		log.Printf("用户%s的%s项目KEY: %s\n", project.Number, project.Title, project.ProjectKey)

		res := project.Create()
		c.JSON(res.Status, res)
	} else {
		log.Println("项目创建失败，解析json参数失败，err:", err)
		log.Printf("number = %s\n", project.Number)
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "项目创建失败，解析json参数失败",
			Error:  err.Error(),
		})
	}
}

func ProjectDelete(c *gin.Context) {
	var project service.ProjectService
	if err := c.ShouldBind(&project); err == nil {
		res := project.Delete()
		c.JSON(res.Status, res)
	} else {
		log.Println("项目删除失败，解析json参数失败，err:", err)
		log.Printf("project_key = %s\n", project.ProjectKey)
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "项目删除失败，解析json参数失败",
			Error:  err.Error(),
		})
	}
}
