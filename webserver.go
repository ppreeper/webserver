package main

import (
	"compress/gzip"
	"errors"
	"flag"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
)

var gzPool = sync.Pool{
	New: func() interface{} {
		w := gzip.NewWriter(io.Discard)
		gzip.NewWriterLevel(w, gzip.BestCompression)
		return w
	},
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Gzip func handler
func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)

		gz.Reset(w)
		defer gz.Close()

		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func main() {
	var addr string
	var dir string
	flag.StringVar(&addr, "addr", ":8080", "addr to start")
	flag.StringVar(&dir, "d", ".", "directory to serve")
	flag.Parse()

	fs := http.FileServer(http.Dir(dir))

	slog.Info("Listening on", "addr", addr)
	err := http.ListenAndServe(addr, Gzip(fs))
	if errors.Is(err, http.ErrServerClosed) {
		slog.Info("server closed")
	} else if err != nil {
		slog.Error("error starting server", "", err)
	}
}
