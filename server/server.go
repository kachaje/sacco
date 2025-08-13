package server

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sacco/utils"
	"sync"

	_ "embed"
	"html/template"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Session struct {
	CurrentMenu string
	Data        map[string]string
}

//go:embed index.html
var indexHTML string

var sessions = make(map[string]*Session)
var mu sync.Mutex
var port int

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
				"4. Banking Details\n" +
				"5. Exit"
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
			session.CurrentMenu = "banking"
			goto rerunSwitch
		case "5":
			response = "END Thank you for using our service"
			mu.Lock()
			delete(sessions, sessionID)
			mu.Unlock()
		}
	case "banking":
		if text == "0" {
			session.CurrentMenu = "main"
			goto rerunSwitch
		} else {
			switch text {
			case "1":
				response = "CON National Bank of Malawi\n" +
					"Name: Kaso SACCO\n" +
					"Number: 0985 242 629\n" +
					"0. Back to Main Menu"
			case "2":
				response = "CON Airtel Money\n" +
					"Name: Kaso SACCO\n" +
					"Number: 1006857589\n" +
					"Branch: Lilongwe\n" +
					"0. Back to Main Menu"
			default:
				response = "CON Banking Details\n" +
					"1. National Bank\n" +
					"2. Airtel Money\n" +
					"0. Back to Main Menu"
			}
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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	fmt.Println(r.Method)

	phoneNumber := r.URL.Query().Get("phoneNumber")
	serviceCode := r.URL.Query().Get("serviceCode")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	sessionId := uuid.NewString()

	var text string

	for {
		data := url.Values{}
		data.Set("sessionId", sessionId)
		data.Set("text", text)
		data.Set("phoneNumber", phoneNumber)
		data.Set("serviceCode", serviceCode)

		encodedData := data.Encode()

		payload := bytes.NewBufferString(encodedData)

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/ussd", port), payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		response := regexp.MustCompile(`^CON\s|^END\s`).ReplaceAllString(string(body), "")

		err = conn.WriteMessage(websocket.TextMessage, []byte(response))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}

		text = string(message)
	}
}

func Main() {
	var err error

	flag.IntVar(&port, "p", port, "server port")

	flag.Parse()

	if port == 0 {
		port, err = utils.GetFreePort()
		if err != nil {
			log.Panic(err)
		}
	}

	http.HandleFunc("/ws", wsHandler)

	indexHTML = regexp.MustCompile("8080").ReplaceAllString(indexHTML, fmt.Sprint(port))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("index").Parse(indexHTML)
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/ussd", ussdHandler)
	log.Printf("USSD server listening on :%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
