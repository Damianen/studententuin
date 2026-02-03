package middlewares

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Code int	    `json:"code"`
	Message string  `json:"message"`
	Data any        `json:"data,omitempty"`
}

func Respond(c *gin.Context, httpStatus int, message string, data any) {
	c.JSON(httpStatus, APIResponse{
		Code: httpStatus,
		Message: message,
		Data: data,
	})
}
