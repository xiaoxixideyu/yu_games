package controllers

import (
	"games/services"

	"github.com/gin-gonic/gin"
)

// WSController WebSocket控制器
type WSController struct {
	wsService *services.WSService
}

// NewWSController 创建WebSocket控制器实例
func NewWSController(wsService *services.WSService) *WSController {
	return &WSController{
		wsService: wsService,
	}
}

// HandleWebSocket 处理WebSocket连接
func (c *WSController) HandleWebSocket(ctx *gin.Context) {
	c.wsService.HandleConnection(ctx)
}
