// www is an HTTP server that runs on port 80 and serves file from $HOME/www.
// If $HOME isn't set, it defaults to /www.
//
// Instructions:
//
//   1) copy the files you want to serve to ~/www:
//   $ mkdir ~/www
//   $ cp *.css *.js ~/www
//
//   2) build the server:
//   $ go build
//
//   3) Give the program Linux capabilities to bind low ports:
//
//   $ sudo setcap 'cap_net_bind_service=+ep' www
//
//   4) Run:
//   $ ./www
//   2013/03/29 15:33:57 Serving files from /www

package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nictuku/mothership/login"
)

func RequireAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		passport, err := login.CurrentPassport(req)
		if err != nil {
			log.Printf("Redirecting to ghlogin: %q. Referrer: %q", err, req.Referer())
			http.Redirect(w, req, "/ghlogin", http.StatusFound)
			return
		}
		if passport.Email == "yves.junqueira@gmail.com" {
			handler.ServeHTTP(w, req)
		} else {
			http.Error(w, "Nope.", http.StatusForbidden)
		}

	})
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func homeDir() string {
	home := os.Getenv("HOME")
	if home == "" {
		home = string(filepath.Separator)
	}
	return filepath.Join(home, "www")
}

func main() {
	var dir string
	switch len(os.Args) {
	case 1:
		dir = homeDir()
	case 2:
		dir = os.Args[1]
	default:
		log.Fatalln("Too many arguments")
		log.Fatalln("Usage: %v [dir]", os.Args[0])
	}
	dir = filepath.Clean(dir)
	log.Printf("Serving files from %v", dir)
	http.Handle("/", RequireAuth(Log(http.FileServer(http.Dir(dir)))))
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Println("Error starting www server:", err)
		// os.IsPermission doesn't match.
		if strings.Contains(err.Error(), "permission denied") {
			log.Println("Try: sudo setcap 'cap_net_bind_service=+ep' www")
		}
	}
}
