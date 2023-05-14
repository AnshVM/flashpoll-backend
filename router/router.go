package router

import (
	"strconv"

	"github.com/AnshVM/flashpoll-backend/controllers"
	"github.com/AnshVM/flashpoll-backend/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	router.Use(cors.New(config))

	wsHub := ws.NewHub()
	go wsHub.Run()

	router.POST("/signup", controllers.Signup)
	router.POST("/login", controllers.Login)
	router.POST("/poll", controllers.Auth, controllers.CreatePoll)
	router.GET("/refresh", controllers.RefreshTokens)
	router.GET("/poll/:id", controllers.Auth, controllers.GetPollById)
	router.POST("logout", controllers.Logout)

	router.GET("/ws/:id", func(ctx *gin.Context) {
		pollID, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
		ws.ServeWs(wsHub, ctx.Writer, ctx.Request, uint(pollID))
	})
	router.POST("/poll/submit", controllers.Auth, func(ctx *gin.Context) {
		controllers.SubmitVote(ctx, wsHub)
	})

	return router
}
