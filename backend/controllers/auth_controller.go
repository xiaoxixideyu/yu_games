package controllers

import (
	"games/services"
	"games/utils"

	"github.com/gin-gonic/gin"
)

// AuthController 用户认证控制器
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController 创建用户认证控制器实例
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	// 获取请求参数
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=20"`
		Password string `json:"password" binding:"required,min=6,max=20"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONErrorWithMsg(ctx, utils.ERROR_INTERNAL, "无效的请求参数")
		return
	}

	// 执行注册
	err := c.authService.Register(req.Username, req.Password)
	if err != nil {
		if err.Error() == "用户名已存在" {
			utils.JSONError(ctx, utils.ERROR_USER_EXISTS)
		} else {
			utils.JSONErrorWithMsg(ctx, utils.ERROR_INTERNAL, err.Error())
		}
		return
	}

	utils.JSONSuccessWithMsg(ctx, nil, "注册成功")
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	// 获取请求参数
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.JSONErrorWithMsg(ctx, utils.ERROR_INTERNAL, "无效的请求参数")
		return
	}

	// 执行登录
	token, user, err := c.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.JSONError(ctx, utils.ERROR_LOGIN_FAIL)
		return
	}

	// 返回令牌和用户信息
	utils.JSONSuccess(ctx, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}
