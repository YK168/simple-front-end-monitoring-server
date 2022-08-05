package model

import (
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Title      string
	ProjectKey string `gorm:"index; not null"`
	User       User   `gorm:"ForeignKey: ProjectKey"`
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
