package main

import (
	"bytes"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

var port = 8080

func handler(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Method + ` ` + r.URL.String())
	if r.Method == "POST" {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		log.Info(buf.String())
	}
	fmt.Fprintf(w, "OK")
}

func main() {
	log.Info(fmt.Sprintf("Server started. port: %d", port))

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
