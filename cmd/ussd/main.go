package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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

	switch session.CurrentMenu {
	case "main":
		switch text {
		case "", "0":
			response = "CON Welcome to our USSD service!\n" +
				"1. Check Balance\n" +
				"2. Buy Airtime\n" +
				"3. Exit"
		case "1":
			session.CurrentMenu = "balance"
			response = "CON Your current balance is $100.00\n" +
				"0. Back to Main Menu"
		case "2":
			session.CurrentMenu = "airtime"
			response = "CON Enter amount to buy airtime:\n" +
				"0. Back to Main Menu"
		case "3":
			response = "END Thank you for using our service!"
			mu.Lock()
			delete(sessions, sessionID)
			mu.Unlock()
		default:
			response = "CON Invalid input. Please try again.\n" +
				"1. Check Balance\n" +
				"2. Buy Airtime\n" +
				"3. Exit"
		}
	case "balance":
		if text == "0" {
			session.CurrentMenu = "main"
			response = "CON Welcome to our USSD service!\n" +
				"1. Check Balance\n" +
				"2. Buy Airtime\n" +
				"3. Exit"
		} else {
			response = "CON Invalid input. Please try again.\n" +
				"0. Back to Main Menu"
		}
	case "airtime":
		if text == "0" {
			session.CurrentMenu = "main"
			response = "CON Welcome to our USSD service!\n" +
				"1. Check Balance\n" +
				"2. Buy Airtime\n" +
				"3. Exit"
		} else if _, err := strconv.Atoi(text); err == nil {
			session.Data["airtimeAmount"] = text
			response = fmt.Sprintf("END You are about to buy airtime worth $%s. Thank you!", text)
			mu.Lock()
			delete(sessions, sessionID)
			mu.Unlock()
		} else {
			response = "CON Invalid amount. Please enter a number:\n" +
				"0. Back to Main Menu"
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, response)
}

func main() {
	http.HandleFunc("/ussd", ussdHandler)
	log.Println("USSD server listening on :8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
