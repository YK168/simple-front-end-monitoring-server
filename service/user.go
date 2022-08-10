package service

import (
	"errors"
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
