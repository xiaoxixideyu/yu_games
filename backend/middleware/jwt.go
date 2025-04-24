package middleware

import (
	"games/config"
	"games/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.JSONError(c, utils.ERROR_AUTH_CHECK_FAIL)
			c.Abort()
			return
		}

		// 检查token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.JSONError(c, utils.ERROR_AUTH_CHECK_FAIL)
			c.Abort()
			return
		}

		// 解析token
		token := parts[1]
		appConfig := config.LoadConfig()
		jwt := utils.NewJWT(appConfig.JWT.Secret)
		claims, err := jwt.ParseToken(token)
		if err != nil {
			utils.JSONError(c, utils.ERROR_AUTH_CHECK_FAIL)
			c.Abort()
			return
		}

		// 将用户ID存储在上下文中
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
