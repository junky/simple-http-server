package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/NYTimes/gziphandler"
	log "github.com/Sirupsen/logrus"
)

var port = 80

func handler(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Method + ` ` + r.URL.String())
	if r.Method == "POST" {
		file, _, err := r.FormFile("data")
		if err != nil {
			log.Info(err.Error())
		}

		if file != nil {
			defer file.Close()

			filename := fmt.Sprintf("./uploads/%d", time.Now().UnixNano())
			f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()
			io.Copy(f, file)

		} else {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			log.Info(buf.String())
		}

	}
	fmt.Fprintf(w, "falcon.aws - OK")
}

func main() {
	log.Info(fmt.Sprintf("Server started. port: %d", port))

	//	http.HandleFunc("/", handler)

	fs := http.FileServer(http.Dir("./static"))
	fsWithGz := gziphandler.GzipHandler(fs)
	http.Handle("/", fsWithGz)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
