package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
)

var filename = "common.js"
var port = 8080

var file_content = ""
var PrevTime time.Time = time.Now().AddDate(0, 0, -1)

func loadFileContent(filename string) {
	t := time.Now()
	if t.Sub(PrevTime).Seconds() > 10 {
		PrevTime = t
		buf, _ := ioutil.ReadFile(filename)
		file_content = string(buf)
		log.Info(fmt.Sprintf("File content updated. Size: %d", len(file_content)))
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	loadFileContent(filename)

	log.Info(fmt.Sprintf("Request served. URL: %s", r.URL.String()))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Age", "444107")
	w.Header().Set("Cache-Control", "public,max-age=86313600")
	w.Header().Set("Content-Type", "application/x-javascript")
	w.Header().Set("Date", "Mon, 26 Sep 2016 10:53:18 GMT")
	w.Header().Set("Expires", "Mon, 17 Jun 2019 06:53:31 GMT")
	w.Header().Set("Last-Modified", "Tue, 20 Sep 2016 18:58:18 GMT")
	w.Header().Set("Server", "Apache")
	w.Header().Set("Vary", "Accept-Encoding")
	fmt.Fprintf(w, "%s", file_content)
}

func main() {
	log.Info(fmt.Sprintf("Server started. port: %d", port))

	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
