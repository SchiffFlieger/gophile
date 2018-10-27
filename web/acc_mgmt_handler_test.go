package web

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/SchiffFlieger/gophile/auth"
	"github.com/SchiffFlieger/gophile/session"
	"github.com/stretchr/testify/assert"
)

type TestRegisterHandler struct{}

func (h TestRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Register(w, r)
}

func TestRegisterGet(t *testing.T) {
	s := httptest.NewServer(&TestRegisterHandler{})
	defer s.Close()

	resp, _ := http.Get(s.URL)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestRegisterPost(t *testing.T) {
	auth.GlobalAccountManager = &TestAuthenticator{}
	GlobalTemplateManager = NewTemplateManager()

	t.Run("ok", registerOk)
	t.Run("not ok", registerNotOk)

	auth.GlobalAccountManager.DeleteFile()
}

func registerOk(t *testing.T) {
	s := httptest.NewServer(&TestRegisterHandler{})
	defer s.Close()

	f := url.Values{}
	f.Add("user", "gooduser")
	f.Add("password", "123")
	f.Add("passwordMatch", "123")
	resp, _ := http.PostForm(s.URL, f)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func registerNotOk(t *testing.T) {
	s := httptest.NewServer(&TestRegisterHandler{})
	defer s.Close()

	f := url.Values{}
	f.Add("user", "baduser")
	f.Add("password", "123")
	f.Add("passwordMatch", "123")
	resp, _ := http.PostForm(s.URL, f)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

type TestLoginHandler struct{}

func (h TestLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Login(w, r)
}

func TestLoginGet(t *testing.T) {
	s := httptest.NewServer(&TestLoginHandler{})
	defer s.Close()

	resp, _ := http.Get(s.URL)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestLoginPost(t *testing.T) {
	auth.GlobalAccountManager = &TestAuthenticator{}
	session.GlobalSessionManager = &TestSessionManagerWithSession{}
	GlobalTemplateManager = NewTemplateManager()

	t.Run("ok", loginOk)
	t.Run("not ok", loginNotOk)

	auth.GlobalAccountManager.DeleteFile()
}

func loginOk(t *testing.T) {
	s := httptest.NewServer(&TestLoginHandler{})
	defer s.Close()

	f := url.Values{}
	f.Add("user", "gooduser")
	f.Add("password", "123")
	resp, _ := http.PostForm(s.URL, f)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func loginNotOk(t *testing.T) {
	s := httptest.NewServer(&TestLoginHandler{})
	defer s.Close()

	f := url.Values{}
	f.Add("user", "gooduser")
	f.Add("password", "1234")
	resp, _ := http.PostForm(s.URL, f)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

type TestLogoutHandler struct{}

func (h TestLogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), SID_CONTEXT, "test-session-id")
	Logout(w, r.WithContext(ctx))
}

func TestLogoutGet(t *testing.T) {
	s := httptest.NewServer(&TestLogoutHandler{})
	defer s.Close()

	resp, _ := http.Get(s.URL)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestLogoutPost(t *testing.T) {
	session.GlobalSessionManager = &TestSessionManagerWithSession{}
	GlobalTemplateManager = NewTemplateManager()
	s := httptest.NewServer(&TestLogoutHandler{})
	defer s.Close()

	resp, _ := http.PostForm(s.URL, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

type TestChangePasswordHandler struct{}

func (h TestChangePasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), USER_CONTEXT, "gooduser")
	ctx = context.WithValue(ctx, PATH_CONTEXT, "")
	ChangePassword(w, r.WithContext(ctx))
}

func TestChangePasswordMethodNotAllowed(t *testing.T) {
	session.GlobalSessionManager = &TestSessionManagerWithSession{}
	GlobalTemplateManager = NewTemplateManager()
	s := httptest.NewServer(&TestChangePasswordHandler{})
	defer s.Close()

	resp, _ := http.Get(s.URL)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestChangePasswordOk(t *testing.T) {
	auth.GlobalAccountManager = &TestAuthenticator{}
	session.GlobalSessionManager = &TestSessionManagerWithSession{}
	GlobalTemplateManager = NewTemplateManager()
	s := httptest.NewServer(&TestChangePasswordHandler{})
	defer s.Close()

	f := url.Values{}
	f.Set("currentPass", "123")
	f.Set("newPass", "abc")
	resp, _ := http.PostForm(s.URL, f)
	data, _ := ioutil.ReadAll(resp.Body)
	html := string(data)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, html, "Password change successful.")
}

func TestChangePasswordNotOk(t *testing.T) {
	auth.GlobalAccountManager = &TestAuthenticator{}
	session.GlobalSessionManager = &TestSessionManagerWithSession{}
	GlobalTemplateManager = NewTemplateManager()
	s := httptest.NewServer(&TestChangePasswordHandler{})
	defer s.Close()

	f := url.Values{}
	f.Set("currentPass", "abc")
	f.Set("newPass", "def")
	resp, _ := http.PostForm(s.URL, f)
	data, _ := ioutil.ReadAll(resp.Body)
	html := string(data)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, html, "Wrong password")
}
