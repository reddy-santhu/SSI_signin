package services

import (
	"strings"
	"testing"
	"time"
)

func TestSessionService_GenerateToken(t *testing.T) {
	s := NewSessionService()
	tok, err := s.GenerateToken()
	if err != nil {
		t.Fatal(err)
	}
	if tok == "" {
		t.Fatal("empty token")
	}
	if strings.ContainsAny(tok, "+/") {
		t.Fatalf("token should be URL-safe base64: %q", tok)
	}
	if len(tok) < 40 {
		t.Fatalf("expected long token, got len %d", len(tok))
	}
}

func TestSessionService_GetExpirationTime(t *testing.T) {
	s := NewSessionService()
	before := time.Now()
	exp := s.GetExpirationTime()
	if !exp.After(before.Add(23 * time.Hour)) {
		t.Fatalf("expiration too soon: %v", exp)
	}
	if !exp.Before(before.Add(25 * time.Hour)) {
		t.Fatalf("expiration too far: %v", exp)
	}
}
