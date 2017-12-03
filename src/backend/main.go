package main

import (
	"./server"
	"net/http"
	"log"
	"flag"
)
var address = flag.String("addr", "localhost:8080", "http service address")
var dir     = flag.String("dir",  "src/frontend/", "path to front files")

func main() {
	flag.Parse()

	log.Println(dir)
	http.Handle("/", http.FileServer(http.Dir(*dir)))
	http.HandleFunc("/ws", server.WSHandler)
	log.Fatal(http.ListenAndServe(*address, nil))
}
