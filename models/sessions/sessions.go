package sessions

import (
	"errors"
	"sync"
	"time"
)

type session struct {
	userID         string
	token          string
	expirationTime time.Time
}

var sessions = make(map[string]*session)

var mutex = &sync.Mutex{}

var (
	// ErrAlreadyExists session is already existing
	ErrAlreadyExists = errors.New("sessions: already exists")
	// ErrNotFound session is not found
	ErrNotFound = errors.New("sessions: not found")
	// ErrNotImplemented not implemented
	ErrNotImplemented = errors.New("sessions: not implemented")
)

// New creates new session
func New(userID string, token string, expirationTime time.Time) error {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := sessions[token]; ok {
		return ErrAlreadyExists
	}
	s := new(session)
	s.token = token
	s.userID = userID
	s.expirationTime = expirationTime
	sessions[token] = s
	return nil
}

// Refresh sets expirationTime
func Refresh(token string, expirationTime time.Time) error {
	mutex.Lock()
	defer mutex.Unlock()

	s, ok := sessions[token]
	if !ok {
		return ErrNotFound
	}

	s.expirationTime = expirationTime
	return nil
}

// Delete deletes session by token
func Delete(token string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := sessions[token]; !ok {
		return ErrNotFound
	}

	delete(sessions, token)
	return nil
}

// IsExpired checks is session expired by token
func IsExpired(token string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	s, ok := sessions[token]
	if !ok {
		return true
	}

	t := time.Now()
	return t.After(s.expirationTime)
}

// GetUser gets user id by token
func GetUser(token string) (userID string, e error) {
	s, ok := sessions[token]
	if !ok {
		e = ErrNotFound
		return
	}
	userID = s.userID
	return
}

// Logout removes all sessions for user
func Logout(userID string) (e error) {
	e = nil

	for k, s := range sessions {
		if s.userID == userID {
			delete(sessions, k)
		}
	}
	return
}
