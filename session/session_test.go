package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var dur time.Duration

func TestRunAllSessionTests(t *testing.T) {
	dur, _ = time.ParseDuration("15m")

	t.Run("store and load value", storeAndLoadValue)
	t.Run("store and delete value", storeAndDeleteValue)
	t.Run("session inactive", sessionInactive)
}

func storeAndLoadValue(t *testing.T) {
	sess := NewSession()
	sess.Set("testkey", true)

	assert.Equal(t, true, sess.Get("testkey"))
	assert.Equal(t, false, sess.Expired(dur))
}

func storeAndDeleteValue(t *testing.T) {
	sess := NewSession()

	sess.Set("testkey", true)
	sess.Delete("testkey")

	assert.Equal(t, nil, sess.Get("testkey"))
	assert.Equal(t, false, sess.Expired(dur))
}

func sessionInactive(t *testing.T) {
	dur, _ := time.ParseDuration("1s")
	sess := NewSession()

	time.Sleep(dur)
	assert.Equal(t, true, sess.Expired(dur))
}
