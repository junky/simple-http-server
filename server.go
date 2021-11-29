package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var port = 8080

func postHandler(w http.ResponseWriter, r *http.Request) {
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

var chunks = []string{}

func chunkHandler(w http.ResponseWriter, r *http.Request) {
	conn := GetConn(r)

	log.Infof("RemoteAddr: %s", conn.RemoteAddr().String())
	log.Infof("LocalAddr: %s", conn.LocalAddr().String())

	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("expected http.ResponseWriter to be an http.Flusher")
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	for _, s := range chunks {
		log.Info(fmt.Sprintf("S: %s", s))
		fmt.Fprint(w, s)
		flusher.Flush() // Trigger "chunked" encoding and send a chunk...
		time.Sleep(50 * time.Millisecond)
	}
}

func getChunks() []string {
	buf, _ := ioutil.ReadFile("./static/001.full.html")
	s := string(buf)
	return strings.Split(s, "<!--chunk-->")
}

type contextKey struct {
	key string
}

var ConnContextKey = &contextKey{"http-conn"}

func SaveConnInContext(ctx context.Context, c net.Conn) context.Context {
	return context.WithValue(ctx, ConnContextKey, c)
}
func GetConn(r *http.Request) net.Conn {
	return r.Context().Value(ConnContextKey).(net.Conn)
}

func main() {
	log.Info(fmt.Sprintf("Server started. port: %d", port))

	chunks = getChunks()
	http.HandleFunc("/", chunkHandler)

	//Handle gziped static content
	/*
		fs := http.FileServer(http.Dir("./static"))
		fsWithGz := gziphandler.GzipHandler(fs)
		http.Handle("/", fsWithGz)
	*/

	//http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	server := http.Server{
		Addr:        fmt.Sprintf(":%d", port),
		ConnContext: SaveConnInContext,
	}
	server.ListenAndServe()
}
