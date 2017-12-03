package server

import "net/http"
import (
	"github.com/gorilla/websocket"
	"log"
	"math/big"
	"backend/calc"
	"fmt"
)

var upgrader = websocket.Upgrader{}


func wsHandler(w http.ResponseWriter, r *http.Request) {
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

	fmt.Printf("Got message: %#v\n", string(input))

	number, parsed := new(big.Int).SetString(string(input), 10)
	if !parsed {
		log.Print("Unable to parse string into integer number with base 10", err)
		return
	}

	calculator := calc.NewCalculator(number, 4)
	defer calculator.Stop()

	for {
		r := <- calculator.DataCh
		err = ws.WriteJSON(r)
		if err != nil {
			log.Print("Sending result:", err)
			return
		}

		//time.Sleep(time.Second)
	}
}

func NewServer(){
	http.HandleFunc("/ws", wsHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}





