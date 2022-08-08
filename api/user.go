package api

import (
	"log"
	"net/http"
	"simple_front_end_monitoring_server/model"
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
	type Info struct {
		Number      string `form:"number" json:"number" binding:"required,min=11,max=11"`
		ProjectName string `form:"project_name" json:"project_name" binding:"required"`
	}
	var info Info
	if err := c.ShouldBind(&info); err == nil {
		// TODO: 查询用户是否存在
		log.Printf("正在生成用户%s的%s项目KEY...", info.Number, info.ProjectName)
		key := utils.MD5(info.Number + info.ProjectName)
		log.Printf("用户%s的%s项目KEY: %s\n", info.Number, info.ProjectName, key)
		// 添加item表记录
		var item model.Item = model.Item{
			Title:      info.ProjectName,
			ProjectKey: key,
			Number:     info.Number,
		}
		err := model.DB.Create(&item).Error
		if err != nil {
			log.Printf("数据库添加用户%s的%s项目记录失败\n", info.Number, info.ProjectName)
			log.Println("err:", err)
			c.JSON(http.StatusOK, utils.Response{
				Status: http.StatusOK,
				Msg:    "数据库添加Item记录失败",
				Data:   key,
			})
			return
		}
		c.JSON(http.StatusOK, utils.Response{
			Status: http.StatusOK,
			Msg:    "项目创建成功",
			Data:   key,
		})
	} else {
		log.Println("项目创建失败，解析json参数失败")
		log.Printf("number = <%s>, project name = <%s>\n", info.Number, info.ProjectName)
		c.JSON(http.StatusBadRequest, utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "项目创建失败，解析json参数失败",
		})
	}
}
