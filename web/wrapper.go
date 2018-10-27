package web

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/SchiffFlieger/gophile/auth"
	"github.com/SchiffFlieger/gophile/session"
)

const USER_CONTEXT = "ctx_user"
const PATH_CONTEXT = "ctx_path"
const SESSION_CONTEXT = "ctx_session"
const SID_CONTEXT = "ctx_sid"
const ACTION_CONTEXT = "ctx_action_text"

// Packs a handler and secures it behind a session. Access to the handler method is only granted
// if the user has an active authenticated session. If there is no active session this wrapper will
// redirect to the landing page.
func WithSession(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, _ := r.Cookie(session.GlobalSessionManager.CookieName())
		if c == nil {
			// set new cookie and create session
			sess, sid := session.GlobalSessionManager.CreateSession()
			c := NewCookie(sid)
			sess.Set(session.SESSION_USERNAME, "")
			sess.Set(session.SESSION_AUTHENTICATED, false)
			sess.Set(session.SESSION_PATH, "")
			http.SetCookie(w, c)

			w.WriteHeader(http.StatusUnauthorized)
			ShowIndex(w, "You do not have permission to view this site.")
		} else {
			// refresh existing cookie
			sess := session.GlobalSessionManager.GetSession(c.Value)
			if sess == nil {
				w.WriteHeader(http.StatusUnauthorized)
				ShowIndex(w, "Your cookie is obsolete: You have been logged out.")
			} else {
				http.SetCookie(w, c)
				if sess.Get(session.SESSION_AUTHENTICATED).(bool) {
					handler(w, r)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					ShowIndex(w, "You do not have permission to view this site")
				}
			}
		}
	}
}

// Injects several values in the context of a handler before calling it.
func WithData(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sm := session.GlobalSessionManager
		c, _ := r.Cookie(sm.CookieName())
		s := sm.GetSession(c.Value)
		u := s.Get(session.SESSION_USERNAME).(string)
		p := s.Get(session.SESSION_PATH).(string)
		ctx := context.WithValue(r.Context(), USER_CONTEXT, u)
		ctx = context.WithValue(ctx, PATH_CONTEXT, p)
		ctx = context.WithValue(ctx, SESSION_CONTEXT, s)
		ctx = context.WithValue(ctx, SID_CONTEXT, c.Value)
		handler(w, r.WithContext(ctx))
	}
}

// Handles the route /basicauth. This route allows an automated download with basic authentication.
func BasicAuth(w http.ResponseWriter, r *http.Request) {
	user, password, ok := r.BasicAuth()
	if ok && auth.GlobalAccountManager.Authenticate(user, password) {
		url := strings.Split(r.URL.String(), "/")
		filepath := path.Join(auth.USERDIR, user)
		for _, elem := range url {
			filepath = path.Join(filepath, elem)
		}
		filepath = strings.Replace(filepath, "basicauth/", "", 1)

		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/file")
		w.Header().Set("Content-Disposition", "attachment; filename="+url[len(url)-1])
		http.ServeContent(w, r, filepath, time.Now(), bytes.NewReader(data))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

// Creates a new cookie for a session.
func NewCookie(sid string) *http.Cookie {
	return &http.Cookie{Name: session.GlobalSessionManager.CookieName(),
		Expires: time.Now().Add(session.GlobalSessionManager.SessionLifetime()),
		Value:   sid}
}
