package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(fs)
	http.Handle("/", r)
	if err := http.ListenAndServe(bindval, handlers.CompressHandler(r)); err != nil {
		panic(err)
	}
}
