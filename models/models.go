package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"unique"`
	Email        string `gorm:"unique"`
	PasswordHash []byte
	Polls        []Poll `gorm:"foreignKey:User"`
}

type Poll struct {
	gorm.Model
	Title   string
	Options []Option `gorm:"foreignKey:Poll"`
	User    uint
}

type Option struct {
	gorm.Model
	Name  string
	Count uint
	Poll  uint
}
