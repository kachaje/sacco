package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sacco/forms"
	"time"

	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gorilla/websocket"
	"html/template"
)

//go:embed index.html
var indexHTML string

var conn *websocket.Conn
var err error

func WaitForPort(host string, port string, timeout time.Duration, retryInterval time.Duration, debug bool) error {
	address := net.JoinHostPort(host, port)
	startTime := time.Now()

	for {
		conn, err := net.DialTimeout("tcp", address, retryInterval)
		if err == nil {
			conn.Close()
			if debug {
				fmt.Printf("Port %s on %s is open.\n", port, host)
			}
			return nil
		}

		if time.Since(startTime) >= timeout {
			return fmt.Errorf("timeout waiting for port %s on %s: %w", port, host, err)
		}

		if debug {
			fmt.Printf("Waiting for port %s on %s... Retrying in %v\n", port, host, retryInterval)
		}

		time.Sleep(retryInterval)
	}
}

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
	if a, err = net.ResolveTCPAddr("tcp", "0.0.0.0:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func main() {
	var port int = 8080
	var interactive bool = true
	var resetting bool

	if port == 0 {
		port, err = GetFreePort()
		if err != nil {
			panic(err)
		}
	}

	flag.IntVar(&port, "p", port, "peer port")
	flag.BoolVar(&interactive, "i", interactive, "interactive mode")

	flag.Parse()

	http.HandleFunc("/ws", wsHandler)

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

	log.Printf("Server started on port %d\n", port)

	if !interactive {
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Panic(err)
		}
	} else {
		go func() {
			err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
			if err != nil {
				log.Panic(err)
			}
		}()

		err := WaitForPort("localhost", fmt.Sprint(port), 30*time.Second, 2*time.Second, false)
		if err != nil {
			log.Fatal(err)
		}

		url := fmt.Sprintf("ws://localhost:%d/ws", port)

		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer ws.Close()

		bot := forms.NewMembershipChatBot()

		var question, input string

		question = bot.ProcessInput(input)

		myApp := app.New()
		myWindow := myApp.NewWindow("Simple Fyne App")
		myWindow.Resize(fyne.NewSize(400, 200))

		questionLabel := widget.NewLabel(question)

		inputEntry := widget.NewEntry()

		outputLabel := widget.NewLabel(fmt.Sprintf("Server running on port :%d", port))

		handleSubmit := func() {
			if resetting {
				return
			}

			input = inputEntry.Text

			question = bot.ProcessInput(input)

			inputEntry.SetText("")

			if question == "" {
				resetting = true
				questionLabel.SetText("Done")

				payload, _ := json.MarshalIndent(bot.Data, "", "  ")

				outputLabel.SetText(string(payload))
			} else {
				resetting = false
				questionLabel.SetText(question)
			}
		}

		submitButton := widget.NewButton("Submit", func() {
			handleSubmit()
		})

		content := container.NewVBox(
			questionLabel,
			inputEntry,
			submitButton,
			outputLabel,
		)

		myWindow.SetContent(content)

		inputEntry.OnSubmitted = func(s string) {
			handleSubmit()
		}

		myWindow.Canvas().SetOnTypedKey(func(keyEvent *fyne.KeyEvent) {
			if keyEvent.Name == fyne.KeyReturn {
				handleSubmit()
			}
		})

		myWindow.ShowAndRun()
	}
}
