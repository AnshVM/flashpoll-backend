package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"unique"`
	Email        string `gorm:"unique"`
	PasswordHash []byte
	PollsCreated []Poll   `gorm:"foreignKey:User"`
	Votes        []Option `gorm:"many2many:user_votes"`
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
