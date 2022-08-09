package service

import (
	"net/http"
	"simple_front_end_monitoring_server/model"
	"simple_front_end_monitoring_server/utils"
)

type JsErrorService struct {
	// 项目名称
	Title string `form:"title" json:"title"`
	// 错误信息
	Message string `form:"message" json:"message"`
	// 报错时路由地址
	URL string `form:"url" json:"url"`
	// 报错代码行数
	Position string `form:"position" json:"position"`
	// 报错文件
	FileName string `form:"filename" json:"filename"`
	// JsError or PromiseError
	ErrType string `form:"errType" json:"errType"`
	// 报错时间
	TimeStamp int64 `form:"timestamp" json:"timestamp"`
	// 根据Cookie来区分不同页面？
	Cookie     string `form:"cookie" json:"cookie"`
	ProjectKey string `form:"projectKey" json:"projectKey"`
}

type ApiErrorService struct {
	// 根据Cookie来区分不同页面？
	Cookie     string `form:"cookie" json:"cookie"`
	ProjectKey string `form:"projectKey" json:"projectKey"`
	// 报错时间
	TimeStamp int64 `form:"timestamp" json:"timestamp"`
	// 项目名称
	Title string `form:"title" json:"title"`
	// 报错时路由地址
	URL     string `form:"url" json:"url"`
	XhrInfo Xhr    `form:"xhr" json:"xhr"`
}

type Xhr struct {
	// API请求耗时
	Duration int `form:"duration" json:"duration"`
	// API请求结果类型
	EventType string `form:"eventType" json:"eventType"`

	Kind string `form:"kind" json:"kind"`
	// POST请求的参数
	Params string `form:"params" json:"params"`
	// API请求的地址，GET请求参数加在后面
	Pathname string `form:"pathname" json:"pathname"`
	Response string `form:"response" json:"response"`
	// API请求的状态
	Status string `form:"status" json:"status"`
	// 什么类型的API请求工具
	ReqType string `form:"type" json:"type"`
}

type SourceErrorService struct {
	// 项目名称
	Title string `form:"title" json:"title"`
	// 报错时路由地址
	URL string `form:"url" json:"url"`
	// 报错文件
	FileName string `form:"filename" json:"filename"`
	// 报错资源标签
	TagName string `form:"tagName" json:"tagName"`
	// JsError or PromiseError
	ErrType string `form:"errType" json:"errType"`
	// 报错时间
	TimeStamp int64 `form:"timestamp" json:"timestamp"`
	// 根据Cookie来区分不同页面？
	Cookie     string `form:"cookie" json:"cookie"`
	ProjectKey string `form:"projectKey" json:"projectKey"`
}

type PerformanceService struct {
	// 项目名称
	Title        string  `form:"title" json:"title"`
	AnalysisTime float32 `form:"analysisTime" json:"analysisTime"`
	AppcacheTime float32 `form:"appcacheTime" json:"appcacheTime"`
	BlankTime    float32 `form:"blankTime" json:"blankTime"`
	DomReadyTime float32 `form:"domReadyTime" json:"domReadyTime"`
	LoadPageTime float32 `form:"loadPageTime" json:"loadPageTime"`
	RedirectTime float32 `form:"redirectTime" json:"redirectTime"`
	ReqTime      float32 `form:"reqTime" json:"reqTime"`
	TcpTime      float32 `form:"tcpTime" json:"tcpTime"`
	TtfbTime     float32 `form:"ttfbTime" json:"ttfbTime"`
	UnloadTim    float32 `form:"uploadTime" json:"uploadTime"`
	// 报错时间
	TimeStamp int64 `form:"timestamp" json:"timestamp"`
	// 根据Cookie来区分不同页面？
	Cookie     string `form:"cookie" json:"cookie"`
	ProjectKey string `form:"projectKey" json:"projectKey"`
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
		URL:        service.URL,
		TimeStamp:  service.TimeStamp,
		Cookie:     service.Cookie,
		ProjectKey: service.ProjectKey,
		Duration:   service.XhrInfo.Duration,
		EventType:  service.XhrInfo.EventType,
		Kind:       service.XhrInfo.Kind,
		Params:     service.XhrInfo.Params,
		Pathname:   service.XhrInfo.Pathname,
		Response:   service.XhrInfo.Response,
		Status:     service.XhrInfo.Status,
		ReqType:    service.XhrInfo.ReqType,
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
