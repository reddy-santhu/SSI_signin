package services

import (
	"sync"
	"time"
)

type ProofRequestStore struct {
	mu      sync.RWMutex
	store   map[string]string
	expiry  map[string]time.Time
}

var globalStore = &ProofRequestStore{
	store:  make(map[string]string),
	expiry: make(map[string]time.Time),
}

func NewProofRequestStore() *ProofRequestStore {
	return globalStore
}

func (s *ProofRequestStore) Set(proofRequestID, sessionToken string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.store[proofRequestID] = sessionToken
	s.expiry[proofRequestID] = time.Now().Add(5 * time.Minute)
	
	go s.cleanup()
}

func (s *ProofRequestStore) Get(proofRequestID string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if expiry, exists := s.expiry[proofRequestID]; exists {
		if time.Now().After(expiry) {
			delete(s.store, proofRequestID)
			delete(s.expiry, proofRequestID)
			return "", false
		}
		token, ok := s.store[proofRequestID]
		return token, ok
	}
	
	return "", false
}

func (s *ProofRequestStore) Exists(proofRequestID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if expiry, exists := s.expiry[proofRequestID]; exists {
		if time.Now().After(expiry) {
			return false
		}
		return true
	}
	
	return false
}

func (s *ProofRequestStore) Delete(proofRequestID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	delete(s.store, proofRequestID)
	delete(s.expiry, proofRequestID)
}

func (s *ProofRequestStore) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	now := time.Now()
	for id, exp := range s.expiry {
		if now.After(exp) {
			delete(s.store, id)
			delete(s.expiry, id)
		}
	}
}


