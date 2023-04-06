package main

import (
	"github.com/AnshVM/flashpoll-backend/controllers"
	"github.com/AnshVM/flashpoll-backend/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	err := db.CreateConnection()

	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)

	router.Run(":8080")
}
