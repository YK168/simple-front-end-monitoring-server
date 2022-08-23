package service

import (
	"net/http"
	"simple_front_end_monitoring_server/model"
	"simple_front_end_monitoring_server/utils"
	"strconv"
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
	XhrInfo Xhr    `form:"xhrInfo" json:"xhrInfo"`
}

type Xhr struct {
	// API请求耗时
	Duration int `form:"duration" json:"duration"`
	// API请求结果类型，error代表请求失败，load代表请求成功
	EventType string `form:"eventType" json:"eventType" binding:"required"`

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
	Title string `form:"title" json:"title"`
	// 上报数据所属URL
	URL string `form:"url" json:"url"`
	// 报错时间
	TimeStamp int64 `form:"timestamp" json:"timestamp"`
	// 根据Cookie来区分不同页面？
	Cookie     string `form:"cookie" json:"cookie"`
	ProjectKey string `form:"projectKey" json:"projectKey"`
	Times      Times  `form:"times" json:"times"`
}

type Times struct {
	AnalysisTime string `form:"analysisTime" json:"analysisTime"`
	AppcacheTime string `form:"appcacheTime" json:"appcacheTime"`
	BlankTime    string `form:"blankTime" json:"blankTime"`
	DnsTime      string `form:"dnsTime" json:"dnsTime"`
	DomReadyTime string `form:"domReadyTime" json:"domReadyTime"`
	LoadPageTime string `form:"loadPageTime" json:"loadPageTime"`
	RedirectTime string `form:"redirectTime" json:"redirectTime"`
	ReqTime      string `form:"reqTime" json:"reqTime"`
	TcpTime      string `form:"tcpTime" json:"tcpTime"`
	TtfbTime     string `form:"ttfbTime" json:"ttfbTime"`
	UnloadTim    string `form:"uploadTime" json:"uploadTime"`
}

type AccessService struct {
	// 项目名称
	Title string `form:"title" json:"title"`
	// 报错时路由地址
	URL string `form:"url" json:"url"`
	// pv
	ErrType string `form:"errType" json:"errType"`
	// 报错时间
	TimeStamp int64 `form:"timestamp" json:"timestamp"`
	// 根据Cookie来区分不同页面？
	Cookie     string `form:"cookie" json:"cookie"`
	ProjectKey string `form:"projectKey" json:"projectKey"`
}

type BlankService struct {
	// 项目名称
	Title string `form:"title" json:"title"`
	// 报错时路由地址
	URL string `form:"url" json:"url"`
	// 是否发生白屏错误
	WhiteErr bool `form:"whiteError" json:"whiteError"`
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
	// 前端保证数据合法性
	analysisTime, _ := strconv.ParseFloat(service.Times.AnalysisTime, 32)
	appcacheTime, _ := strconv.ParseFloat(service.Times.AppcacheTime, 32)
	blankTime, _ := strconv.ParseFloat(service.Times.BlankTime, 32)
	dnsTime, _ := strconv.ParseFloat(service.Times.DnsTime, 32)
	domReadyTime, _ := strconv.ParseFloat(service.Times.DomReadyTime, 32)
	loadPageTime, _ := strconv.ParseFloat(service.Times.LoadPageTime, 32)
	redirectTime, _ := strconv.ParseFloat(service.Times.RedirectTime, 32)
	reqTime, _ := strconv.ParseFloat(service.Times.ReqTime, 32)
	tcpTime, _ := strconv.ParseFloat(service.Times.TcpTime, 32)
	ttfbTime, _ := strconv.ParseFloat(service.Times.TtfbTime, 32)
	unloadTim, _ := strconv.ParseFloat(service.Times.UnloadTim, 32)
	performance := model.Performance{
		Title:        service.Title,
		AnalysisTime: float32(analysisTime),
		AppcacheTime: float32(appcacheTime),
		BlankTime:    float32(blankTime),
		DnsTime:      float32(dnsTime),
		DomReadyTime: float32(domReadyTime),
		LoadPageTime: float32(loadPageTime),
		RedirectTime: float32(redirectTime),
		ReqTime:      float32(reqTime),
		TcpTime:      float32(tcpTime),
		TtfbTime:     float32(ttfbTime),
		UnloadTim:    float32(unloadTim),
		URL:          service.URL,
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

func (service *AccessService) Report() utils.Response {
	access := model.Access{
		Title:      service.Title,
		URL:        service.URL,
		ErrType:    service.ErrType,
		TimeStamp:  service.TimeStamp,
		Cookie:     service.Cookie,
		ProjectKey: service.ProjectKey,
	}
	err := model.DB.Create(&access).Error
	if err != nil {
		return utils.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库操作失败，添加PV/UV监控记录失败",
		}
	}
	return utils.Response{
		Status: http.StatusOK,
		Msg:    "PV/UV监控记录添加成功",
	}
}

func (service *BlankService) Report() utils.Response {
	blankErr := model.BlankError{
		Title:      service.Title,
		URL:        service.URL,
		WhiteErr:   service.WhiteErr,
		TimeStamp:  service.TimeStamp,
		Cookie:     service.Cookie,
		ProjectKey: service.ProjectKey,
	}
	err := model.DB.Create(&blankErr).Error
	if err != nil {
		return utils.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库操作失败，添加白屏Error监控记录失败",
		}
	}
	return utils.Response{
		Status: http.StatusOK,
		Msg:    "白屏Error监控记录添加成功",
	}
}
