package main

import (
	"fmt"
	"games/config"
	"games/controllers"
	"games/models"
	"games/routes"
	"games/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	appConfig := config.LoadConfig()

	// 初始化数据库
	config.InitDatabase(appConfig.Database)

	// 自动迁移数据库表结构
	err := config.DB.AutoMigrate(
		&models.User{},
		&models.Game{},
		&models.Room{},
		&models.RoomPlayer{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 初始化预设游戏
	initGames()

	// 创建Gin实例
	router := gin.Default()

	// 创建服务实例
	authService := services.NewAuthService()
	gameService := services.NewGameService()
	roomService := services.NewRoomService()
	pokerService := services.NewPokerService(roomService)
	wsService := services.NewWSService(roomService, pokerService)

	// 创建控制器实例
	authController := controllers.NewAuthController(authService)
	gameController := controllers.NewGameController(gameService, roomService)
	wsController := controllers.NewWSController(wsService)

	// 设置路由
	routes.SetupRoutes(router, authController, gameController, wsController)

	// 启动服务器
	addr := fmt.Sprintf(":%d", appConfig.Server.Port)
	log.Printf("服务器启动，监听地址: %s", addr)
	err = router.Run(addr)
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// 初始化预设游戏
func initGames() {
	var count int64
	config.DB.Model(&models.Game{}).Count(&count)
	if count > 0 {
		return
	}

	games := []models.Game{
		{
			Name:        "德州扑克",
			Description: "一种风靡全球的扑克游戏，玩家可以利用自己的两张底牌和五张公共牌组成最好的牌型",
			ImageURL:    "/images/games/poker.jpg",
			MinPlayers:  2,
			MaxPlayers:  10,
		},
		{
			Name:        "斗地主",
			Description: "中国流行的扑克游戏，三人游戏，一方为地主，另外两人为农民",
			ImageURL:    "/images/games/doudizhu.jpg",
			MinPlayers:  3,
			MaxPlayers:  3,
		},
	}

	for _, game := range games {
		config.DB.Create(&game)
	}
}
