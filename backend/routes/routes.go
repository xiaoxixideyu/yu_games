package routes

import (
	"games/controllers"
	"games/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(
	router *gin.Engine,
	authController *controllers.AuthController,
	gameController *controllers.GameController,
	wsController *controllers.WSController,
) {
	// API路由组
	api := router.Group("/api")
	{
		// 用户认证相关路由
		api.POST("/register", authController.Register)
		api.POST("/login", authController.Login)

		// 游戏相关路由，需要认证
		game := api.Group("/games", middleware.JWTAuth())
		{
			game.GET("", gameController.GetGameList)
			game.GET("/:gameId/rooms", gameController.GetRoomList)
			game.POST("/:gameId/rooms", gameController.CreateRoom)
		}

		// 房间相关路由，需要认证
		room := api.Group("/rooms", middleware.JWTAuth())
		{
			room.POST("/:roomId/join", gameController.JoinRoom)
		}
	}

	// WebSocket路由
	router.GET("/ws", middleware.JWTAuth(), wsController.HandleWebSocket)
}
