package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AUTH_SERVER   = "http://auth-server-container:8080"
	CLIENT_ID     = "app-two-client"
	CLIENT_SECRET = "app-two-secret"
)

func HandleAuthCallback(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}

	if err := c.BindJSON(&req); err != nil {
		log.Printf("Error binding JSON in HandleAuthCallback: %v", err)
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Exchange code for token with auth server
	tokenResp, err := exchangeCodeForToken(req.Code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		c.JSON(500, gin.H{"error": "Token exchange failed"})
		return
	}

	c.JSON(200, tokenResp)
}

func exchangeCodeForToken(code string) (map[string]interface{}, error) {
	// Create form data
	formData := map[string][]string{
		"code":          {code},
		"client_id":     {CLIENT_ID},
		"client_secret": {CLIENT_SECRET},
		"grant_type":    {"authorization_code"},
	}

	// Send request to auth server
	resp, err := http.PostForm(AUTH_SERVER+"/oauth/token", formData)
	if err != nil {
		log.Printf("Error sending POST request to auth server: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding response from auth server: %v", err)
		return nil, err
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Auth server returned non-200 status: %d, response: %v", resp.StatusCode, result)
		return nil, fmt.Errorf("auth server returned status: %d", resp.StatusCode)
	}

	return result, nil
}

func GetMe(c *gin.Context) {
	// Mengambil token dari header request
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Println("Authorization header is missing in GetMe")
		c.JSON(401, gin.H{"error": "No authorization header"})
		return
	}

	// Membuat request baru ke AUTH_SERVER
	req, err := http.NewRequest("GET", AUTH_SERVER+"/api/me", nil)
	if err != nil {
		log.Printf("Error creating request in GetMe: %v", err)
		c.JSON(500, gin.H{"error": "Failed to create request"})
		return
	}

	// Meneruskan header Authorization
	req.Header.Set("Authorization", authHeader)

	// Melakukan request dengan http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to auth server in GetMe: %v", err)
		c.JSON(500, gin.H{"error": "Failed to fetch user data"})
		return
	}
	defer resp.Body.Close()

	// Check response status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Auth server returned non-200 status in GetMe: %d", resp.StatusCode)
		c.JSON(resp.StatusCode, gin.H{"error": "Unauthorized"})
		return
	}

	// Decode response
	type User struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	type apiResponse struct {
		Status     string `json:"status"`
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
		Data       struct {
			User User `json:"user"`
		} `json:"data"`
	}

	var userData apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		log.Printf("Error decoding user data in GetMe: %v", err)
		c.JSON(500, gin.H{"error": "Failed to decode user data"})
		return
	}

	c.JSON(200, gin.H{
		"user_id":  userData.Data.User.ID,
		"username": userData.Data.User.Username,
	})
}
