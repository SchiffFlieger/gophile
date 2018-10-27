package web

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/SchiffFlieger/gophile/auth"
	"github.com/SchiffFlieger/gophile/session"

	"github.com/stretchr/testify/assert"
)

type TestGophilePageHandler struct{}

func (h TestGophilePageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), SID_CONTEXT, "test-session-id")
	ctx = context.WithValue(ctx, USER_CONTEXT, "testuser")
	ctx = context.WithValue(ctx, PATH_CONTEXT, "")
	ctx = context.WithValue(ctx, ACTION_CONTEXT, "TestAction")

	GophilePage(w, r.WithContext(ctx))
}

func TestGophilePageGet(t *testing.T) {
	if GlobalTemplateManager == nil {
		GlobalTemplateManager = NewTemplateManager()
	}
	s := httptest.NewServer(&TestGophilePageHandler{})
	defer s.Close()

	resp, _ := http.Get(s.URL)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func TestPrintLink(t *testing.T) {
	if GlobalTemplateManager == nil {
		GlobalTemplateManager = NewTemplateManager()
	}
	s := httptest.NewServer(&TestGophilePageHandler{})
	defer s.Close()
	form := url.Values{}
	form.Add("LinkFile", "testfile")

	resp, _ := http.PostForm(s.URL, form)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	html := string(data)

	assert.Contains(t, html, strings.Replace(s.URL, "http", "https", 1)+"/basicauth/testfile")
}

func TestDownloadFile(t *testing.T) {
	setUpUserFolder()
	if GlobalTemplateManager == nil {
		GlobalTemplateManager = NewTemplateManager()
	}
	s := httptest.NewServer(&TestGophilePageHandler{})
	defer s.Close()

	form := url.Values{}
	form.Add("DownloadFile", "staticTestFile.txt")

	resp, _ := http.PostForm(s.URL, form)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteFile(t *testing.T) {
	if GlobalTemplateManager == nil {
		GlobalTemplateManager = NewTemplateManager()
	}
	setUpUserFolder()
	defer cleanUpUserFolder()

	s := httptest.NewServer(&TestGophilePageHandler{})
	defer s.Close()

	_, err := os.Stat(path.Join(auth.USERDIR, "testuser", "staticTestFile.txt"))
	assert.False(t, os.IsNotExist(err))

	form := url.Values{}
	form.Add("DeleteFile", "staticTestFile.txt")

	resp, _ := http.PostForm(s.URL, form)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_, err = os.Stat(path.Join(auth.USERDIR, "testuser", "staticTestFile.txt"))
	assert.True(t, os.IsNotExist(err))
}

func TestRefresh(t *testing.T) {
	setUpUserFolder()
	defer cleanUpUserFolder()
	if GlobalTemplateManager == nil {
		GlobalTemplateManager = NewTemplateManager()
	}

	s := httptest.NewServer(&TestGophilePageHandler{})

	resp, _ := http.Get(s.URL)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)
	html1 := string(data)

	resp, _ = http.Get(s.URL)

	data, _ = ioutil.ReadAll(resp.Body)
	html2 := string(data)

	assert.Equal(t, html1, html2)
}

func TestCreateFolder(t *testing.T) {
	setUpUserFolder()
	defer cleanUpUserFolder()
	if GlobalTemplateManager == nil {
		GlobalTemplateManager = NewTemplateManager()
	}
	s := httptest.NewServer(&TestGophilePageHandler{})
	defer s.Close()

	form := url.Values{}
	form.Add("FolderName", "testfolder")

	resp, _ := http.PostForm(s.URL, form)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_, err := os.Stat(path.Join(auth.USERDIR, "testuser", "testfolder"))
	assert.False(t, os.IsNotExist(err))
}

type TestGoToParentDirHandler struct{}

func (h TestGoToParentDirHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), SID_CONTEXT, "test-session-id")
	ctx = context.WithValue(ctx, USER_CONTEXT, "testuser")
	ctx = context.WithValue(ctx, PATH_CONTEXT, "testfolder")
	ctx = context.WithValue(ctx, ACTION_CONTEXT, "TestAction")
	s := session.NewSession()
	s.Set(session.SESSION_PATH, "testfolder")
	ctx = context.WithValue(ctx, SESSION_CONTEXT, s)
	GophilePage(w, r.WithContext(ctx))
}

func TestGoToParentDir(t *testing.T) {
	setUpUserFolder()
	defer cleanUpUserFolder()
	if GlobalTemplateManager == nil {
		GlobalTemplateManager = NewTemplateManager()
	}
	s := httptest.NewServer(&TestGoToParentDirHandler{})

	resp, _ := http.Get(s.URL)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	html := string(data)
	assert.Contains(t, html, `<span id="folderPath">/testfolder</span>`)

	form := url.Values{}
	form.Add("Back", "true")
	resp, _ = http.PostForm(s.URL, form)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	html = string(data)
	assert.Contains(t, html, `<span id="folderPath">/.</span>`)
}

type TestBrowseFolderHandler struct{}

func (h TestBrowseFolderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), SID_CONTEXT, "test-session-id")
	ctx = context.WithValue(ctx, USER_CONTEXT, "testuser")
	ctx = context.WithValue(ctx, PATH_CONTEXT, ".")
	ctx = context.WithValue(ctx, ACTION_CONTEXT, "TestAction")
	s := session.NewSession()
	s.Set(session.SESSION_PATH, ".")
	ctx = context.WithValue(ctx, SESSION_CONTEXT, s)
	GophilePage(w, r.WithContext(ctx))
}

func TestBrowseFolder(t *testing.T) {
	setUpUserFolder()
	defer cleanUpUserFolder()
	if GlobalTemplateManager == nil {
		GlobalTemplateManager = NewTemplateManager()
	}
	s := httptest.NewServer(&TestBrowseFolderHandler{})

	resp, _ := http.Get(s.URL)

	data, _ := ioutil.ReadAll(resp.Body)
	html := string(data)
	assert.Contains(t, html, `<span id="folderPath">/.</span>`)

	form := url.Values{}
	form.Add("BrowseFolder", "testfolder")
	resp, _ = http.PostForm(s.URL, form)

	data, _ = ioutil.ReadAll(resp.Body)
	html = string(data)

	assert.Contains(t, html, `<span id="folderPath">/testfolder</span>`)
}

func setUpUserFolder() {
	os.MkdirAll(path.Join(auth.USERDIR, "testuser"), os.ModePerm)

	fileTextBA := []byte("This is the text in the file.")
	ioutil.WriteFile(path.Join(auth.USERDIR, "testuser", "staticTestFile.txt"), fileTextBA, os.ModePerm)
}

func cleanUpUserFolder() {
	os.RemoveAll(auth.USERDIR)
}
