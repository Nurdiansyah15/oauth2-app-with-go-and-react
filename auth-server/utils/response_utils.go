package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type ErrorResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"status_code"`
	ErrorCode  string      `json:"error_code"`
	Error      interface{} `json:"error"`
}

func ApiResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := Response{
		Status:     "success",
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}

	c.JSON(statusCode, response)
}

func ApiErrorResponse(c *gin.Context, statusCode int, errorCode string, error interface{}) {

	response := ErrorResponse{
		Status:     "error",
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Error:      error,
	}

	c.JSON(statusCode, response)
}
