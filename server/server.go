package server

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sacco/parser"
	"sacco/server/menus"
	"sacco/utils"
	"sync"

	_ "embed"
	"html/template"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

//go:embed index.html
var indexHTML string

//go:embed workflows/membership/personalInformation.yml
var PITemplate string

//go:embed workflows/membership/occupationalDetails.yml
var occupationTemplate string

//go:embed workflows/membership/contactDetails.yml
var contactsTemplate string

//go:embed workflows/membership/nomineeDetails.yml
var nomineeTemplate string

//go:embed workflows/membership/beneficiaries.yml
var beneficiariesTemplate string

//go:embed workflows/preferences/language.yml
var languageTemplate string

var mu sync.Mutex
var port int
var personalInformationData map[string]any
var languageData map[string]any
var occupationData map[string]any
var contactsData map[string]any
var nomineeData map[string]any
var beneficiariesData map[string]any

var preferencesFolder = filepath.Join(".", "settings")

func init() {
	var err error

	personalInformationData, err = utils.LoadYaml(PITemplate)
	if err != nil {
		panic(err)
	}

	languageData, err = utils.LoadYaml(languageTemplate)
	if err != nil {
		panic(err)
	}

	occupationData, err = utils.LoadYaml(occupationTemplate)
	if err != nil {
		panic(err)
	}

	contactsData, err = utils.LoadYaml(contactsTemplate)
	if err != nil {
		panic(err)
	}

	nomineeData, err = utils.LoadYaml(nomineeTemplate)
	if err != nil {
		panic(err)
	}

	beneficiariesData, err = utils.LoadYaml(beneficiariesTemplate)
	if err != nil {
		panic(err)
	}

}

func saveData(data map[string]any, model, phoneNumber *string) {
	switch *model {
	case "preferredLanguage":
		if data["language"] != nil && phoneNumber != nil {
			language, ok := data["language"].(string)
			if ok {
				savePreference(*phoneNumber, "language", language)
			}
		}
	default:
		fmt.Println("##########", data)
	}
}

func savePreference(phoneNumber, key, value string) error {
	settingsFile := filepath.Join(preferencesFolder, phoneNumber)

	data := map[string]any{}

	_, err := os.Stat(settingsFile)
	if !os.IsNotExist(err) {
		content, err := os.ReadFile(settingsFile)
		if err != nil {
			log.Println(err)
			return err
		}

		err = json.Unmarshal(content, &data)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	data[key] = value

	payload, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
		return err
	}

	return os.WriteFile(settingsFile, payload, 0644)
}

func ussdHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.FormValue("sessionId")
	serviceCode := r.FormValue("serviceCode")
	phoneNumber := r.FormValue("phoneNumber")
	text := r.FormValue("text")

	if phoneNumber == "" {
		phoneNumber = "default"
	}

	log.Printf("Received USSD request: SessionID=%s, ServiceCode=%s, PhoneNumber=%s, Text=%s",
		sessionID, serviceCode, phoneNumber, text)

	preferredLanguage := menus.CheckPreferredLanguage(phoneNumber, preferencesFolder)

	mu.Lock()
	session, exists := menus.Sessions[sessionID]
	if !exists {
		session = &menus.Session{
			CurrentMenu: "main",
			Data:        make(map[string]string),

			LanguageWorkflow: parser.NewWorkflow(languageData, saveData, preferredLanguage, &phoneNumber),

			PIWorkflow: parser.NewWorkflow(personalInformationData, saveData, preferredLanguage, &phoneNumber),

			OccupationWorkflow: parser.NewWorkflow(occupationData, saveData, preferredLanguage, &phoneNumber),

			ContactsWorkflow: parser.NewWorkflow(contactsData, saveData, preferredLanguage, &phoneNumber),

			NomineeWorkflow: parser.NewWorkflow(nomineeData, saveData, preferredLanguage, &phoneNumber),

			BeneficiariesWorkflow: parser.NewWorkflow(beneficiariesData, saveData, preferredLanguage, &phoneNumber),
		}

		if preferredLanguage != nil {
			session.PreferredLanguage = *preferredLanguage
		}

		menus.Sessions[sessionID] = session
	}
	mu.Unlock()

	response := menus.MainMenu(session, phoneNumber, text, sessionID, *preferredLanguage)

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

	_, err = os.Stat(preferencesFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(preferencesFolder, 0755)
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
