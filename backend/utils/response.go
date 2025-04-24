package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 响应状态码
const (
	SUCCESS               = 0
	ERROR_AUTH_CHECK_FAIL = 1001
	ERROR_LOGIN_FAIL      = 1002
	ERROR_USER_EXISTS     = 1003
	ERROR_GAME_NOT_FOUND  = 2001
	ERROR_ROOM_NOT_FOUND  = 2002
	ERROR_ROOM_FULL       = 2003
	ERROR_ALREADY_IN_ROOM = 2004
	ERROR_INTERNAL        = 9999
)

// 响应状态码对应的消息
var CodeMsg = map[int]string{
	SUCCESS:               "成功",
	ERROR_AUTH_CHECK_FAIL: "用户未登录或token无效",
	ERROR_LOGIN_FAIL:      "用户名或密码错误",
	ERROR_USER_EXISTS:     "用户名已存在",
	ERROR_GAME_NOT_FOUND:  "游戏不存在",
	ERROR_ROOM_NOT_FOUND:  "房间不存在",
	ERROR_ROOM_FULL:       "房间已满",
	ERROR_ALREADY_IN_ROOM: "已在房间中",
	ERROR_INTERNAL:        "服务器内部错误",
}

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// JSONSuccess 返回成功响应
func JSONSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Data: data,
		Msg:  CodeMsg[SUCCESS],
	})
}

// JSONSuccessWithMsg 返回带自定义消息的成功响应
func JSONSuccessWithMsg(c *gin.Context, data interface{}, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Data: data,
		Msg:  msg,
	})
}

// JSONError 返回错误响应
func JSONError(c *gin.Context, code int) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: nil,
		Msg:  CodeMsg[code],
	})
}

// JSONErrorWithMsg 返回带自定义消息的错误响应
func JSONErrorWithMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Data: nil,
		Msg:  msg,
	})
}
