package main

import (
	"flag"
	"net/http"

	"github.com/NYTimes/gziphandler"
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
	mux := http.NewServeMux()
	mux.Handle("/", fs)
	withGz := gziphandler.GzipHandler(mux)
	http.Handle("/", withGz)
	if err := http.ListenAndServe(bindval, nil); err != nil {
		panic(err)
	}
}
