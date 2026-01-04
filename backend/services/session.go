package services

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type SessionService struct{}

func NewSessionService() *SessionService {
	return &SessionService{}
}

func (s *SessionService) GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *SessionService) GetExpirationTime() time.Time {
	return time.Now().Add(24 * time.Hour)
}

