package web

import (
	"net/http"
)

// Handles the route /. Shows a default landing page.
func LandingPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		ShowIndex(w, "")
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// Handles the route /impressum.
func Impressum(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		ShowImpressum(w)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
