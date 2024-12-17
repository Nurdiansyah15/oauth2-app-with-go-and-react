package utils

import (
	"auth-server/models"
	"math/rand"
	"time"
)

type LoginInput struct {
	Username    string `form:"username"`
	Password    string `form:"password"`
}

func GenerateAuthCode(userID string) models.AuthCode {
	return models.AuthCode{
		Code:      GenerateRandomString(32),
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Minute * 10),
	}
}

func GenerateRandomString(length int) string {
	var letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// func RenderLoginError(c *gin.Context, status int, message string, input interface{}) {
// 	if loginInput, ok := input.(LoginInput); ok {
// 		c.HTML(status, "login.html", gin.H{
// 			"error":        message,
// 			"username":     loginInput.Username,
// 			"clientID":    loginInput.ClientID,
// 			"redirectURI": loginInput.RedirectURI,
// 		})
// 		return
// 	}
// 	// Fallback if input is not LoginInput
// 	c.HTML(status, "login.html", gin.H{
// 		"error": message,
// 	})
// }
