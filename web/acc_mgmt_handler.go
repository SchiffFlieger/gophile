package web

import (
	"encoding/base64"
	"net/http"
	"path"

	"github.com/SchiffFlieger/gophile/auth"
	"github.com/SchiffFlieger/gophile/session"
)

// Handles the route /login. Checks the credentials and creates a new session as needed.
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		user := r.Form.Get("user")
		password := r.Form.Get("password")
		if auth.GlobalAccountManager.Authenticate(user, password) {
			sess, sid := session.GlobalSessionManager.CreateSession()
			sess.Set(session.SESSION_USERNAME, user)
			sess.Set(session.SESSION_AUTHENTICATED, true)
			sess.Set(session.SESSION_PATH, path.Join())
			http.SetCookie(w, NewCookie(sid))
			w.WriteHeader(http.StatusOK)
			ShowGophile(w, CreateFileData(user, "", "Welcome, "+user+"!"))
		} else {
			w.WriteHeader(http.StatusForbidden)
			ShowIndex(w, "User or password incorrect.")
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Handles the route /register. Creates a new user account.
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		user := r.Form.Get("user")
		password := r.Form.Get("password")
		passwordMatch := r.Form.Get("passwordMatch")

		if password != passwordMatch {
			w.WriteHeader(http.StatusBadRequest)
			ShowIndex(w, "Passwords do not match.")
			return
		}

		// easter egg
		if base64.StdEncoding.EncodeToString([]byte(user)) == "cmljaw==" {
			url, _ := base64.StdEncoding.DecodeString("aHR0cDovL3JpY2tyb2xsb21hdGljLmNvbQ==")
			rq, _ := http.NewRequest("get", string(url), nil)
			http.Redirect(w, rq, string(url), http.StatusMovedPermanently)
		}
		err := auth.GlobalAccountManager.CreateAccount(user, password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ShowIndex(w, err.Error())
		} else {
			w.WriteHeader(http.StatusOK)
			ShowIndex(w, "User "+user+" created.")
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Handles the route /logout. Deletes the active user session.
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		sid := r.Context().Value(SID_CONTEXT).(string)
		session.GlobalSessionManager.DropSession(sid)
		w.WriteHeader(http.StatusOK)
		ShowIndex(w, "User logged out.")
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Handles the route /changePassword. Sets a new password for an existing user.
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		cp := r.Form.Get("currentPass")
		np := r.Form.Get("newPass")
		u := r.Context().Value(USER_CONTEXT).(string)
		p := r.Context().Value(PATH_CONTEXT).(string)
		err := auth.GlobalAccountManager.ChangePassword(u, cp, np)
		w.WriteHeader(http.StatusOK)
		if err != nil {
			ShowIndex(w, "Wrong password")
		} else {
			ShowGophile(w, CreateFileData(u, p, "Password change successful."))
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
