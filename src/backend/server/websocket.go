package server

import (
	"fmt"
	"log"
	"math/big"
	"net/http"

	"../calc"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}


func WSHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Connection")

	ws, err := upgrader.Upgrade(w, r, nil)
	defer ws.Close()
	if err != nil {
		log.Print("Http Upgrade:", err)
		return
	}

	_, input, err := ws.ReadMessage()
	if err != nil {
		log.Print("Reading input number:", err)
		return
	}


	number, parsed := new(big.Int).SetString(string(input), 10)
	if !parsed {
		log.Print("Unable to parse string into integer number with base 10", err)
		return
	}

	calculator := calc.NewCalculator(number, 4)
	defer calculator.Stop()

	for {
		r := <-calculator.DataCh
		err = ws.WriteJSON(r)
		if err != nil {
			log.Print("Sending result:", err)
			return
		}

		//time.Sleep(time.Second)
	}
}
