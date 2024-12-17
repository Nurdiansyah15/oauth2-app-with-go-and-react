package handlers

import (
	"auth-server/services"
	"auth-server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUserInfoHandler(c *gin.Context)
}

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (h *userHandler) GetUserInfoHandler(c *gin.Context) {
	userID, isExist := c.Get("user_id")
	if !isExist {
		utils.ApiErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}
	user, errCode := h.userService.GetUserById(userID.(string))
	if errCode != "" {
		switch errCode {
		case "USER_NOT_FOUND_404":
			utils.ApiErrorResponse(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
		default:
			utils.ApiErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "Failed to get user info")
		}
		return
	}

	utils.ApiResponse(c, http.StatusOK, "User info retrieved successfully", gin.H{
		"user": user,
	})
}
