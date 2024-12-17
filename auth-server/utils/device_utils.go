package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

func GenerateDeviceID(c *gin.Context) string {
	userAgent := c.Request.UserAgent()
	ip := c.ClientIP()
	// Bisa tambah identifikasi lain seperti Accept-Language, dll

	// Generate unique hash
	data := userAgent + ip
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
