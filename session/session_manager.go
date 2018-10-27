package session

import (
	"sync"
	"time"

	"github.com/SchiffFlieger/gophile/random"
)

var GlobalSessionManager SessionManagerInterface

type SessionManagerInterface interface {
	CreateSession() (*Session, string)
	GetSession(sid string) *Session
	DropSession(sid string)
	CookieName() string
	SessionLifetime() time.Duration
}

type SessionManager struct {
	cookieName string
	lock       sync.Mutex
	lifetime   time.Duration
	sidLength  int
	sessions   map[string]*Session
}

// Creates a new session manager. The session manager has an automatic garbage collector for expired sessions.
// Name is the name of the cookie that is set in the browser. Lifetime is the lifetime for the cookie. SidLength
// is the length of the session id.
func NewSessionManager(name string, lifetime time.Duration, sidLength int) *SessionManager {
	sm := SessionManager{cookieName: name,
		lifetime:  lifetime,
		sessions:  make(map[string]*Session),
		sidLength: sidLength}
	go sm.dropExpiredSessions()
	return &sm
}

// Creates and returns a new session and its session id.
func (sm *SessionManager) CreateSession() (*Session, string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	sid := random.Rnd.RandomString(sm.sidLength)
	sm.sessions[sid] = NewSession()
	return sm.sessions[sid], sid
}

// Returns the session corresponding to the given session id.
func (sm *SessionManager) GetSession(sid string) *Session {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	return sm.sessions[sid]
}

// Deletes an active session.
func (sm *SessionManager) DropSession(sid string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	delete(sm.sessions, sid)
}

func (sm *SessionManager) CookieName() string {
	return sm.cookieName
}

func (sm *SessionManager) SessionLifetime() time.Duration {
	return sm.lifetime
}

// This is the garbage collector for expired sessions. This method calls itself after a
// certain amount of time.
func (sm *SessionManager) dropExpiredSessions() {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	for sid, sess := range sm.sessions {
		if sess.Expired(sm.lifetime) {
			delete(sm.sessions, sid)
		}
	}

	time.AfterFunc(sm.lifetime, sm.dropExpiredSessions)
}
