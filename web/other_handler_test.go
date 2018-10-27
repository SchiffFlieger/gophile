package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestLandingPageHandler struct{}

func (h TestLandingPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	LandingPage(w, r)
}

type TestImpressumHandler struct{}

func (h TestImpressumHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Impressum(w, r)
}

func TestLandingPageGet(t *testing.T) {
	s := httptest.NewServer(&TestLandingPageHandler{})
	defer s.Close()

	resp, _ := http.Get(s.URL)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestLandingPagePost(t *testing.T) {
	s := httptest.NewServer(&TestLandingPageHandler{})
	defer s.Close()

	resp, _ := http.PostForm(s.URL, nil)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestImpressumGet(t *testing.T) {
	s := httptest.NewServer(&TestImpressumHandler{})
	defer s.Close()

	resp, _ := http.Get(s.URL)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestImpressumPost(t *testing.T) {
	s := httptest.NewServer(&TestImpressumHandler{})
	defer s.Close()

	resp, _ := http.PostForm(s.URL, nil)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}
