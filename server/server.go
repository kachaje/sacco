package server

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sacco/utils"
	"sync"
)

type Session struct {
	CurrentMenu string
	Data        map[string]string
}

var sessions = make(map[string]*Session)
var mu sync.Mutex

func ussdHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.FormValue("sessionId")
	serviceCode := r.FormValue("serviceCode")
	phoneNumber := r.FormValue("phoneNumber")
	text := r.FormValue("text")

	log.Printf("Received USSD request: SessionID=%s, ServiceCode=%s, PhoneNumber=%s, Text=%s",
		sessionID, serviceCode, phoneNumber, text)

	mu.Lock()
	session, exists := sessions[sessionID]
	if !exists {
		session = &Session{
			CurrentMenu: "main",
			Data:        make(map[string]string),
		}
		sessions[sessionID] = session
	}
	mu.Unlock()

	var response string

rerunSwitch:
	switch session.CurrentMenu {
	case "main":
		switch text {
		case "", "0":
			response = "CON Welcome to Kaso SACCO\n" +
				"1. Membership Application\n" +
				"2. Loan Application\n" +
				"3. Check Balance\n" +
				"4. Exit"
		case "1":
			session.CurrentMenu = "registration"
			goto rerunSwitch
		case "2":
			session.CurrentMenu = "loan"
			goto rerunSwitch
		case "3":
			session.CurrentMenu = "balance"
			goto rerunSwitch
		case "4":
			response = "END Thank you for using our service"
			mu.Lock()
			delete(sessions, sessionID)
			mu.Unlock()
		}
	case "registration":
		if text == "0" {
			session.CurrentMenu = "main"
			goto rerunSwitch
		} else {
			response = "CON Membership Application\n" +
				"0. Back to Main Menu"
		}
	case "loan":
		if text == "0" {
			session.CurrentMenu = "main"
			goto rerunSwitch
		} else {
			response = "CON Loan Application\n" +
				"0. Back to Main Menu"
		}
	case "balance":
		if text == "0" {
			session.CurrentMenu = "main"
			goto rerunSwitch
		} else {
			response = "CON Check Balance\n" +
				"0. Back to Main Menu"
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, response)
}

func Main() {
	var port int
	var err error

	flag.IntVar(&port, "p", port, "server port")

	flag.Parse()

	if port == 0 {
		port, err = utils.GetFreePort()
		if err != nil {
			log.Panic(err)
		}
	}

	http.HandleFunc("/ussd", ussdHandler)
	log.Printf("USSD server listening on :%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
