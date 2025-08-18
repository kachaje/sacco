package server

import (
	"bytes"
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
	"sacco/server/database"
	"sacco/server/menus"
	"sacco/utils"
	"strings"
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

//go:embed workflows/membership/occupationDetails.yml
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
var cacheFolder = filepath.Join(".", "data", "cache")

var db *database.Database

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
		session = &parser.Session{
			CurrentMenu: "main",
			Data:        make(map[string]string),
			SessionId:   sessionID,
			PhoneNumber: phoneNumber,

			LanguageWorkflow: parser.NewWorkflow(languageData, SaveData, preferredLanguage, &phoneNumber, &sessionID, &cacheFolder, &preferencesFolder, db.AddMember, menus.Sessions),

			PIWorkflow: parser.NewWorkflow(personalInformationData, SaveData, preferredLanguage, &phoneNumber, &sessionID, &cacheFolder, &preferencesFolder, db.AddMember, menus.Sessions),

			OccupationWorkflow: parser.NewWorkflow(occupationData, SaveData, preferredLanguage, &phoneNumber, &sessionID, &cacheFolder, &preferencesFolder, db.AddMember, menus.Sessions),

			ContactsWorkflow: parser.NewWorkflow(contactsData, SaveData, preferredLanguage, &phoneNumber, &sessionID, &cacheFolder, &preferencesFolder, db.AddMember, menus.Sessions),

			NomineeWorkflow: parser.NewWorkflow(nomineeData, SaveData, preferredLanguage, &phoneNumber, &sessionID, &cacheFolder, &preferencesFolder, db.AddMember, menus.Sessions),

			BeneficiariesWorkflow: parser.NewWorkflow(beneficiariesData, SaveData, preferredLanguage, &phoneNumber, &sessionID, &cacheFolder, &preferencesFolder, db.AddMember, menus.Sessions),
		}

		if preferredLanguage != nil {
			session.PreferredLanguage = *preferredLanguage
		}

		menus.Sessions[sessionID] = session
	}
	mu.Unlock()

	if phoneNumber != "default" {
		go func() {
			data, err := db.MemberByDefaultPhoneNumber(phoneNumber)
			if err == nil {
				session.ActiveMemberData = data

				session.UpdateSessionFlags()
			} else {
				if !strings.HasSuffix(err.Error(), "sql: no rows in result set") {
					log.Println(err)
				}

				session.LoadMemberCache(phoneNumber, cacheFolder)
			}
		}()
	}

	response := menus.MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)

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

	phoneNumber := r.URL.Query().Get("phoneNumber")
	serviceCode := r.URL.Query().Get("serviceCode")
	sessionId := r.URL.Query().Get("sessionId")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	if sessionId == "" {
		sessionId = uuid.NewString()
	}

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
	var dbname string = ":memory:"

	flag.IntVar(&port, "p", port, "server port")
	flag.StringVar(&dbname, "n", dbname, "database name")

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

	_, err = os.Stat(cacheFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(cacheFolder, 0755)
	}

	db = database.NewDatabase(dbname)

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
