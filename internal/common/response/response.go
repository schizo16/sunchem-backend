package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Code    string      `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, APIResponse{Data: data})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, APIResponse{Message: message, Data: data})
}
