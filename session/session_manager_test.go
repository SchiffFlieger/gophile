package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var manager SessionManagerInterface

func TestRunAllSessionManagerTests(t *testing.T) {
	timeout, _ := time.ParseDuration("15m")
	manager = NewSessionManager("testcookie", timeout, 64)

	t.Run("create and verify session", createAndVerifySession)
	t.Run("create and drop session", createAndDropSession)
	t.Run("check cookie name", checkCookieName)
	t.Run("check session lifetime", checkLifetime)
}

func createAndVerifySession(t *testing.T) {
	session1, sid := manager.CreateSession()
	session2 := manager.GetSession(sid)

	assert.Equal(t, session1, session2)
}

func createAndDropSession(t *testing.T) {
	_, sid := manager.CreateSession()
	manager.DropSession(sid)

	res := manager.GetSession(sid)
	if res == nil {
		assert.Equal(t, true, true)
	} else {
		assert.Equal(t, true, false)
	}
}

func checkCookieName(t *testing.T) {
	assert.Equal(t, "testcookie", manager.CookieName())
}

func checkLifetime(t *testing.T) {
	dur, _ := time.ParseDuration("15m")
	assert.Equal(t, dur, manager.SessionLifetime())
}

func TestCreateSessionWithTimeout(t *testing.T) {
	timeout, _ := time.ParseDuration("1s")
	sleeptime, _ := time.ParseDuration("2s")
	local_manager := NewSessionManager("testcookie", timeout, 64)
	_, sid := local_manager.CreateSession()
	time.Sleep(sleeptime)

	res := local_manager.GetSession(sid)
	if res == nil {
		assert.Equal(t, true, true)
	} else {
		assert.Equal(t, true, false)
	}
}
