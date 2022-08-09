package model

import (
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Title      string
	ProjectKey string `gorm:"index; not null"`
	// 对应的用户号码
	Number string `gorm:"index; not null"`
}

type JSError struct {
	gorm.Model
	// 项目名称
	Title string
	// 错误信息
	Message string
	// 报错时路由地址
	URL string
	// 报错代码行数
	Position string
	// 报错文件
	FileName string
	// JsError or PromiseError
	ErrType string
	// 报错时间
	TimeStamp int64
	// 根据Cookie来区分不同页面？
	Cookie     string
	ProjectKey string `gorm:"index; not null"`
}

type APIError struct {
	gorm.Model
	// 项目名称
	Title string
	// 报错时路由地址
	URL string
	// 报错时间
	TimeStamp int64
	// 根据Cookie来区分不同页面？
	Cookie     string
	ProjectKey string `gorm:"index; not null"`
	// API请求耗时
	Duration int
	// API请求结果类型
	EventType string
	Kind      string
	// POST请求的参数
	Params string
	// API请求的地址，GET请求参数加在后面
	Pathname string
	Response string
	// API请求的状态
	Status string
	// 什么类型的API请求工具
	ReqType string
}

type SourceError struct {
	gorm.Model
	// 项目名称
	Title string
	// 报错时路由地址
	URL string
	// 报错文件
	FileName string
	// 报错资源标签
	TagName string
	// JsError or PromiseError
	ErrType string
	// 报错时间
	TimeStamp int64
	// 根据Cookie来区分不同页面？
	Cookie     string
	ProjectKey string `gorm:"index; not null"`
}

type Performance struct {
	gorm.Model
	// 项目名称
	Title        string
	AnalysisTime float32
	AppcacheTime float32
	BlankTime    float32
	DomReadyTime float32
	LoadPageTime float32
	RedirectTime float32
	ReqTime      float32
	TcpTime      float32
	TtfbTime     float32
	UnloadTim    float32
	// 报错时间
	TimeStamp int64
	// 根据Cookie来区分不同页面？
	Cookie     string
	ProjectKey string `gorm:"index; not null"`
}
