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
	"sacco/server/database"
	filehandling "sacco/server/fileHandling"
	"sacco/server/menus"
	menufuncs "sacco/server/menus/menuFuncs"
	"sacco/server/parser"
	"sacco/utils"
	"strings"
	"sync"

	"html/template"

	_ "embed"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

//go:embed index.html
var indexHTML string

var mu sync.Mutex
var port int
var demoMode bool

var preferencesFolder = filepath.Join(".", "settings")

var activeMenu *menus.Menus

func ussdHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.FormValue("sessionId")
	serviceCode := r.FormValue("serviceCode")
	phoneNumber := r.FormValue("phoneNumber")
	text := r.FormValue("text")

	defaultPhoneNumber := "000000000"

	if phoneNumber == "" {
		phoneNumber = defaultPhoneNumber
	}

	log.Printf("Received USSD request: SessionID=%s, ServiceCode=%s, PhoneNumber=%s, Text=%s",
		sessionID, serviceCode, phoneNumber, text)

	preferredLanguage := menufuncs.CheckPreferredLanguage(phoneNumber, preferencesFolder)

	mu.Lock()
	session, exists := menufuncs.Sessions[phoneNumber]
	if !exists {
		session = parser.NewSession(menufuncs.DB.MemberByPhoneNumber, &phoneNumber, &sessionID)

		for model, data := range menufuncs.WorkflowsData {
			session.WorkflowsMapping[model] = parser.NewWorkflow(data, filehandling.SaveData, preferredLanguage, &phoneNumber, &sessionID, &preferencesFolder, menufuncs.DB.GenericsSaveData, menufuncs.Sessions, nil)
		}

		if preferredLanguage != nil {
			session.PreferredLanguage = *preferredLanguage
		}

		if demoMode {
			defaultUser := "admin"
			defaultUserId := int64(1)
			defaultRole := "admin"

			session.SessionUser = &defaultUser
			session.SessionUserId = &defaultUserId
			session.SessionUserRole = &defaultRole
		}

		menufuncs.Sessions[phoneNumber] = session
	}
	mu.Unlock()

	go func() {
		_, err := session.RefreshSession()
		if err == nil {
			session.UpdateSessionFlags()
		} else {
			if !strings.HasSuffix(err.Error(), "sql: no rows in result set") {
				log.Println(err)
			}
		}
	}()

	response := activeMenu.LoadMenu(session.CurrentMenu, session, phoneNumber, text, preferencesFolder)

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
	var devMode bool

	flag.IntVar(&port, "p", port, "server port")
	flag.StringVar(&dbname, "n", dbname, "database name")
	flag.BoolVar(&devMode, "d", devMode, "dev mode")
	flag.BoolVar(&demoMode, "o", demoMode, "demo mode")

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

	menufuncs.DB = database.NewDatabase(dbname)

	_, err = menufuncs.DB.DB.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		panic(err)
	}

	activeMenu = menus.NewMenus(&devMode, &demoMode)

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
