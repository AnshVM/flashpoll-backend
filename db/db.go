package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/AnshVM/flashpoll-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func CreateConnection() error {
	dsn := strings.Join([]string{
		fmt.Sprintf("host=%s", os.Getenv("DB_HOST")),
		fmt.Sprintf("user=%s", os.Getenv("DB_USER")),
		fmt.Sprintf("password=%s", os.Getenv("DB_PASSWORD")),
		fmt.Sprintf("dbname=%s", os.Getenv("DB_NAME")),
		fmt.Sprintf("port=%s", os.Getenv("DB_PORT")),
		"sslmode=require",
		"TimeZone=Asia/Kolkata",
	}, " ")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return err
	}

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Option{})
	DB.AutoMigrate(&models.Poll{})

	return nil
}

func Create(val interface{}) error {
	return DB.Create(val).Error
}

// more types to be added as more db models are created
func FindOneWhere[T models.User](query *T, dest *T) *gorm.DB {
	return DB.Where(query).First(dest)
}

func FindById[T models.User | models.Poll | models.Option](ID uint, dest *T) error {
	return DB.First(dest, ID).Error
}
