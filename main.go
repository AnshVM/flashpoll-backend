package main

import (
	"fmt"
	"net/http"

	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"unique"`
	PasswordHash []byte
}

type Signup struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {

	dsn := strings.Join([]string{
		"host=localhost",
		"user=postgres",
		"password=GunJedi_99",
		"dbname=flashpoll",
		"port=5432",
		"sslmode=disable",
		"TimeZone=Asia/Kolkata",
	}, " ")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
		return
	}

	db.AutoMigrate(&User{})

	router := gin.Default()

	router.POST("/signup", func(ctx *gin.Context) {
		var req Signup
		fmt.Println("recieve requst")
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request Body"})
		}
		fmt.Println(1)
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
		if err != nil {
			if errors.Is(err, bcrypt.ErrPasswordTooLong) {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "PasswordTooLong"})
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "An error occured"})
			return
		}

		err = db.Create(&User{Email: req.Email, PasswordHash: hash}).Error

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "EmailAlreadyInUse"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "request recived"})
	})

	router.Run(":8080")

}
