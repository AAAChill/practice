package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"practice/model"
)

// SuccessResponse
// @Description 成功返回
func SuccessResponse(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, model.BaseResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// ErrorResponse
// @Description 错误返回
func ErrorResponse(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(code, model.BaseResponse{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
