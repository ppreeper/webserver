package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/handlers"
)

var port string
var dir string
var bindval string

func init() {
	flag.StringVar(&port, "port", "8080", "port to start")
	flag.StringVar(&dir, "d", ".", "directory to serve")
	flag.Parse()
}

func main() {
	bindval = ":" + port
	fs := http.FileServer(http.Dir(dir))
	if err := http.ListenAndServe(bindval, MaxAge(handlers.CompressHandler(fs))); err != nil {
		panic(err)
	}
}

// MaxAge sets expire headers based on extension
func MaxAge(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var age time.Duration
		ext := filepath.Ext(r.URL.String())

		// Timings are based on github.com/h5bp/server-configs-nginx

		switch ext {
		case ".rss", ".atom":
			age = time.Hour / time.Second
		case ".css", ".js":
			age = (time.Hour * 24 * 365) / time.Second
		case ".jpg", ".jpeg", ".gif", ".png", ".ico", ".cur", ".gz", ".svg", ".svgz", ".mp4", ".ogg", ".ogv", ".webm", ".htc":
			age = (time.Hour * 24 * 30) / time.Second
		default:
			age = 0
		}

		if age > 0 {
			w.Header().Add("Cache-Control", fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", age))
		}

		h.ServeHTTP(w, r)
	})
}
