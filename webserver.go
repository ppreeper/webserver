package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	var addr string
	var dir string
	var logging bool
	flag.StringVar(&addr, "addr", ":8080", "addr to start")
	flag.StringVar(&dir, "d", ".", "directory to serve")
	flag.BoolVar(&logging, "l", false, "enable logging")
	flag.Parse()

	r := chi.NewRouter()
	if logging {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.Compress(5, "text/html", "text/css"))
	r.Use(middleware.Recoverer)

	filesDir := http.Dir(dir)
	FileServer(r, "/", filesDir)

	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Println("webserver closed: %w", err)
		return
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
