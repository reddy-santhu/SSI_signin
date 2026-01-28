package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"ssi-signin/backend/models"
)

type SessionFinder interface {
	FindByToken(token string) (*models.Session, error)
}

type AuthMiddleware struct {
	sessions SessionFinder
}

func NewAuthMiddleware(sessions SessionFinder) *AuthMiddleware {
	return &AuthMiddleware{
		sessions: sessions,
	}
}

func (m *AuthMiddleware) RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid authorization header format",
			})
		}

		token := parts[1]
		session, err := m.sessions.FindByToken(token)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to validate session",
			})
		}

		if session == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid or expired session",
			})
		}

		if time.Now().After(session.ExpiresAt) {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Session expired",
			})
		}

		c.Set("session", session)
		c.Set("user_id", session.UserID)

		return next(c)
	}
}
