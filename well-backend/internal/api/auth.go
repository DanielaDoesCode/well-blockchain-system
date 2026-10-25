package api

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type SessionEntry struct {
	User       string
	Token      string
	ExpiresAt  time.Time
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*SessionEntry
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*SessionEntry),
	}
}

func (s *SessionStore) Add(user, token string, expiresAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[token] = &SessionEntry{
		User:      user,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	log.Printf("session added: user=%s expires=%s", user, expiresAt.Format(time.RFC3339))
}

func (s *SessionStore) Revoke(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
	log.Printf("session revoked: token=%s", token)
}

func (s *SessionStore) IsValid(ctx context.Context, user, token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, ok := s.sessions[token]
	if !ok {
		return false
	}
	if entry.User != user {
		return false
	}
	if time.Now().After(entry.ExpiresAt) {
		return false
	}
	return true
}

type AuthContextKey string

const (
	AuthCtxKey AuthContextKey = "auth_user"
)

func ExtractUser(ctx context.Context) (string, error) {
	v := ctx.Value(AuthCtxKey)
	if v == nil {
		return "", fmt.Errorf("no auth in context")
	}
	user, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("invalid auth value")
	}
	return user, nil
}
