package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
)

var port string
var bindval string

func init() {
	flag.StringVar(&port, "port", "8080", "port to start")
	flag.Parse()
}

func main() {
	bindval = ":" + port
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))
	http.Handle("/", r)
	if err := http.ListenAndServe(bindval, nil); err != nil {
		panic(err)
	}
}
