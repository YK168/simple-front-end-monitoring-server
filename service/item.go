package service

import (
	"log"
	"net/http"
	"simple_front_end_monitoring_server/model"
	"simple_front_end_monitoring_server/utils"
)

type ProjectService struct {
	Number     string `form:"number" json:"number"`
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
	if service.ProjectKey == "" {
		return utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "删除项目失败，projectKey为空",
		}
	}
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

func (service *ProjectService) Search(number string) utils.Response {
	if number == "" {
		return utils.Response{
			Status: http.StatusBadRequest,
			Msg:    "用户号码为空",
		}
	}
	data := []ProjectService{}
	err := model.DB.Model(&model.Item{}).Where("number = ?", number).Find(&data).Error
	if err != nil {
		return utils.Response{
			Status: http.StatusInternalServerError,
			Msg:    "获取项目失败",
			Error:  err.Error(),
		}
	}
	return utils.Response{
		Status: http.StatusOK,
		Msg:    "获取项目成功",
		Data:   data,
	}
}
