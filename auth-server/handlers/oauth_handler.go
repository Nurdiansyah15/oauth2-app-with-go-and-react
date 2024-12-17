package handlers

import (
	"auth-server/models"
	"auth-server/repositories"
	"auth-server/services"
	"auth-server/utils"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type OauthHandler interface {
	LoginPage(c *gin.Context)
	LoginHandler(c *gin.Context)
	LogoutHandler(c *gin.Context)
	AuthorizeHandler(c *gin.Context)
	TokenHandler(c *gin.Context)
	ValidateTokenHandler(c *gin.Context)
}

type oauthHandler struct {
	clientRepo   repositories.ClientRepository
	sessionRepo  repositories.SessionRepository
	authCodeRepo repositories.AuthCodeRepository
	authService  services.AuthService
	oauthService services.OAuthService
	jwtService   services.JWTService
}

func NewOauthHandler(
	clientRepo repositories.ClientRepository,
	sessionRepo repositories.SessionRepository,
	authCodeRepo repositories.AuthCodeRepository,
	authService services.AuthService,
	oauthService services.OAuthService,
	jwtService services.JWTService,
) OauthHandler {
	return &oauthHandler{
		clientRepo:   clientRepo,
		sessionRepo:  sessionRepo,
		authCodeRepo: authCodeRepo,
		authService:  authService,
		oauthService: oauthService,
		jwtService:   jwtService,
	}
}

func (h *oauthHandler) LoginPage(c *gin.Context) {
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")

	// Validasi client
	client, errCode := h.clientRepo.GetClientByID(clientID)
	if errCode != "" || client.RedirectURI != redirectURI {
		utils.ApiErrorResponse(c, http.StatusBadRequest, "INVALID_CLIENT", "Invalid client or redirect URI")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"clientID":    clientID,
		"redirectURI": redirectURI,
	})
}

func (h *oauthHandler) LoginHandler(c *gin.Context) {
	// 	RedirectURI string `form:"redirect_uri"`
	// ClientID    string `form:"client_id"`

	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")

	client, errCode := h.clientRepo.GetClientByID(clientID)
	if errCode != "" || client.RedirectURI != redirectURI {
		utils.ApiErrorResponse(c, http.StatusBadRequest, "INVALID_CLIENT", "Invalid client or redirect URI")
		return
	}

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
	// Create session
	session := &models.Session{
		ID:           uuid.New().String(),
		UserID:       user.ID,
		ClientID:     clientID,
		SessionToken: uuid.New().String(),
		BrowserInfo:  c.Request.UserAgent(),
		IPAddress:    c.ClientIP(),
		LastActivity: time.Now(),
		CreatedAt:    time.Now(),
		ExpiredAt:    time.Now().Add(24 * time.Hour),
		IsActive:     true,
	}

	if errCode := h.sessionRepo.CreateSession(session); errCode != "" {
		utils.ApiErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "Failed to create session")
		return
	}

	// Set session cookie
	c.SetCookie(
		"session_token",
		session.SessionToken,
		int(time.Until(session.ExpiredAt).Seconds()),
		"/",
		"",    // domain
		false, // secure (set to true in production)
		true,  // httpOnly
	)

	// Generate auth code
	authCode, errCode := h.oauthService.GenerateAuthCodeWithSession(user.ID, clientID)
	if errCode != "" {
		utils.ApiErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to create auth code")
		return
	}

	// Redirect with auth code
	// c.Redirect(http.StatusFound, fmt.Sprintf("%s?code=%s", redirectURI, authCode.Code))

	utils.ApiResponse(c, http.StatusOK, "Login Successfully", gin.H{
		"auth_code":    authCode.Code,
		"redirect_uri": redirectURI,
		"redirect_to":  redirectURI + "?code=" + authCode.Code,
	})
}

func (h *oauthHandler) LogoutHandler(c *gin.Context) {
	// Get session token from cookie
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		// utils.ApiErrorResponse(c, http.StatusBadRequest, "INVALID_SESSION", "No active session")
		// return
		redirectURI := c.Query("redirect_uri")
		if redirectURI != "" {
			c.Redirect(http.StatusFound, redirectURI)
			return
		}

		utils.ApiResponse(c, http.StatusOK, "Logout successful", nil)
	}

	// Deactivate session
	errCode := h.sessionRepo.DeactivateSessionByToken(sessionToken)
	if errCode != "" {
		utils.ApiErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "Failed to logout")
		return
	}

	// Clear session cookie
	c.SetCookie(
		"session_token",
		"", // empty value
		-1, // negative maxAge = delete cookie
		"/",
		"",    // domain
		false, // secure
		true,  // httpOnly
	)

	// Redirect ke login page atau response success
	redirectURI := c.Query("redirect_uri")
	if redirectURI != "" {
		c.Redirect(http.StatusFound, redirectURI)
		return
	}

	utils.ApiResponse(c, http.StatusOK, "Logout successful", nil)
}

func (h *oauthHandler) AuthorizeHandler(c *gin.Context) {
	clientID := c.Query("client_id")
	redirectURI := c.Query("redirect_uri")

	// Validasi client
	client, errCode := h.clientRepo.GetClientByID(clientID)
	if errCode != "" || client.RedirectURI != redirectURI {
		utils.ApiErrorResponse(c, http.StatusBadRequest, "INVALID_CLIENT", "Invalid client or redirect URI")
		return
	}

	// Get session token dari cookie
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		c.Redirect(http.StatusFound, "/login?client_id="+clientID+"&redirect_uri="+redirectURI)
		return
	}

	// Cek session berdasarkan token
	session, errCode := h.sessionRepo.GetActiveSessionByToken(sessionToken)
	if errCode != "" {
		c.Redirect(http.StatusFound, "/login?client_id="+clientID+"&redirect_uri="+redirectURI)
		return
	}

	// Update last activity
	if errCode := h.sessionRepo.UpdateSessionActivity(session.ID); errCode != "" {
		utils.ApiErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "Failed to update session")
		return
	}

	// Generate auth code
	authCode, errCode := h.oauthService.GenerateAuthCodeWithSession(session.UserID, clientID)
	if errCode != "" {
		utils.ApiErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "Failed to generate auth code")
		return
	}

	c.Redirect(http.StatusFound, redirectURI+"?code="+authCode.Code)
	// utils.ApiResponse(c, http.StatusOK, "Login Successfully",gin.H{
	// 	"auth_code" : authCode.Code,
	// 	"redirect_uri" : redirectURI,
	// 	"redirect_to" : redirectURI+"?code="+authCode.Code,
	// })

}

func (h *oauthHandler) TokenHandler(c *gin.Context) {
	code := c.PostForm("code")
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")

	// Validate client dan secret
	client, errCode := h.clientRepo.GetClientByID(clientID)
	if errCode != "" {
		utils.ApiErrorResponse(c, http.StatusUnauthorized, "INVALID_CLIENT", "Invalid client ID")
		return
	}

	if client.Secret != clientSecret {
		utils.ApiErrorResponse(c, http.StatusUnauthorized, "INVALID_SECRET", "Invalid client secret")
		return
	}

	// Validate auth code
	authCode, errCode := h.authCodeRepo.GetAuthCodeByCode(code)
	if errCode != "" {
		utils.ApiErrorResponse(c, http.StatusBadRequest, "INVALID_CODE", "Invalid authorization code")
		return
	}

	if authCode.ClientID != clientID {
		utils.ApiErrorResponse(c, http.StatusUnauthorized, "INVALID_CLIENT", "Auth code not issued for this client")
		return
	}

	if time.Now().After(authCode.ExpiresAt) {
		utils.ApiErrorResponse(c, http.StatusBadRequest, "CODE_EXPIRED", "Authorization code expired")
		return
	}

	// Generate JWT
	token, err := h.jwtService.GenerateToken(authCode.UserID)
	if err != "" {
		utils.ApiErrorResponse(c, http.StatusInternalServerError, "TOKEN_ERROR", "Could not generate token")
		return
	}

	// Delete used auth code
	if errCode := h.authCodeRepo.DeleteAuthCode(code); errCode != "" {
		utils.ApiErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "Failed to invalidate auth code")
		return
	}

	utils.ApiResponse(c, http.StatusOK, "Token generated successfully", gin.H{
		"access_token": token,
		"token_type":   "Bearer",
	})
}

func (h *oauthHandler) ValidateTokenHandler(c *gin.Context) {
	tokenString := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)

	claims, err := h.jwtService.ValidateToken(tokenString)
	if err != "" {
		utils.ApiErrorResponse(c, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid token")
		return
	}

	userID := claims["user_id"].(string)

	utils.ApiResponse(c, http.StatusOK, "Token validated successfully", gin.H{
		"valid":   true,
		"user_id": userID,
	})
}
