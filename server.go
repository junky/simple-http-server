package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var port = 80
var ServerId = "1"

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
	chunks := []string{}
	chunks = append(chunks, "Hello, ServerId: "+ServerId+"!")
	chunks = append(chunks, "<--! chunk -->")
	return chunks
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

	EnvServerId := os.Getenv("SERVER_ID")
	if EnvServerId != "" {
		ServerId = EnvServerId
	}
	log.Info(fmt.Sprintf("Server started. port: %d, ServerId: %s", port, ServerId))

	chunks = getChunks()
	http.HandleFunc("/", chunkHandler)

	server := http.Server{
		Addr:        fmt.Sprintf(":%d", port),
		ConnContext: SaveConnInContext,
	}
	server.ListenAndServe()
}
