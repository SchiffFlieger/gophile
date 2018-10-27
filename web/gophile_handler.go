package web

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/SchiffFlieger/gophile/auth"
	"github.com/SchiffFlieger/gophile/session"
)

// Handles the route /gophile. Passes the different form actions to the corresponding handler.
func GophilePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		refresh(w, r)
	} else {
		r.ParseForm()
		r.ParseMultipartForm(16000000) // 16 MB
		if r.Form["Back"] != nil {
			goToParentDir(w, r)
		} else if r.Form["Refresh"] != nil {
			refresh(w, r)
		} else if r.Form["FolderName"] != nil {
			createFolder(w, r)
		} else if r.Form["BrowseFolder"] != nil {
			browseFolder(w, r)
		} else if r.Form["DownloadFile"] != nil {
			downloadFile(w, r)
			refresh(w, r)
		} else if r.Form["LinkFile"] != nil {
			printLink(w, r)
		} else if r.Form["DeleteFile"] != nil {
			deleteFile(w, r)
		} else {
			uploadFile(w, r)
		}
	}
}

// Creates a url for downloading files via basic auth.
func printLink(w http.ResponseWriter, r *http.Request) {
	dir := r.Form.Get("LinkFile")
	p := r.Context().Value(PATH_CONTEXT).(string)

	if p != "" {
		p = p + "/"
	}
	if strings.Contains(p, "./") {
		p = strings.Replace(p, "./", "", -1)
	}

	link := "https://" + r.Host + "/basicauth/" + p + dir
	ctx := context.WithValue(r.Context(), ACTION_CONTEXT, "BasicAuth-Link: "+link)

	refresh(w, r.WithContext(ctx))
}

// Direct download of a file.
func downloadFile(w http.ResponseWriter, r *http.Request) {
	dir := r.Form.Get("DownloadFile")
	u := r.Context().Value(USER_CONTEXT).(string)
	p := r.Context().Value(PATH_CONTEXT).(string)
	filepath := path.Join(auth.USERDIR, u, p, dir)

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/file")
	w.Header().Set("Content-Disposition", "attachment; filename="+dir)
	http.ServeContent(w, r, filepath, time.Now(), bytes.NewReader(data))
}

// Deletes a file or a folder with all its contents.
func deleteFile(w http.ResponseWriter, r *http.Request) {
	dir := r.Form.Get("DeleteFile")
	//defer refresh(w, r, "File/Folder " +  + " deleted.")
	u := r.Context().Value(USER_CONTEXT).(string)
	p := r.Context().Value(PATH_CONTEXT).(string)

	toDelete := path.Join(auth.USERDIR, u, p, dir)
	err := os.RemoveAll(toDelete)
	if err != nil {
		panic(err)
	}
	ctx := context.WithValue(r.Context(), ACTION_CONTEXT, "File/Folder "+dir+" deleted.")
	refresh(w, r.WithContext(ctx))
}

// Receives an uploaded file and stores it on the hard drive.
func uploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, _ := r.FormFile("UploadFile")
	defer file.Close()
	actionText := "File " + handler.Filename + " uploaded."

	u := r.Context().Value(USER_CONTEXT).(string)
	p := r.Context().Value(PATH_CONTEXT).(string)

	toWrite := path.Join(auth.USERDIR, u, p, handler.Filename)

	if _, err := os.Stat(toWrite); err == nil {
		actionText = "File already exists"
	} else {
		f, err := os.OpenFile(toWrite, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			ctx := context.WithValue(r.Context(), ACTION_CONTEXT, "Cannot upload file")
			refresh(w, r.WithContext(ctx))
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}

	ctx := context.WithValue(r.Context(), ACTION_CONTEXT, actionText)
	refresh(w, r.WithContext(ctx))
}

// Refreshes the current page.
func refresh(w http.ResponseWriter, r *http.Request) {
	var t string
	u := r.Context().Value(USER_CONTEXT).(string)
	p := r.Context().Value(PATH_CONTEXT).(string)
	val := r.Context().Value(ACTION_CONTEXT)
	if val == nil {
		t = ""
	} else {
		t = val.(string)
	}
	ShowGophile(w, CreateFileData(u, p, t))
}

// Creates a new folder with the given name.
func createFolder(w http.ResponseWriter, r *http.Request) {
	dir := r.Form.Get("FolderName")
	u := r.Context().Value(USER_CONTEXT).(string)
	p := r.Context().Value(PATH_CONTEXT).(string)
	err := os.Mkdir(path.Join(auth.USERDIR, u, p, dir), os.ModePerm)
	var ctx context.Context
	if err != nil {
		ctx = context.WithValue(r.Context(), ACTION_CONTEXT, "Cannot create folder.")
	} else {
		ctx = context.WithValue(r.Context(), ACTION_CONTEXT, "Folder "+dir+" created.")
	}

	refresh(w, r.WithContext(ctx))
}

// Navigates to the parent directory of the current folder.
func goToParentDir(w http.ResponseWriter, r *http.Request) {
	p := path.Dir(r.Context().Value(PATH_CONTEXT).(string))
	u := r.Context().Value(USER_CONTEXT).(string)
	r.Context().Value(SESSION_CONTEXT).(*session.Session).Set(session.SESSION_PATH, p)
	ShowGophile(w, CreateFileData(u, p, ""))
}

// Navigates inside the chosen folder.
func browseFolder(w http.ResponseWriter, r *http.Request) {
	dir := r.Form.Get("BrowseFolder")
	p := path.Join(r.Context().Value(PATH_CONTEXT).(string), dir)
	u := r.Context().Value(USER_CONTEXT).(string)
	r.Context().Value(SESSION_CONTEXT).(*session.Session).Set(session.SESSION_PATH, p)
	ShowGophile(w, CreateFileData(u, p, ""))
}
