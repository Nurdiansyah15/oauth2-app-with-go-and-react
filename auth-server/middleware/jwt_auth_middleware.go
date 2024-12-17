package middleware

import (
	"auth-server/utils"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// sessionToken, err := c.Cookie("session_token")
		// fmt.Println("Session", sessionToken)

		// if err != nil {
		// 	c.Abort()
		// 	return
		// }

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ApiErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
			c.Abort()
			return
		}

		// Remove 'Bearer ' prefix
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			utils.ApiErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
			c.Next()
		} else {
			utils.ApiErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
			c.Abort()
			return
		}
	}
}
