package server

import (
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "index.html")
}