package menus

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sacco/server/database"
	menufuncs "sacco/server/menus/menuFuncs"
	"sacco/server/parser"
	"sacco/utils"
	"slices"
	"strings"
	"sync"
)

//go:embed menus/*
var menuFiles embed.FS

//go:embed templates/member.template.json
var menuTemplateContent []byte

//go:embed templates/member.template.json
var templateContent []byte

var DB *database.Database

var Sessions map[string]*parser.Session
var templateData map[string]any

type Menus struct {
	ActiveMenus   map[string]any
	Titles        map[string]string
	Workflows     map[string]any
	Functions     map[string]any
	FunctionsMap  map[string]func(map[string]any) string
	TargetKeys    map[string][]string
	LabelWorkflow map[string]any

	mu sync.Mutex

	DevModeActive bool
	DemoMode      bool

	Cache      map[string]string
	LastPrompt string
}

var menuTemplateData map[string]any

func init() {
	err := json.Unmarshal(menuTemplateContent, &menuTemplateData)
	if err != nil {
		log.Fatalf("server.menus.init: %s", err.Error())
	}

	err = json.Unmarshal(templateContent, &templateData)
	if err != nil {
		log.Fatalf("server.menus.init: %s", err.Error())
	}

	Sessions = map[string]*parser.Session{}
}

func NewMenus(devMode, demoMode *bool) *Menus {
	m := &Menus{
		ActiveMenus:   map[string]any{},
		Titles:        map[string]string{},
		Workflows:     map[string]any{},
		Functions:     map[string]any{},
		FunctionsMap:  map[string]func(map[string]any) string{},
		TargetKeys:    map[string][]string{},
		LabelWorkflow: map[string]any{},
		mu:            sync.Mutex{},

		Cache:      map[string]string{},
		LastPrompt: "username",
	}

	if devMode != nil {
		m.DevModeActive = *devMode
	}
	if demoMode != nil {
		m.DemoMode = *demoMode
	}

	m.FunctionsMap["doExit"] = func(data map[string]any) string {
		return m.doExit(data)
	}
	m.FunctionsMap["businessSummary"] = func(data map[string]any) string {
		return m.businessSummary(data)
	}
	m.FunctionsMap["employmentSummary"] = func(data map[string]any) string {
		return m.employmentSummary(data)
	}
	m.FunctionsMap["checkBalance"] = func(data map[string]any) string {
		return m.checkBalance(data)
	}
	m.FunctionsMap["bankingDetails"] = func(data map[string]any) string {
		return m.bankingDetails(data)
	}
	m.FunctionsMap["viewMemberDetails"] = func(data map[string]any) string {
		return m.viewMemberDetails(data)
	}
	m.FunctionsMap["devConsole"] = func(data map[string]any) string {
		return menufuncs.DevConsole(DB, data)
	}
	m.FunctionsMap["memberLoansSummary"] = func(data map[string]any) string {
		return m.memberLoansSummary(data)
	}
	m.FunctionsMap["signIn"] = func(data map[string]any) string {
		return menufuncs.SignIn(m.LoadMenu, DB, data)
	}
	m.FunctionsMap["listUsers"] = func(data map[string]any) string {
		return m.listUsers(data)
	}
	m.FunctionsMap["blockUser"] = func(data map[string]any) string {
		return m.blockUser(data)
	}
	m.FunctionsMap["editUser"] = func(data map[string]any) string {
		return m.editUser(data)
	}
	m.FunctionsMap["changePassword"] = func(data map[string]any) string {
		return m.changePassword(data)
	}
	m.FunctionsMap["signUp"] = func(data map[string]any) string {
		return menufuncs.SignUp(m.LoadMenu, DB, data)
	}
	m.FunctionsMap["landing"] = func(data map[string]any) string {
		return menufuncs.Landing(m.LoadMenu, DB, data)
	}

	err := fs.WalkDir(menuFiles, ".", func(file string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		content, err := menuFiles.ReadFile(file)
		if err != nil {
			return err
		}

		data, err := utils.LoadYaml(string(content))
		if err != nil {
			log.Fatal(err)
		}

		re := regexp.MustCompile("Menu$")

		group := re.ReplaceAllLiteralString(strings.Split(filepath.Base(file), ".")[0], "")

		if val, ok := data["title"].(string); ok {
			m.Titles[group] = val
		}

		m.LabelWorkflow[group] = map[string]any{}

		if val, ok := data["fields"].(map[string]any); ok {
			m.ActiveMenus[group] = map[string]any{}

			keys := []string{}
			values := []string{}
			kv := map[string]any{}
			devMenus := map[string]any{}

			for key, row := range val {
				keys = append(keys, key)

				if val, ok := row.(map[string]any); ok {
					if val["id"] != nil && val["label"] != nil && val["label"].(map[string]any)["en"] != nil {
						if val["devOnly"] != nil && !m.DevModeActive {
							continue
						}

						id := fmt.Sprintf("%v", val["id"])
						label := fmt.Sprintf("%v", val["label"].(map[string]any)["en"])

						kv[key] = id

						value := fmt.Sprintf("%v. %v\n", key, label)

						values = append(values, value)

						if val["workflow"] != nil {
							if v, ok := val["workflow"].(string); ok {
								m.Workflows[id] = v

								m.LabelWorkflow[group].(map[string]any)[value] = map[string]any{
									"model": v,
									"id":    id,
								}
							}
						}
						if val["function"] != nil {
							if v, ok := val["function"].(string); ok {
								m.Functions[id] = v
							}
						}
						if val["targetKeys"] != nil {
							if v, ok := val["targetKeys"].([]any); ok {
								m.TargetKeys[id] = []string{}

								for _, e := range v {
									if s, ok := e.(string); ok {
										m.TargetKeys[id] = append(m.TargetKeys[id], s)
									}
								}
							}
						}
					}
				}
			}

			m.ActiveMenus[group].(map[string]any)["keys"] = keys
			m.ActiveMenus[group].(map[string]any)["kv"] = kv
			m.ActiveMenus[group].(map[string]any)["values"] = values
			m.ActiveMenus[group].(map[string]any)["devMenus"] = devMenus
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return m
}

func (m *Menus) LoadMenu(menuName string, session *parser.Session, phoneNumber, text, preferencesFolder, cacheFolder string) string {
	var response string

	preferredLanguage := CheckPreferredLanguage(phoneNumber, preferencesFolder)

	if preferredLanguage != nil {
		session.PreferredLanguage = *preferredLanguage
	}

	if session == nil {
		return response
	}

	if session.SessionToken == nil && !m.DemoMode {
		switch session.CurrentMenu {
		case "signIn":
			return menufuncs.SignIn(
				m.LoadMenu, DB,
				map[string]any{
					"phoneNumber":       phoneNumber,
					"cacheFolder":       cacheFolder,
					"session":           session,
					"preferredLanguage": preferredLanguage,
					"preferencesFolder": preferencesFolder,
					"text":              text,
				})
		case "signUp":
			return menufuncs.SignUp(
				m.LoadMenu, DB,
				map[string]any{
					"phoneNumber":       phoneNumber,
					"cacheFolder":       cacheFolder,
					"session":           session,
					"preferredLanguage": preferredLanguage,
					"preferencesFolder": preferencesFolder,
					"text":              text,
				})
		default:
			return menufuncs.Landing(
				m.LoadMenu, DB,
				map[string]any{
					"phoneNumber":       phoneNumber,
					"cacheFolder":       cacheFolder,
					"session":           session,
					"preferredLanguage": preferredLanguage,
					"preferencesFolder": preferencesFolder,
					"text":              text,
				})
		}
	}

	keys := []string{}
	values := []string{}
	kv := map[string]string{}

	if val, ok := m.ActiveMenus[menuName].(map[string]any); ok {
		if val["keys"] != nil {
			if v, ok := val["keys"].([]string); ok {
				keys = v
			}
		}
		if val["values"] != nil {
			if v, ok := val["values"].([]string); ok {
				values = v
			}
		}
		if val["kv"] != nil {
			if v, ok := val["kv"].(map[string]any); ok {
				for key, value := range v {
					if vs, ok := value.(string); ok {
						kv[key] = vs
					}
				}
			}
		}
	}

	if slices.Contains(keys, text) {
		target := text
		text = "000"

		if session != nil {
			session.CurrentMenu = kv[target]
		}

		return m.LoadMenu(kv[target], session, phoneNumber, text, preferencesFolder, cacheFolder)
	} else if session != nil && m.Workflows[session.CurrentMenu] != nil {
		workingMenu := session.CurrentMenu
		model := fmt.Sprintf("%v", m.Workflows[workingMenu])

		if session.ActiveData != nil {
			if regexp.MustCompile(`^\d+$`).MatchString(phoneNumber) && session.WorkflowsMapping != nil &&
				session.WorkflowsMapping[model] != nil {
				if m.TargetKeys[workingMenu] != nil {
					targetKeys := m.TargetKeys[workingMenu]

					updateArrayRow := func(row map[string]any, i int) {
						for key, value := range row {
							localKey := fmt.Sprintf("%s%d", key, i+1)

							if slices.Contains(targetKeys, key) && session.WorkflowsMapping[model].Data[localKey] == nil {

								session.WorkflowsMapping[model].Data[localKey] = fmt.Sprintf("%v", value)
							}
						}
					}

					if session.ActiveData[model] != nil {
						if val, ok := session.ActiveData[model].(map[string]any); ok {
							for key, value := range val {
								if slices.Contains(targetKeys, key) && session.WorkflowsMapping[model].Data[key] == nil {
									session.WorkflowsMapping[model].Data[key] = fmt.Sprintf("%v", value)
								}
							}
						} else if val, ok := session.ActiveData[model].([]any); ok {
							for i, row := range val {
								if rowVal, ok := row.(map[string]any); ok {
									updateArrayRow(rowVal, i)
								}
							}
						} else if val, ok := session.ActiveData[model].([]map[string]any); ok {
							for i, row := range val {
								updateArrayRow(row, i)
							}
						}
					} else {
						for key, value := range session.ActiveData {
							if slices.Contains(targetKeys, key) && session.WorkflowsMapping[model].Data[key] == nil {
								session.WorkflowsMapping[model].Data[key] = fmt.Sprintf("%v", value)
							}
						}
					}
				}

				if session.WorkflowsMapping[model].Data["phoneNumber"] == nil {
					session.WorkflowsMapping[model].Data["phoneNumber"] = phoneNumber
				}
			}
		}

		if session.WorkflowsMapping != nil &&
			session.WorkflowsMapping[model] != nil {
			response = session.WorkflowsMapping[model].NavNext(text)

			if text == "00" {
				session.CurrentMenu = "main"
				text = "0"
				return m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, preferencesFolder, cacheFolder)
			} else if strings.TrimSpace(response) == "" {
				if text == "0" {
					session.AddedModels[model] = true
				}

				parentMenu := "main"

				if regexp.MustCompile(`\.\d+$`).MatchString(session.CurrentMenu) {
					parentMenu = regexp.MustCompile(`\.\d+$`).ReplaceAllLiteralString(session.CurrentMenu, "")
				}

				session.CurrentMenu = parentMenu
				text = ""

				return m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, preferencesFolder, cacheFolder)
			}
		} else {
			if text == "00" {
				session.CurrentMenu = "main"
				text = "0"
				return m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, preferencesFolder, cacheFolder)
			}

			response = "NOT IMPLEMENTED YET\n\n" +
				"00. Main Menu\n"
		}

	} else {
		var menuRoot string

		if session != nil {
			menuRoot = session.CurrentMenu

			if regexp.MustCompile(`\.\d+$`).MatchString(session.CurrentMenu) {
				menuRoot = strings.Split(session.CurrentMenu, ".")[0]
			}

			if m.Functions[menuRoot] == nil && m.Functions[session.CurrentMenu] != nil {
				menuRoot = session.CurrentMenu
			}
		}

		if session != nil && m.Functions[menuRoot] != nil {
			if text == "00" {
				session.CurrentMenu = "main"
				text = "0"
				return m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, preferencesFolder, cacheFolder)
			} else {
				if fnName, ok := m.Functions[menuRoot].(string); ok && m.FunctionsMap[fnName] != nil {
					response = m.FunctionsMap[fnName](map[string]any{
						"phoneNumber":       phoneNumber,
						"cacheFolder":       cacheFolder,
						"session":           session,
						"preferredLanguage": preferredLanguage,
						"preferencesFolder": preferencesFolder,
						"text":              text,
					})
				} else {
					response = fmt.Sprintf("Function %s not found\n\n", m.Functions[menuRoot]) +
						"00. Main Menu\n"
				}
			}
		} else {
			newValues := []string{}

			if m.LabelWorkflow[menuName] != nil && session != nil {
				for _, value := range values {
					if m.LabelWorkflow[menuName].(map[string]any)[value] != nil {
						model := fmt.Sprintf("%v", m.LabelWorkflow[menuName].(map[string]any)[value].(map[string]any)["model"])

						suffix := ""

						if session.AddedModels[model] {
							suffix = "(*)"
						}

						newValues = append(newValues, fmt.Sprintf("%s %s\n", strings.TrimSpace(value), suffix))
					} else {
						newValues = append(newValues, value)
					}
				}
			} else {
				newValues = values
			}

			slices.Sort(newValues)

			index := utils.Index(newValues, "99. Cancel\n")

			if index >= 0 {
				newValues = append(newValues[:index], newValues[index+1:]...)

				newValues = append(newValues, "\n99. Cancel")
			}

			index = utils.Index(newValues, "00. Main Menu\n")

			if index >= 0 {
				newValues = append(newValues[:index], newValues[index+1:]...)

				newValues = append(newValues, "\n00. Main Menu\n")
			}

			response = fmt.Sprintf("CON %s\n%s", m.Titles[menuName], strings.Join(newValues, ""))
		}
	}

	return response
}

func (m *Menus) doExit(data map[string]any) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var session *parser.Session
	var phoneNumber string
	var cacheFolder string

	if data != nil {
		if data["session"] != nil {
			if val, ok := data["session"].(*parser.Session); ok {
				session = val
			}
		}
		if data["phoneNumber"] != nil {
			if val, ok := data["phoneNumber"].(string); ok {
				phoneNumber = val
			}
		}
		if data["cacheFolder"] != nil {
			if val, ok := data["cacheFolder"].(string); ok {
				cacheFolder = val
			}
		}

		m.Cache = map[string]string{}
		m.LastPrompt = "username"
		session.SessionToken = nil

		if phoneNumber != "" {
			delete(Sessions, phoneNumber)

			if cacheFolder != "" {
				folderName := filepath.Join(cacheFolder, phoneNumber)

				_, err := os.Stat(folderName)
				if !os.IsNotExist(err) {
					files, err := os.ReadDir(folderName)
					if err == nil && len(files) == 0 {
						err = os.RemoveAll(folderName)
						if err != nil {
							log.Printf("server.menus.menu.removeFolder: %s\n", err.Error())
						}
					}
				}
			}
		}
	}

	return "END Thank you for using our service"
}

func (m *Menus) businessSummary(data map[string]any) string {
	var result string = "Business Summary\n\n" +
		"00. Main Menu\n"

	_ = data

	return result
}

func (m *Menus) employmentSummary(data map[string]any) string {
	var result string = "Employment Summary\n\n" +
		"00. Main Menu\n"

	_ = data
	return result
}

func (m *Menus) checkBalance(data map[string]any) string {
	var result string = "Check Balance\n\n" +
		"00. Main Menu\n"

	_ = data

	return result
}

func (m *Menus) bankingDetails(data map[string]any) string {
	var preferredLanguage *string
	var response, text string

	if data["preferredLanguage"] != nil {
		if val, ok := data["preferredLanguage"].(*string); ok {
			preferredLanguage = val
		}
	}
	if data["text"] != nil {
		if val, ok := data["text"].(string); ok {
			text = val
		}
	}

	firstLine := "CON Banking Details\n"
	lastLine := "00. Main Menu\n"
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
			fmt.Sprintf("%-8s: Kaso SACCO\n", name) +
			fmt.Sprintf("%-8s: 1006857589\n", number) +
			fmt.Sprintf("%-8s: Lilongwe\n", branch) +
			"\n99. Cancel\n" +
			lastLine
	case "2":
		response = "CON Airtel Money\n" +
			fmt.Sprintf("%-8s: Kaso SACCO\n", name) +
			fmt.Sprintf("%-8s: 0985 242 629\n", number) +
			"\n99. Cancel\n" +
			lastLine
	default:
		response = firstLine +
			"1. National Bank\n" +
			"2. Airtel Money\n" +
			"\n" +
			lastLine
	}

	return response
}

func (m *Menus) viewMemberDetails(data map[string]any) string {
	var session *parser.Session
	var preferredLanguage *string
	var response string
	var phoneNumber, text, preferencesFolder, cacheFolder string

	if data["session"] != nil {
		if val, ok := data["session"].(*parser.Session); ok {
			session = val
		}
	}
	if data["preferredLanguage"] != nil {
		if val, ok := data["preferredLanguage"].(*string); ok {
			preferredLanguage = val
		}
	}
	if data["phoneNumber"] != nil {
		if val, ok := data["phoneNumber"].(string); ok {
			phoneNumber = val
		}
	}
	if data["text"] != nil {
		if val, ok := data["text"].(string); ok {
			text = val
		}
	}
	if data["preferencesFolder"] != nil {
		if val, ok := data["preferencesFolder"].(string); ok {
			preferencesFolder = val
		}
	}
	if data["cacheFolder"] != nil {
		if val, ok := data["cacheFolder"].(string); ok {
			cacheFolder = val
		}
	}

	if session != nil {
		if strings.TrimSpace(text) == "99" {
			parentMenu := "main"

			if regexp.MustCompile(`\.\d+$`).MatchString(session.CurrentMenu) {
				parentMenu = regexp.MustCompile(`\.\d+$`).ReplaceAllLiteralString(session.CurrentMenu, "")
			}

			session.CurrentMenu = parentMenu
			text = ""
			return m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, preferencesFolder, cacheFolder)
		} else {
			data = LoadTemplateData(session.ActiveData, templateData)

			table := TabulateData(data)

			tableString := strings.Join(table, "\n")

			if preferredLanguage != nil && *preferredLanguage == "ny" {
				response = "CON Zambiri za Membala\n" +
					"\n" +
					fmt.Sprintf("%s\n", tableString) +
					"\n" +
					"99. Basi\n" +
					"00. Tiyambirenso"
			} else {
				response = "CON Member Details\n" +
					"\n" +
					fmt.Sprintf("%s\n", tableString) +
					"\n" +
					"99. Cancel\n" +
					"00. Main Menu"
			}
		}
	} else {
		response = "Member Details\n\n" +
			"00. Main Menu\n"
	}

	return response
}

func (m *Menus) memberLoansSummary(data map[string]any) string {
	var response string = "Member Loans Summary\n\n00. Main Menu"

	_ = data

	return response
}

func (m *Menus) listUsers(_ map[string]any) string {
	var response, content string

	title := "Users List\n----------"

	result, err := DB.SQLQuery("SELECT id, username, role, created_at, updated_at FROM user WHERE active = 1")
	if err != nil {
		content = fmt.Sprintf("%s\n", err.Error())
	} else {
		rows := []string{
			fmt.Sprintf("%2s | %-10s | %-8s | %-20s | %-20s", "id", "username", "role", "created_at", "updated_at"),
			strings.Repeat("-", 70),
		}

		for _, row := range result {
			id := fmt.Sprintf("%v", row["id"])
			username := fmt.Sprintf("%v", row["username"])
			role := fmt.Sprintf("%v", row["role"])
			createdAt := fmt.Sprintf("%v", row["created_at"])
			updatedAt := fmt.Sprintf("%v", row["updated_at"])

			entry := fmt.Sprintf("%2s | %-10s | %-8s | %-20s | %-20s", id, username, role, createdAt, updatedAt)

			rows = append(rows, entry)
		}

		content = strings.Join(rows, "\n")
	}

	response = fmt.Sprintf("%s\n\n%s\n\n00. Main Menu\n", title, content)

	return response
}

func (m *Menus) blockUser(data map[string]any) string {
	var response string = "Block User\n\n00. Main Menu"

	_ = data

	return response
}

func (m *Menus) editUser(data map[string]any) string {
	var response string = "Edit User\n\n00. Main Menu"

	_ = data

	return response
}

func (m *Menus) changePassword(data map[string]any) string {
	var response, text string
	var session *parser.Session

	if data["session"] != nil {
		if val, ok := data["session"].(*parser.Session); ok {
			session = val
		}
	}

	if data["text"] != nil {
		if val, ok := data["text"].(string); ok {
			text = val
		}
	}

	currentPassword := func(msg string) string {
		return fmt.Sprintf("Change Password\n"+
			"---------------\n\n"+
			"Current Password %s\n", msg)
	}
	newPassword := func(msg string) string {
		return fmt.Sprintf("Change Password\n"+
			"---------------\n\n"+
			"New Password %s\n", msg)
	}
	confirmPassword := func(msg string) string {
		return fmt.Sprintf("Change Password\n"+
			"---------------\n\n"+
			"Confirm Password %s\n", msg)
	}

	switch m.LastPrompt {
	case "currentPassword":
		if text == "" {
			response = currentPassword("(Required Field)")
		} else {
			session.Cache["currentPassword"] = text

			m.LastPrompt = "newPassword"

			response = newPassword("")
		}
	case "newPassword":
		if text == "" {
			response = newPassword("(Required Field)")
		} else {
			session.Cache["newPassword"] = text

			m.LastPrompt = "confirmPassword"

			response = confirmPassword("")
		}
	case "confirmPassword":
		if text == "" {
			response = confirmPassword("(Required Field)")
		} else {
			session.Cache["confirmPassword"] = text

			text = ""

			if session.Cache["newPassword"] != session.Cache["confirmPassword"] {
				m.LastPrompt = "newPassword"

				response = newPassword("(Password Mismatch!)")
			} else {
				if id, ok := DB.ValidatePassword(*session.SessionUser, session.Cache["currentPassword"]); ok {
					err := DB.GenericModels["user"].UpdateRecord(map[string]any{
						"password": session.Cache["newPassword"],
					}, *id)
					if err != nil {
						response = fmt.Sprintf("Change Password\n"+
							"-------------\n"+
							"ERROR: %s\n\n00. Main Menu", err.Error())
					} else {
						session.Cache = map[string]string{}
						session.LastPrompt = ""

						response = "Change Password\n" +
							"-------------\n" +
							"Password Changed!\n\n00. Main Menu"
					}
				} else {
					m.LastPrompt = "currentPassword"

					response = currentPassword("(Invalid credentials)")
				}
			}
		}
	default:
		m.LastPrompt = "currentPassword"

		response = currentPassword("")
	}

	return response
}
