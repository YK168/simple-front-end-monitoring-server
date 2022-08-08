package service

import (
	"errors"
	"log"
	"net/http"
	"simple_front_end_monitoring_server/model"
	"simple_front_end_monitoring_server/utils"

	"gorm.io/gorm"
)

type UserService struct {
	Number string `form:"number" json:"number" binding:"required,min=11,max=11"`
	Passwd string `form:"passwd" json:"passwd" binding:"required,min=6,max=25"`
}

func (service *UserService) Register() utils.Response {
	var user model.User
	var count int64
	model.DB.Model(&model.User{}).Where("number = ?", service.Number).First(&user).Count(&count)
	if count == 1 {
		return utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "用户已存在",
		}
	}
	user.Number = service.Number
	// 密码加密
	if err := user.SetPasswd(service.Passwd); err != nil {
		return utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "密码加密失败",
			Error:  err.Error(),
		}
	}
	// 创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		return utils.Response{
			Status: http.StatusInternalServerError,
			Msg:    "创建用户失败",
			Error:  err.Error(),
		}
	}
	return utils.Response{
		Status: http.StatusOK,
		Msg:    "用户注册成功",
	}
}

func (service *UserService) Login() utils.Response {
	var user model.User
	// 查找用户是否存在
	err := model.DB.Where("number = ?", service.Number).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.Response{
				Status: http.StatusBadRequest,
				Msg:    "用户不存在",
			}
		}
		return utils.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库操作失败",
			Error:  err.Error(),
		}
	}
	// 密码校验
	if !user.CheckPasswd(service.Passwd) {
		return utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "密码错误",
		}
	}
	// 颁发Token
	token, err := utils.GenerateToken(user.ID, service.Number, service.Passwd)
	if err != nil {
		return utils.Response{
			Status: http.StatusInternalServerError,
			Msg:    "生成Token错误",
			Error:  err.Error(),
		}
	}
	return utils.Response{
		Status: http.StatusOK,
		Msg:    "登录成功",
		Data: utils.TokenData{
			User:  utils.BuildUser(user),
			Token: token,
		},
	}
}

type ProjectService struct {
	Number     string
	Title      string `form:"title" json:"title"`
	ProjectKey string `form:"project_key" json:"project_key"`
}

func (service *ProjectService) Create() utils.Response {
	// 查询project key是否已存在
	var count int64
	model.DB.Model(&model.Item{}).Where("project_key = ?", service.ProjectKey).Count(&count)
	if count >= 1 {
		return utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "项目名重复",
		}
	}

	// 添加item表记录
	var item model.Item = model.Item{
		Title:      service.Title,
		ProjectKey: service.ProjectKey,
		Number:     service.Number,
	}
	err := model.DB.Create(&item).Error
	if err != nil {
		log.Printf("数据库添加用户%s的%s项目记录失败\n", service.Number, service.Title)
		log.Println("err:", err)
		return utils.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库添加Item记录失败",
		}
	}
	return utils.Response{
		Status: http.StatusOK,
		Msg:    "项目创建成功",
		Data:   service.ProjectKey,
	}
}

func (service *ProjectService) Delete() utils.Response {
	err := model.DB.Model(&model.Item{}).Where("project_key = ?", service.ProjectKey).
		Delete(&model.Item{}).Error
	if err != nil {
		return utils.Response{
			Status: http.StatusInternalServerError,
			Msg:    "删除项目失败",
			Error:  err.Error(),
		}
	}
	return utils.Response{
		Status: http.StatusOK,
		Msg:    "删除项目成功",
	}
}
