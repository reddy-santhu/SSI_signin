package services

import (
	"testing"
	"time"
)

func TestProofRequestStore_SetGetDelete(t *testing.T) {
	s := &ProofRequestStore{
		store:  make(map[string]string),
		expiry: make(map[string]time.Time),
	}
	s.Set("pr1", "")
	tok, ok := s.Get("pr1")
	if !ok || tok != "" {
		t.Fatalf("Get pending: ok=%v tok=%q", ok, tok)
	}
	s.Set("pr1", "session-token")
	tok, ok = s.Get("pr1")
	if !ok || tok != "session-token" {
		t.Fatalf("Get completed: ok=%v tok=%q", ok, tok)
	}
	s.Delete("pr1")
	_, ok = s.Get("pr1")
	if ok {
		t.Fatal("expected not found after delete")
	}
}

func TestProofRequestStore_Exists(t *testing.T) {
	s := &ProofRequestStore{
		store:  make(map[string]string),
		expiry: make(map[string]time.Time),
	}
	s.Set("x", "")
	if !s.Exists("x") {
		t.Fatal("expected exists")
	}
	s.Delete("x")
	if s.Exists("x") {
		t.Fatal("expected not exists")
	}
}
