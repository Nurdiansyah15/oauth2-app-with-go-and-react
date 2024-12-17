package handlers

import (
	"auth-server/services"
	"auth-server/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type authHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *authHandler {
	return &authHandler{authService: authService}
}

func (h *authHandler) LoginHandler(c *gin.Context) {
	type LoginInput struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var input LoginInput

	// Handle JSON binding error
	if err := c.ShouldBindJSON(&input); err != nil {
		// Specific validation errors
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)

			for _, e := range validationErrors {
				field := strings.ToLower(e.Field()) // lowercase field name for consistency
				switch {
				case field == "username" && e.Tag() == "required":
					errorMessages[field] = "Username is required"
				case field == "password" && e.Tag() == "required":
					errorMessages[field] = "Password is required"
				default:
					errorMessages[field] = fmt.Sprintf("Invalid value for %s", field)
				}
			}

			utils.ApiErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", errorMessages)
			return
		}

		// General JSON parsing error
		utils.ApiErrorResponse(c, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format")
		return
	}

	// Trim spaces from input
	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)

	// Call login service
	user, errCode := h.authService.Login(input.Username, input.Password)
	if errCode != "" {
		switch errCode {
		case "USER_NOT_FOUND_404", "WRONG_PASSWORD_401":
			utils.ApiErrorResponse(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid username or password")
		default:
			utils.ApiErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "An error occurred during login")
		}
		return
	}

	utils.ApiResponse(c, http.StatusOK, "Login successful", user)
}
