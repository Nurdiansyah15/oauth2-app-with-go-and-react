// app-one/backend/middleware/auth.go
package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const AUTH_SERVER = "http://auth-server-container:8080"

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "No authorization header"})
			c.Abort()
			return
		}

		// Hit endpoint validasi
		req, err := http.NewRequest("GET", AUTH_SERVER+"/oauth/validate", nil)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create request"})
			c.Abort()
			return
		}

		fmt.Println("token", token)
		fmt.Println("token", req)

		req.Header.Set("Authorization", token)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to validate token"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		type result struct {
			Valid  bool   `json:"valid"`
			UserID string `json:"user_id"`
		}

		var apiResponse struct {
			Status     string `json:"status"`
			StatusCode int    `json:"status_code"`
			Message    string `json:"message"`
			Data       result `json:"data"`
		}

		// Parse response untuk dapat user_id

		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			c.JSON(500, gin.H{"error": "Failed to parse validation response"})
			c.Abort()
			return
		}

		// Set user_id ke context untuk digunakan handler
		c.Set("user_id", apiResponse.Data.UserID)
		c.Next()
	}
}
