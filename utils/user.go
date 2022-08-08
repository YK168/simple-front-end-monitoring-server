package utils

import (
	"crypto/md5"
	"fmt"
	"simple_front_end_monitoring_server/model"
)

type User struct {
	ID       uint   `json:"id" form:"id" example:"1"`                   // 用户ID
	Number   string `json:"number" form:"number" example:"12345678910"` // 用户手机号
	Status   string `json:"status" form:"status"`                       // 用户状态
	CreateAt int64  `json:"create_at" form:"create_at"`                 // 创建时间
}

func BuildUser(user model.User) User {
	return User{
		ID:       user.ID,
		Number:   user.Number,
		CreateAt: user.CreatedAt.Unix(),
	}
}

func MD5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
