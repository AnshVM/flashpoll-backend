package main

import (
	"github.com/AnshVM/flashpoll-backend/controllers"
	"github.com/AnshVM/flashpoll-backend/db"
	"github.com/gin-contrib/cors"
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
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(config))

	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)

	router.Run(":8080")
}
