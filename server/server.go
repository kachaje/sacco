package server

import (
	"bytes"
	"context"
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sacco/server/database"
	filehandling "sacco/server/fileHandling"
	"sacco/server/menus"
	"sacco/server/parser"
	"sacco/utils"
	"strings"
	"sync"
	"time"

	"html/template"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

//go:embed index.html
var indexHTML string

//go:embed workflows/*
var RawWorkflows embed.FS

var mu sync.Mutex
var port int

var workflowsData map[string]map[string]any

var preferencesFolder = filepath.Join(".", "settings")
var cacheFolder = filepath.Join(".", "data", "cache")

var ctx context.Context

var activeMenu *menus.Menus

func init() {
	var err error

	workflowsData = map[string]map[string]any{}

	err = fs.WalkDir(RawWorkflows, ".", func(file string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(file, ".yml") {
			return nil
		}

		content, err := RawWorkflows.ReadFile(file)
		if err != nil {
			return err
		}

		data, err := utils.LoadYaml(string(content))
		if err != nil {
			log.Fatal(err)
		}

		model := strings.Split(filepath.Base(file), ".")[0]

		workflowsData[model] = data

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	_, err = os.Stat(cacheFolder)
	if !os.IsNotExist(err) {
		os.MkdirAll(cacheFolder, 0755)
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
	session, exists := menus.Sessions[phoneNumber]
	if !exists {
		session = parser.NewSession(menus.DB.MemberByPhoneNumber, &phoneNumber, &sessionID)

		for model, data := range workflowsData {
			session.WorkflowsMapping[model] = parser.NewWorkflow(data, filehandling.SaveData, preferredLanguage, &phoneNumber, &sessionID, &cacheFolder, &preferencesFolder, menus.DB.GenericsSaveData, menus.Sessions, nil)
		}

		if preferredLanguage != nil {
			session.PreferredLanguage = *preferredLanguage
		}

		menus.Sessions[phoneNumber] = session
	}
	mu.Unlock()

	if phoneNumber != "default" {
		go func() {
			_, err := session.RefreshSession()
			if err == nil {
				session.UpdateSessionFlags()
			} else {
				if !strings.HasSuffix(err.Error(), "sql: no rows in result set") {
					log.Println(err)
				}
			}

			err = session.LoadMemberCache(phoneNumber, cacheFolder)
			if err != nil {
				log.Printf("server.ussdHandler: %s", err.Error())
			}
		}()
	}

	response := activeMenu.LoadMenu(session.CurrentMenu, session, phoneNumber, text, preferencesFolder, cacheFolder)

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
	var cancel context.CancelFunc
	var devMode bool

	flag.IntVar(&port, "p", port, "server port")
	flag.StringVar(&dbname, "n", dbname, "database name")
	flag.BoolVar(&devMode, "d", devMode, "dev mode")

	flag.Parse()

	if port == 0 {
		port, err = utils.GetFreePort()
		if err != nil {
			log.Panic(err)
		}
	}

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	_, err = os.Stat(preferencesFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(preferencesFolder, 0755)
	}

	_, err = os.Stat(cacheFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(cacheFolder, 0755)
	}

	menus.DB = database.NewDatabase(dbname)

	activeMenu = menus.NewMenus(&devMode)

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

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if _, err := os.Stat(cacheFolder); !os.IsNotExist(err) {
					files, err := os.ReadDir(cacheFolder)

					if err == nil {
						for _, file := range files {
							phoneNumber := file.Name()

							filehandling.RerunFailedSaves(&phoneNumber, &cacheFolder, menus.DB.GenericsSaveData, menus.Sessions)
						}
					}
				}

				time.Sleep(60 * time.Second)
			}
		}
	}()

	http.HandleFunc("/ussd", ussdHandler)
	log.Printf("USSD server listening on :%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
