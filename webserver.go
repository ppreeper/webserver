package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var addr string
var dir string
var log bool
var bindval string

func main() {
	flag.StringVar(&addr, "addr", ":8080", "addr to start")
	flag.StringVar(&dir, "d", ".", "directory to serve")
	flag.BoolVar(&log, "l", false, "enable logging")
	flag.Parse()

	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	if log {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Compression on data
	r.Use(middleware.DefaultCompress)

	r.Mount("/debug", middleware.Profiler()) // ..routes return r })

	fs := http.FileServer(http.Dir(dir))
	r.Handle("/*", fs)

	if err := http.ListenAndServe(addr, MaxAge(r)); err != nil {
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
