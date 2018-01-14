package server

import (
	"log"
	"math/big"
	"net/http"

	"../calc"
	"../calc/cache"
	"github.com/gorilla/websocket"
	"time"
	"strings"
)

const base = 10
const workers = 4
const maxMsgRatio = time.Millisecond * 200

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {return true},
}

var pathCache = cache.NewCache()

func WSHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Accepted Connection")

	//Upgrading protocol to websocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	defer ws.Close()

	if err != nil {
		log.Print("Http Upgrade:", err)
		return
	}

	//Waiting for string with initial number
	_, inputBytes, err := ws.ReadMessage()
	if err != nil {
		log.Print("Reading input number:", err)
		return
	}

	input := string(inputBytes)

	// Check that input is non-zero
	if strings.HasPrefix(input, "0") {
		log.Print("Zero as input numbers are forbidden")
		return
	}

	// Check that input is positive
	if strings.HasPrefix(input, "-") {
		log.Print("Negative input numbers are forbidden")
		return
	}

	//
	number, parsed := new(big.Int).SetString(input, base)
	if !parsed {
		log.Print("Unable to parse string into integer number", err)
		return
	}

	// Gorutine to handle close messages from client
	// From gorilla library websocket docs:
	//   "If the application is not otherwise interested in messages from the peer,
	//   then the application should start a goroutine to read and discard messages from the peer"
	go func(ws *websocket.Conn){
		for {
			if _, _, err := ws.NextReader(); err != nil {
				ws.Close()
				return
			}
		}
	}(ws)

	//Create calculator with start number
	calculator := calc.NewCalculator(number, workers, pathCache)
	defer calculator.Stop()

	//Cause we can be too fast in sending results to client, we limit it
	limiter := time.NewTicker(maxMsgRatio)
	defer limiter.Stop()

	//Main handler loop
	for result := range calculator.DataCh{
		<-limiter.C //Waiting for limit

		err = ws.WriteJSON(result)
		if err != nil {
			log.Print("Sending result:", err)
			return
		}
	}
}
