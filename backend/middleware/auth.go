package middleware

import (
	"net/http"
	"ssi-signin/backend/repositories"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	sessionRepo *repositories.SessionRepository
}

func NewAuthMiddleware(sessionRepo *repositories.SessionRepository) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRepo: sessionRepo,
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
		session, err := m.sessionRepo.FindByToken(token)
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

