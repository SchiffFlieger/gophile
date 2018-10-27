package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/SchiffFlieger/gophile/auth"
	"github.com/SchiffFlieger/gophile/session"
	"github.com/stretchr/testify/assert"
)

type TestSessionWrapper struct {
	name, value string
}

func (tsw TestSessionWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.AddCookie(&http.Cookie{
		Name:  tsw.name,
		Value: tsw.value,
	})
	WithSession(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "teststring")
	})(w, r)
}

func TestWithSessionOk(t *testing.T) {
	session.GlobalSessionManager = &TestSessionManagerWithSession{}
	GlobalTemplateManager = NewTemplateManager()

	s := httptest.NewServer(&TestSessionWrapper{name: "test-cookie", value: "test-session-ok"})
	defer s.Close()

	resp, _ := http.Get(s.URL)
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "teststring", string(data))
}

func TestWithMissingCookie(t *testing.T) {
	session.GlobalSessionManager = &TestSessionManagerWithSession{}
	GlobalTemplateManager = NewTemplateManager()

	s := httptest.NewServer(&TestSessionWrapper{name: "wrong-cookie", value: "test-session-not-ok"})
	defer s.Close()

	resp, _ := http.Get(s.URL)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestWithMissingSession(t *testing.T) {
	session.GlobalSessionManager = &TestSessionManagerWithSession{}
	GlobalTemplateManager = NewTemplateManager()

	s := httptest.NewServer(&TestSessionWrapper{name: "test-cookie", value: "test-session-not-ok"})
	defer s.Close()

	resp, _ := http.Get(s.URL)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestWithSessionUnauthorized(t *testing.T) {
	session.GlobalSessionManager = &TestSessionManagerWithSession{}
	GlobalTemplateManager = NewTemplateManager()

	s := httptest.NewServer(&TestSessionWrapper{name: "test-cookie", value: "test-session-unauthorized"})
	defer s.Close()

	resp, _ := http.Get(s.URL)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

type TestDataWrapperStruct struct {
	h http.HandlerFunc
}

func (tsw TestDataWrapperStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.AddCookie(&http.Cookie{
		Name:  "test-cookie",
		Value: "session-ok",
	})
	WithData(tsw.h)(w, r)
}

func TestDataWrapper(t *testing.T) {
	session.GlobalSessionManager = &TestSessionManagerWithData{}

	s := httptest.NewServer(&TestDataWrapperStruct{h: func(w http.ResponseWriter, r *http.Request) {
		res := r.Context().Value(USER_CONTEXT).(string) + "\n"
		res += r.Context().Value(PATH_CONTEXT).(string) + "\n"
		res += r.Context().Value(SID_CONTEXT).(string)
		fmt.Fprint(w, res)
	}})
	defer s.Close()

	resp, _ := http.Get(s.URL)
	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "testuser\ntestpath\nsession-ok", string(data))
}

type TestBasicAuthStruct struct {
	user, pass string
}

func (tba TestBasicAuthStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.SetBasicAuth(tba.user, tba.pass)
	BasicAuth(w, r)
}

func TestBasicAuthOk(t *testing.T) {
	auth.GlobalAccountManager = &TestAuthenticator{}
	p := path.Join(auth.USERDIR, "gooduser", "testpath")
	os.MkdirAll(p, os.ModePerm)
	f, _ := os.Create(path.Join(p, "testfile"))
	f.Close()

	s := httptest.NewServer(&TestBasicAuthStruct{user: "gooduser", pass: "123"})
	defer s.Close()

	resp, _ := http.Get(s.URL + "/basicauth/testpath/testfile")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestBasicAuthUnauthorized(t *testing.T) {
	auth.GlobalAccountManager = &TestAuthenticator{}
	defer os.RemoveAll(auth.USERDIR)
	p := path.Join(auth.USERDIR, "testuser", "testpath")
	os.MkdirAll(p, os.ModePerm)
	f, _ := os.Create(path.Join(p, "testfile"))
	defer f.Close()

	s := httptest.NewServer(&TestBasicAuthStruct{user: "gooduser", pass: "abc"})
	defer s.Close()

	resp, _ := http.Get(s.URL + "/basicauth/testpath/testfile")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestBasicAuthFileNotFound(t *testing.T) {
	auth.GlobalAccountManager = &TestAuthenticator{}
	s := httptest.NewServer(&TestBasicAuthStruct{user: "gooduser", pass: "123"})
	defer s.Close()

	resp, _ := http.Get(s.URL + "/basicauth/testpath/testfile")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
