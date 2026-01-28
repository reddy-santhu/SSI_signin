package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"ssi-signin/backend/models"
)

type stubSessions struct {
	session *models.Session
	err     error
}

func (s *stubSessions) FindByToken(token string) (*models.Session, error) {
	if s.err != nil {
		return nil, s.err
	}
	if token == "good-token" {
		return s.session, nil
	}
	return nil, nil
}

func TestRequireAuth_MissingHeader(t *testing.T) {
	e := echo.New()
	mw := NewAuthMiddleware(&stubSessions{})
	h := mw.RequireAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := h(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("code=%d", rec.Code)
	}
}

func TestRequireAuth_InvalidBearerFormat(t *testing.T) {
	e := echo.New()
	mw := NewAuthMiddleware(&stubSessions{})
	h := mw.RequireAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "not-bearer x")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := h(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("code=%d", rec.Code)
	}
}

func TestRequireAuth_ExpiredSession(t *testing.T) {
	e := echo.New()
	sess := &models.Session{
		UserID:    1,
		Token:     "good-token",
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	mw := NewAuthMiddleware(&stubSessions{session: sess})
	h := mw.RequireAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer good-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := h(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("code=%d", rec.Code)
	}
}

func TestRequireAuth_RepoError(t *testing.T) {
	e := echo.New()
	mw := NewAuthMiddleware(&stubSessions{err: errors.New("db down")})
	h := mw.RequireAuth(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer good-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := h(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("code=%d", rec.Code)
	}
}

func TestRequireAuth_OK(t *testing.T) {
	e := echo.New()
	sess := &models.Session{
		UserID:    42,
		Token:     "good-token",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	mw := NewAuthMiddleware(&stubSessions{session: sess})
	var sawUserID int
	h := mw.RequireAuth(func(c echo.Context) error {
		sawUserID, _ = c.Get("user_id").(int)
		return c.String(http.StatusOK, "ok")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer good-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := h(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("code=%d", rec.Code)
	}
	if sawUserID != 42 {
		t.Fatalf("user_id=%d", sawUserID)
	}
}
