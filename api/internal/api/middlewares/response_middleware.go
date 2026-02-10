package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Respond(c *gin.Context, httpStatus int, message string, data any) {
	c.JSON(httpStatus, APIResponse{
		Code:    httpStatus,
		Message: message,
		Data:    data,
	})
}

func BindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindBodyWithJSON(obj); err != nil {
		fmt.Println(err.Error())
		Respond(c, http.StatusBadRequest, "invalid JSON or missing values", nil)
		return false
	}
	return true
}

func CheckOwnership(c *gin.Context, check func(ctx context.Context, userID, resourceID string) (bool, error), userID, resourceID, resourceName string) bool {
	allowed, err := check(c.Request.Context(), userID, resourceID)
	if errors.Is(err, gorm.ErrRecordNotFound) || (err != nil && strings.Contains(err.Error(), "invalid input syntax")) {
		Respond(c, http.StatusNotFound, resourceName+" not found", nil)
		return false
	}
	if err != nil {
		fmt.Println(err.Error())
		Respond(c, http.StatusInternalServerError, "internal server error", nil)
		return false
	}
	if !allowed {
		Respond(c, http.StatusForbidden, "unauthorized", nil)
		return false
	}
	return true
}
