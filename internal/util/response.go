package util

import "github.com/gin-gonic/gin"

const (
	CodeSuccess = 0
	CodeError   = 1
)

type Response struct {
	Code       int    `json:"code"`
	Data       any    `json:"data,omitempty"`
	Pagination any    `json:"pagination,omitempty"`
	Message    string `json:"message,omitempty"`
}

func SuccessResponse(c *gin.Context, data any) {
	c.JSON(200, Response{
		Code: 0,
		Data: data,
	})
	c.Abort()
}

func SuccessData(c *gin.Context, data any, pagination any) {
	c.JSON(200, Response{
		Code:       0,
		Data:       data,
		Pagination: pagination,
	})
	c.Abort()
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
	})
	c.Abort()
}
