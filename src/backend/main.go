package main

import (
	"./server"
	"./calc/cache"
	"flag"
	"log"
	"net/http"
	"time"
)

var address     = flag.String( "address",      "0.0.0.0:80",    "external address for server")
var static      = flag.String( "static",       "src/frontend/", "path to static files")
var workers     = flag.Int   ( "workers",      4,               "max workers per user")
var checkOrigin = flag.Bool  ( "check_origin", false ,          "check origin for ws")
var maxMsgRatio = flag.Duration("msg_ration", time.Millisecond * 200, "max ration of messages to client")

func main() {
	flag.Parse()

	pathCache := cache.NewCacheDefault()

	wsServer := server.NewWSServer(
		*workers,
		*maxMsgRatio,
		*checkOrigin,
		pathCache,
	)

	http.Handle("/", http.FileServer(http.Dir(*static)))
	http.HandleFunc("/ws", wsServer.WSHandler)

	log.Fatal(http.ListenAndServe(*address, nil))
}
