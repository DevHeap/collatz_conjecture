package main

import (
	"./server"
	"flag"
	"log"
	"net/http"
)

var address = flag.String("address", "0.0.0.0:80", "external address for server")
var static = flag.String("static", "src/frontend/", "path to static files")

func main() {
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*static)))
	http.HandleFunc("/ws", server.WSHandler)
	log.Fatal(http.ListenAndServe(*address, nil))
	}
