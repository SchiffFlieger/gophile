package web

import (
	"flag"
	"log"
	"path"
	"strconv"
	"time"

	"github.com/SchiffFlieger/gophile/auth"
	"github.com/SchiffFlieger/gophile/session"
)

var PortString = ":8080"

// Initializes all variables and flags.
func Init() {
	var portFlag = flag.Int("port", 9090, "Defines the Port the server is run on.")
	var sessionLifetimeString = flag.String("sessionLifetime", "15m", "Defines the time in minutes until a user is automatically being logged out")
	flag.Parse()

	if *portFlag < 0 || *portFlag > 65536 {
		log.Fatal("Invalid Port")
	}

	sessionLifetime, timeErr := time.ParseDuration(*sessionLifetimeString)
	if timeErr != nil {
		log.Fatal("Invalid Session Lifetime format")
	}
	PortString = ":" + strconv.Itoa(*portFlag)
	session.GlobalSessionManager = session.NewSessionManager("gophile-session", sessionLifetime, 64)
	auth.GlobalAccountManager = auth.NewAccountManager(path.Join(auth.USERDIR, "accounts.xml"), 16)
	GlobalTemplateManager = NewTemplateManager()
}
