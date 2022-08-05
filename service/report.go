package service

import (
	"net/http"
	"simple_front_end_monitoring_server/model"
	"simple_front_end_monitoring_server/utils"
)

type JsErrorService struct {
	// 项目名称
	Title string `form:"title" json:"title" binding:"required"`
	// 错误信息
	Message string `form:"message" json:"message" binding:"required"`
	// 报错时路由地址
	URL string `form:"url" json:"url" binding:"required"`
	// 报错代码行数
	Position string `form:"position" json:"position" binding:"required"`
	// 报错文件
	FileName string `form:"file_name" json:"file_name" binding:"required"`
	// JsError or PromiseError
	ErrType string `form:"err_type" json:"err_type" binding:"required"`
	// 报错时间
	TimeStamp int64 `form:"time_stamp" json:"time_stamp" binding:"required"`
	// 根据Cookie来区分不同页面？
	Cookie     string `form:"cookie" json:"cookie" binding:"required"`
	ProjectKey string `form:"project_key" json:"project_key" binding:"required"`
}

type ApiErrorService struct {
	// 项目名称
	Title string `form:"title" json:"title" binding:"required"`
	// 错误信息
	Message string `form:"message" json:"message" binding:"required"`
	// 报错时路由地址
	URL string `form:"url" json:"url" binding:"required"`
	// 报错代码行数
	Position string `form:"position" json:"position" binding:"required"`
	// 报错文件
	FileName string `form:"file_name" json:"file_name" binding:"required"`
	// JsError or PromiseError
	ErrType string `form:"err_type" json:"err_type" binding:"required"`
	// 报错时间
	TimeStamp int64 `form:"time_stamp" json:"time_stamp" binding:"required"`
	// 根据Cookie来区分不同页面？
	Cookie     string `form:"cookie" json:"cookie" binding:"required"`
	ProjectKey string `form:"project_key" json:"project_key" binding:"required"`
}

func (service *JsErrorService) Report() utils.Response {
	jsErr := model.JSError{
		Title:      service.Title,
		Message:    service.Message,
		URL:        service.URL,
		Position:   service.Position,
		FileName:   service.FileName,
		ErrType:    service.ErrType,
		TimeStamp:  service.TimeStamp,
		Cookie:     service.Cookie,
		ProjectKey: service.ProjectKey,
	}
	err := model.DB.Create(&jsErr).Error
	if err != nil {
		return utils.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库操作失败，添加Js Error监控记录失败",
		}
	}
	return utils.Response{
		Status: http.StatusOK,
		Msg:    "Js Error监控记录添加成功",
	}
}

func (service *ApiErrorService) Report() utils.Response {
	apiErr := model.APIError{
		Title:      service.Title,
		Message:    service.Message,
		URL:        service.URL,
		Position:   service.Position,
		FileName:   service.FileName,
		ErrType:    service.ErrType,
		TimeStamp:  service.TimeStamp,
		Cookie:     service.Cookie,
		ProjectKey: service.ProjectKey,
	}
	err := model.DB.Create(&apiErr).Error
	if err != nil {
		return utils.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库操作失败，添加Api Error监控记录失败",
		}
	}
	return utils.Response{
		Status: http.StatusOK,
		Msg:    "Api Error监控记录添加成功",
	}
}
