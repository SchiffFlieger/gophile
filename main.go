package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SchiffFlieger/gophile/web"
)

func main() {
	web.Init()

	http.HandleFunc("/basicauth/", web.BasicAuth)
	http.HandleFunc("/login", web.Login)
	http.HandleFunc("/register", web.Register)
	http.HandleFunc("/impressum", web.Impressum)
	http.HandleFunc("/gophile", web.WithSession(web.WithData(web.GophilePage)))
	http.HandleFunc("/logout", web.WithSession(web.WithData(web.Logout)))
	http.HandleFunc("/changePassword", web.WithSession(web.WithData(web.ChangePassword)))
	http.HandleFunc("/", web.LandingPage)

	fmt.Println("URL: http://localhost" + web.PortString)

	err := http.ListenAndServe(web.PortString, nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}
