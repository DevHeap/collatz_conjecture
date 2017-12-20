package server

import (
	"log"
	"math/big"
	"net/http"

	"../calc"
	"github.com/gorilla/websocket"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {return true},
}


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
	_, input, err := ws.ReadMessage()
	if err != nil {
		log.Print("Reading input number:", err)
		return
	}

	//
	number, parsed := new(big.Int).SetString(string(input), 10)
	if !parsed {
		log.Print("Unable to parse string into integer number with base 10", err)
		return
	}

	//Create calculator with start number
	calculator := calc.NewCalculator(number, 4)
	defer calculator.Stop()

	//Cause we can be too fast in sending results to client, we limit it
	limiter := time.NewTicker(time.Millisecond * 200)
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
