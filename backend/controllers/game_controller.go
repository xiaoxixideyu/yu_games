package controllers

import (
	"games/services"
	"games/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GameController 游戏控制器
type GameController struct {
	gameService *services.GameService
	roomService *services.RoomService
}

// NewGameController 创建游戏控制器实例
func NewGameController(gameService *services.GameService, roomService *services.RoomService) *GameController {
	return &GameController{
		gameService: gameService,
		roomService: roomService,
	}
}

// GetGameList 获取游戏列表
func (c *GameController) GetGameList(ctx *gin.Context) {
	// 获取游戏列表
	games, err := c.gameService.GetGameList()
	if err != nil {
		utils.JSONErrorWithMsg(ctx, utils.ERROR_INTERNAL, "获取游戏列表失败")
		return
	}

	utils.JSONSuccess(ctx, games)
}

// GetRoomList 获取游戏房间列表
func (c *GameController) GetRoomList(ctx *gin.Context) {
	// 获取游戏ID
	gameIDStr := ctx.Param("gameId")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil {
		utils.JSONErrorWithMsg(ctx, utils.ERROR_INTERNAL, "无效的游戏ID")
		return
	}

	// 获取房间列表
	rooms, err := c.roomService.GetRoomList(uint(gameID))
	if err != nil {
		if err.Error() == "游戏不存在" {
			utils.JSONError(ctx, utils.ERROR_GAME_NOT_FOUND)
		} else {
			utils.JSONErrorWithMsg(ctx, utils.ERROR_INTERNAL, "获取房间列表失败")
		}
		return
	}

	utils.JSONSuccess(ctx, rooms)
}

// CreateRoom 创建游戏房间
func (c *GameController) CreateRoom(ctx *gin.Context) {
	// 获取游戏ID
	gameIDStr := ctx.Param("gameId")
	gameID, err := strconv.ParseUint(gameIDStr, 10, 64)
	if err != nil {
		utils.JSONErrorWithMsg(ctx, utils.ERROR_INTERNAL, "无效的游戏ID")
		return
	}

	// 获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.JSONError(ctx, utils.ERROR_AUTH_CHECK_FAIL)
		return
	}

	// 获取请求参数
	var req struct {
		Name string `json:"name" binding:"required,min=1,max=30"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONErrorWithMsg(ctx, utils.ERROR_INTERNAL, "无效的请求参数")
		return
	}

	// 创建房间
	room, err := c.roomService.CreateRoom(req.Name, uint(gameID), userID.(uint))
	if err != nil {
		if err.Error() == "游戏不存在" {
			utils.JSONError(ctx, utils.ERROR_GAME_NOT_FOUND)
		} else {
			utils.JSONErrorWithMsg(ctx, utils.ERROR_INTERNAL, "创建房间失败: "+err.Error())
		}
		return
	}

	utils.JSONSuccessWithMsg(ctx, room, "创建成功")
}

// JoinRoom 加入游戏房间
func (c *GameController) JoinRoom(ctx *gin.Context) {
	// 获取房间ID
	roomIDStr := ctx.Param("roomId")
	roomID, err := strconv.ParseUint(roomIDStr, 10, 64)
	if err != nil {
		utils.JSONErrorWithMsg(ctx, utils.ERROR_INTERNAL, "无效的房间ID")
		return
	}

	// 获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.JSONError(ctx, utils.ERROR_AUTH_CHECK_FAIL)
		return
	}

	// 加入房间
	room, err := c.roomService.JoinRoom(uint(roomID), userID.(uint))
	if err != nil {
		switch err.Error() {
		case "房间不存在":
			utils.JSONError(ctx, utils.ERROR_ROOM_NOT_FOUND)
		case "房间已满":
			utils.JSONError(ctx, utils.ERROR_ROOM_FULL)
		case "已在房间中":
			utils.JSONError(ctx, utils.ERROR_ALREADY_IN_ROOM)
		default:
			utils.JSONErrorWithMsg(ctx, utils.ERROR_INTERNAL, "加入房间失败: "+err.Error())
		}
		return
	}

	utils.JSONSuccessWithMsg(ctx, room, "加入成功")
}
