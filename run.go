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
)

func main() {
	home := os.Getenv("HOME")
	if home == "" {
		home = string(filepath.Separator)
	}
	dir := filepath.Join(home, "www")
	log.Printf("Serving files from %v", dir)

	err := http.ListenAndServe(":80", http.FileServer(http.Dir(dir)))
	if err != nil {
		log.Println("Error starting www server:", err)
	}
}
