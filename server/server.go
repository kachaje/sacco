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
	"sacco/utils"
	"strings"
	"sync"

	_ "embed"
	"html/template"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Session struct {
	CurrentMenu           string
	Data                  map[string]string
	PIWorkflow            *parser.WorkFlow
	LanguageWorkflow      *parser.WorkFlow
	OccupationWorkflow    *parser.WorkFlow
	ContactsWorkflow      *parser.WorkFlow
	NomineeWorkflow       *parser.WorkFlow
	BeneficiariesWorkflow *parser.WorkFlow
	PreferredLanguage     string
}

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

var sessions = make(map[string]*Session)
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

func checkPreferredLanguage(phoneNumber string) *string {
	settingsFile := filepath.Join(preferencesFolder, phoneNumber)

	_, err := os.Stat(settingsFile)
	if !os.IsNotExist(err) {
		content, err := os.ReadFile(settingsFile)
		if err != nil {
			log.Println(err)
			return nil
		}

		data := map[string]any{}

		err = json.Unmarshal(content, &data)
		if err != nil {
			log.Println(err)
			return nil
		}

		var preferredLanguage string

		if data["language"] != nil {
			val, ok := data["language"].(string)
			if ok {
				preferredLanguage = val
			}
		}

		return &preferredLanguage
	}

	return nil
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

	preferredLanguage := checkPreferredLanguage(phoneNumber)

	mu.Lock()
	session, exists := sessions[sessionID]
	if !exists {
		session = &Session{
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

		sessions[sessionID] = session
	}
	mu.Unlock()

	var runSwitch func(session *Session) string

	runSwitch = func(session *Session) string {
		preferredLanguage := checkPreferredLanguage(phoneNumber)

		if preferredLanguage != nil {
			session.PreferredLanguage = *preferredLanguage
		}

		var response string

		switch session.CurrentMenu {
		case "main":
			switch text {
			case "", "0":
				if preferredLanguage != nil && *preferredLanguage == "ny" {
					response = "CON Takulandilani ku Kaso SACCO\n" +
						"1. Membala Watsopano\n" +
						"2. Tengani Ngongole\n" +
						"3. Balansi\n" +
						"4. Matumizidwe\n" +
						"5. Chiyankhulo\n" +
						"6. Malizani"
				} else {
					response = "CON Welcome to Kaso SACCO\n" +
						"1. Membership Application\n" +
						"2. Loan Application\n" +
						"3. Check Balance\n" +
						"4. Banking Details\n" +
						"5. Preferred Language\n" +
						"6. Exit"
				}
			case "1":
				text = "000"
				session.CurrentMenu = "registration"
				return runSwitch(session)
			case "2":
				text = "000"
				session.CurrentMenu = "loan"
				return runSwitch(session)
			case "3":
				text = "000"
				session.CurrentMenu = "balance"
				return runSwitch(session)
			case "4":
				text = "000"
				session.CurrentMenu = "banking"
				return runSwitch(session)
			case "5":
				session.CurrentMenu = "language"
				return runSwitch(session)
			case "6":
				if preferredLanguage != nil && *preferredLanguage == "ny" {
					response = "END Zikomo potidalila"
				} else {
					response = "END Thank you for using our service"
				}
				mu.Lock()
				delete(sessions, sessionID)
				mu.Unlock()
			}
		case "language":
			if text == "" {
				session.CurrentMenu = "main"
				return runSwitch(session)
			} else {
				response = session.LanguageWorkflow.NavNext(text)

				if strings.TrimSpace(response) == "" {
					session.CurrentMenu = "main"
					text = ""
					return runSwitch(session)
				}
			}
		case "banking":
			if text == "0" {
				session.CurrentMenu = "main"
				return runSwitch(session)
			} else {
				firstLine := "CON Banking Details\n"
				lastLine := "0. Back to Main Menu"
				name := "Name"
				number := "Number"
				branch := "Branch"

				if preferredLanguage != nil && *preferredLanguage == "ny" {
					firstLine = "CON Matumizidwe\n"
					lastLine = "0. Bwererani Pofikira"
					name = "Dzina"
					number = "Nambala"
					branch = "Buranchi"
				}

				switch text {
				case "1":
					response = "CON National Bank of Malawi\n" +
						fmt.Sprintf("%8s: Kaso SACCO\n", name) +
						fmt.Sprintf("%8s: 1006857589\n", number) +
						fmt.Sprintf("%8s: Lilongwe\n", branch) +
						lastLine
				case "2":
					response = "CON Airtel Money\n" +
						fmt.Sprintf("%8s: Kaso SACCO\n", name) +
						fmt.Sprintf("%8s: 0985 242 629\n", number) +
						lastLine
				default:
					response = firstLine +
						"1. National Bank\n" +
						"2. Airtel Money\n" +
						lastLine
				}
			}
		case "registration":
			switch text {
			case "00":
				session.PIWorkflow.NavNext(text)
				session.CurrentMenu = "main"
				text = "0"
				return runSwitch(session)
			case "1":
				session.CurrentMenu = "registration.1"
				return runSwitch(session)
			default:
				if preferredLanguage != nil && *preferredLanguage == "ny" {
					response = "CON Sankhani Zochita\n" +
						"1. Zokhudza Membala\n" +
						"2. Zokhudza Ntchito\n" +
						"3. Adiresi Yamembela\n" +
						"4. Wachibale wa Membala\n" +
						"5. Odzalandila\n" +
						"\n" +
						"00. Tiyambirenso"
				} else {
					response = "CON Choose Activity\n" +
						"1. Add Member Details\n" +
						"2. Add Occupation Details\n" +
						"3. Add Contact Details\n" +
						"4. Add Next of Kin Details\n" +
						"5. Add Beneficiaries\n" +
						"\n" +
						"00. Main Menu"
				}
			}
		case "registration.1":
			response = session.PIWorkflow.NavNext(text)

			if text == "00" {
				session.CurrentMenu = "main"
				text = "0"
				return runSwitch(session)
			} else if strings.TrimSpace(response) == "" {
				session.CurrentMenu = "registration"
				text = ""
				return runSwitch(session)
			}
		case "loan":
			if text == "0" {
				session.CurrentMenu = "main"
				return runSwitch(session)
			} else {
				response = "CON Loan Application\n" +
					"0. Back to Main Menu"
			}
		case "balance":
			if text == "0" {
				session.CurrentMenu = "main"
				return runSwitch(session)
			} else {
				response = "CON Check Balance\n" +
					"0. Back to Main Menu"
			}
		}

		return response
	}

	response := runSwitch(session)

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
