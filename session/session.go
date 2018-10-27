package session

import (
	"time"
)

type Session struct {
	lastAccess time.Time
	values     map[string]interface{}
}

const SESSION_USERNAME = "username"
const SESSION_AUTHENTICATED = "auth"
const SESSION_PATH = "currentpath"

// Creates a new session. The session contains the last access time and a key-value map.
func NewSession() *Session {
	return &Session{lastAccess: time.Now(), values: make(map[string]interface{})}
}

// Sets a new key-value pair in the session.
func (s *Session) Set(key string, value interface{}) {
	defer s.refresh()
	s.values[key] = value
}

// Returns the stored value corresponding to the given key.
func (s *Session) Get(key string) interface{} {
	defer s.refresh()
	return s.values[key]
}

// Deletes a key-value pair from the session.
func (s *Session) Delete(key string) {
	defer s.refresh()
	delete(s.values, key)
}

// Checks if the session is inactive for the given duration.
func (s *Session) Expired(dur time.Duration) bool {
	return time.Now().After(s.lastAccess.Add(dur))
}

// Refreshes the last access time of the session.
func (s *Session) refresh() {
	s.lastAccess = time.Now()
}
