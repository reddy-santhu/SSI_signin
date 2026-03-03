package services

import (
	"sync"
	"time"
)

type ProofRequestStore struct {
	mu     sync.RWMutex
	store  map[string]string
	expiry map[string]time.Time
}

var (
	globalStore = &ProofRequestStore{
		store:  make(map[string]string),
		expiry: make(map[string]time.Time),
	}
	proofStoreCleanupOnce sync.Once
)

func NewProofRequestStore() *ProofRequestStore {
	return globalStore
}

func (s *ProofRequestStore) Set(proofRequestID, sessionToken string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[proofRequestID] = sessionToken
	s.expiry[proofRequestID] = time.Now().Add(5 * time.Minute)

	proofStoreCleanupOnce.Do(func() {
		go s.runPeriodicCleanup()
	})
}

func (s *ProofRequestStore) Get(proofRequestID string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expiry, exists := s.expiry[proofRequestID]
	if !exists {
		return "", false
	}
	if time.Now().After(expiry) {
		return "", false
	}
	return s.store[proofRequestID], true
}

func (s *ProofRequestStore) Exists(proofRequestID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	expiry, exists := s.expiry[proofRequestID]
	if !exists {
		return false
	}
	return !time.Now().After(expiry)
}

func (s *ProofRequestStore) Delete(proofRequestID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.store, proofRequestID)
	delete(s.expiry, proofRequestID)
}

func (s *ProofRequestStore) runPeriodicCleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.cleanup()
	}
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
