package model

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	// 手机号码
	Number string `gorm:"unique"`
	// 密文
	Passwd string
}

// 密码加密
func (user *User) SetPasswd(passwd string) error {
	b, err := bcrypt.GenerateFromPassword([]byte(passwd), 12)
	if err != nil {
		return err
	}
	user.Passwd = string(b)
	return nil
}

// 校验密码
func (user *User) CheckPasswd(passwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(passwd))
	if err != nil {
		log.Println("校验密码失败:", err)
		return false
	}
	return true
}
