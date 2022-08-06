package model

import "log"

func migrate() {
	// 自动迁移模式
	// 建立与数据库的映射
	err := DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(
		&User{}, &Item{}, &JSError{}, &APIError{}, &SourceError{}, &Performance{},
	)
	if err != nil {
		log.Fatalln("数据库迁移失败", err)
	}
}
