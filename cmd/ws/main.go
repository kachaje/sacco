package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sacco/forms"

	"github.com/gorilla/websocket"
)

var conn *websocket.Conn
var err error

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	bot := forms.NewMembershipChatBot()

	var input string

	for {
		question := bot.ProcessInput(input)

		if question == "" {
			payload, _ := json.MarshalIndent(bot.Data, "", "  ")

			err := conn.WriteMessage(websocket.TextMessage, payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				break
			}

			break
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(question))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		input = string(message)
	}
}

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func main() {
	var port int

	port, err = GetFreePort()
	if err != nil {
		panic(err)
	}

	flag.IntVar(&port, "p", port, "peer port")

	flag.Parse()

	http.HandleFunc("/ws", wsHandler)

	log.Printf("Server started on port %d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Panic(err)
	}
}
