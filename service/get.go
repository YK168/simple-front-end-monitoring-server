package service

import (
	"log"
	"simple_front_end_monitoring_server/model"
)

type Searcher struct {
	ProjectKey     string
	StartTimeStamp int64
	EndTimeStamp   int64
}

func (s *Searcher) Search(mode any, data any) {
	err := model.DB.Model(mode).
		Where("project_key = ?", s.ProjectKey).
		Where("time_stamp > ?", s.StartTimeStamp-1).
		Where("time_stamp < ?", s.EndTimeStamp+1).
		Order("time_stamp").
		Find(data).Error
	if err != nil {
		log.Fatalln("查询数据库失败", err)
	}
}
