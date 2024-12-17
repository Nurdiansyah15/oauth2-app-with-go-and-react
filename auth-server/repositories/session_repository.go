package repositories

import (
	"auth-server/models"
	"time"

	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(session *models.Session) string
	GetSessionByDeviceID(deviceID string) (*models.Session, string)
	DeactivateSessionByToken(sessionToken string) string
	UpdateSessionActivity(sessionID string) string
	GetActiveSessionByToken(sessionToken string) (*models.Session, string)
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (s *sessionRepository) CreateSession(session *models.Session) string {
	if err := s.db.Create(session).Error; err != nil {
		return "DATABASE_ERROR_500"
	}
	return ""
}

func (s *sessionRepository) GetSessionByDeviceID(deviceID string) (*models.Session, string) {
	var session models.Session

	if err := s.db.Where("device_id = ? AND is_active = ?", deviceID, true).
		Order("created_at DESC").
		First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "SESSION_NOT_FOUND_404"
		}
		return nil, "DATABASE_ERROR_500"
	}

	return &session, ""
}

func (s *sessionRepository) DeactivateSessionByToken(sessionToken string) string {
	result := s.db.Model(&models.Session{}).
		Where("session_token = ?", sessionToken).
		Updates(map[string]interface{}{
			"is_active":  false,
			"expired_at": time.Now(),
		})

	if result.Error != nil {
		return "DATABASE_ERROR_500"
	}

	if result.RowsAffected == 0 {
		return "SESSION_NOT_FOUND_404"
	}

	return ""
}

func (s *sessionRepository) UpdateSessionActivity(sessionID string) string {
	result := s.db.Model(&models.Session{}).
		Where("id = ?", sessionID).
		Update("last_activity", time.Now())

	if result.Error != nil {
		return "DATABASE_ERROR_500"
	}

	if result.RowsAffected == 0 {
		return "SESSION_NOT_FOUND_404"
	}

	return ""
}

func (s *sessionRepository) GetActiveSessionByToken(sessionToken string) (*models.Session, string) {
	var session models.Session

	if err := s.db.Where(
		"session_token = ? AND is_active = ? AND expired_at > ?",
		sessionToken,
		true,
		time.Now(),
	).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "SESSION_NOT_FOUND_404"
		}
		return nil, "DATABASE_ERROR_500"
	}

	return &session, ""
}
