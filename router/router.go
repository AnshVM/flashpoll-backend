package router

import (
	"github.com/AnshVM/flashpoll-backend/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	router.Use(cors.New(config))

	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.POST("/poll", controllers.CreatePoll)
	router.GET("/refresh", controllers.RefreshTokens)
	router.GET("/poll/:id", controllers.GetPollById)
	router.POST("/poll/submit", controllers.SubmitVote)
	router.POST("logout", controllers.Logout)

	return router
}
