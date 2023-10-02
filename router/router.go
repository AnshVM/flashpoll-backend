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
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	router.Use(cors.New(config))

	wsHub := ws.NewHub()
	go wsHub.Run()

	api := router.Group("/api")
	{
		api.POST("/signup", controllers.Signup)
		api.POST("/login", controllers.Login)
		api.POST("/poll", controllers.Auth, controllers.CreatePoll)
		api.POST("/refresh", controllers.RefreshTokens)
		api.GET("/poll/:id", controllers.Auth, controllers.GetPollById)
		api.POST("logout", controllers.Logout)

		api.GET("/ws/:id", func(ctx *gin.Context) {
			pollID, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
			ws.ServeWs(wsHub, ctx.Writer, ctx.Request, uint(pollID))
		})
		api.POST("/poll/submit", controllers.Auth, func(ctx *gin.Context) {
			controllers.SubmitVote(ctx, wsHub)
		})
	}

	return router
}
