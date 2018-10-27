package web

import (
	"errors"
	"time"

	"github.com/SchiffFlieger/gophile/session"
)

type TestAuthenticator struct{}
type TestSessionManagerWithSession struct{}
type TestSessionManagerWithData struct{}

func (ta *TestAuthenticator) CreateAccount(user, password string) error {
	if user == "baduser" {
		return errors.New("error")
	} else {
		return nil
	}
}
func (ta *TestAuthenticator) DeleteAccount(user, password string) error { return nil }
func (ta *TestAuthenticator) ChangePassword(user, oldPw, newPw string) error {
	if user == "gooduser" && oldPw == "123" {
		return nil
	} else {
		return errors.New("error")
	}
}
func (ta *TestAuthenticator) Authenticate(user, password string) bool {
	if user == "gooduser" && password == "123" {
		return true
	} else {
		return false
	}
}
func (ta *TestAuthenticator) DeleteFile() {}

func (tsm *TestSessionManagerWithSession) CreateSession() (*session.Session, string) {
	return session.NewSession(), "test-session"
}
func (tsm *TestSessionManagerWithSession) GetSession(sid string) *session.Session {
	if sid == "test-session-not-ok" {
		return nil
	}

	s := session.NewSession()
	s.Set(session.SESSION_AUTHENTICATED, sid == "test-session-ok")
	return s
}
func (tsm *TestSessionManagerWithSession) DropSession(sid string)         {}
func (tsm *TestSessionManagerWithSession) CookieName() string             { return "test-cookie" }
func (tsm *TestSessionManagerWithSession) SessionLifetime() time.Duration { return time.Minute * 15 }

func (tsm *TestSessionManagerWithData) CreateSession() (*session.Session, string) {
	return session.NewSession(), "test-session"
}
func (tsm *TestSessionManagerWithData) GetSession(sid string) *session.Session {
	s := session.NewSession()
	s.Set(session.SESSION_USERNAME, "testuser")
	s.Set(session.SESSION_PATH, "testpath")
	return s
}
func (tsm *TestSessionManagerWithData) DropSession(sid string)         {}
func (tsm *TestSessionManagerWithData) CookieName() string             { return "test-cookie" }
func (tsm *TestSessionManagerWithData) SessionLifetime() time.Duration { return time.Minute * 15 }
