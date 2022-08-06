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

type SourceErrorService struct {
	// 项目名称
	Title string `form:"title" json:"title" binding:"required"`
	// 报错时路由地址
	URL string `form:"url" json:"url" binding:"required"`
	// 报错文件
	FileName string `form:"file_name" json:"file_name" binding:"required"`
	// 报错资源标签
	TagName string `form:"tag_name" json:"tag_name" binding:"required"`
	// JsError or PromiseError
	ErrType string `form:"err_type" json:"err_type" binding:"required"`
	// 报错时间
	TimeStamp int64 `form:"time_stamp" json:"time_stamp" binding:"required"`
	// 根据Cookie来区分不同页面？
	Cookie     string `form:"cookie" json:"cookie" binding:"required"`
	ProjectKey string `form:"project_key" json:"project_key" binding:"required"`
}

type PerformanceService struct {
	// 项目名称
	Title        string  `form:"title" json:"title" binding:"required"`
	AnalysisTime float32 `form:"analysis_time" json:"analysis_time"`
	AppcacheTime float32 `form:"appcache_time" json:"appcache_time"`
	BlankTime    float32 `form:"blank_time" json:"blank_time"`
	DomReadyTime float32 `form:"dom_ready_time" json:"dom_ready_time"`
	LoadPageTime float32 `form:"load_page_time" json:"load_page_time"`
	RedirectTime float32 `form:"redirect_time" json:"redirect_time"`
	ReqTime      float32 `form:"req_time" json:"req_time"`
	TcpTime      float32 `form:"tcp_time" json:"tcp_time"`
	TtfbTime     float32 `form:"ttfb_time" json:"ttfb_time"`
	UnloadTim    float32 `form:"upload_time" json:"upload_time"`
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

func (service *SourceErrorService) Report() utils.Response {
	sourceErr := model.SourceError{
		Title:      service.Title,
		URL:        service.URL,
		TagName:    service.TagName,
		FileName:   service.FileName,
		ErrType:    service.ErrType,
		TimeStamp:  service.TimeStamp,
		Cookie:     service.Cookie,
		ProjectKey: service.ProjectKey,
	}
	err := model.DB.Create(&sourceErr).Error
	if err != nil {
		return utils.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库操作失败，添加Source Error监控记录失败",
		}
	}
	return utils.Response{
		Status: http.StatusOK,
		Msg:    "Source Error监控记录添加成功",
	}
}

func (service *PerformanceService) Report() utils.Response {
	performance := model.Performance{
		Title:        service.Title,
		AnalysisTime: service.AnalysisTime,
		AppcacheTime: service.AppcacheTime,
		BlankTime:    service.BlankTime,
		DomReadyTime: service.DomReadyTime,
		LoadPageTime: service.LoadPageTime,
		RedirectTime: service.RedirectTime,
		ReqTime:      service.ReqTime,
		TcpTime:      service.TcpTime,
		TtfbTime:     service.TtfbTime,
		UnloadTim:    service.UnloadTim,
		TimeStamp:    service.TimeStamp,
		Cookie:       service.Cookie,
		ProjectKey:   service.ProjectKey,
	}
	err := model.DB.Create(&performance).Error
	if err != nil {
		return utils.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库操作失败，添加Performance监控记录失败",
		}
	}
	return utils.Response{
		Status: http.StatusOK,
		Msg:    "Performance监控记录添加成功",
	}
}
