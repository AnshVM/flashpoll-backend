package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/AnshVM/flashpoll-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

func CreateConnection() error {
	dsn := strings.Join([]string{
		"host=localhost",
		"user=postgres",
		fmt.Sprintf("password=%s", os.Getenv("DB_PASSWORD")),
		"dbname=flashpoll",
		"port=5432",
		"sslmode=disable",
		"TimeZone=Asia/Kolkata",
	}, " ")

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	db.AutoMigrate(&models.User{})

	return nil
}

func Create(val interface{}) error {
	return db.Create(val).Error
}

func FindOneWhere(query *models.User, dest *models.User) *gorm.DB {
	return db.Where(query).First(dest)
}
