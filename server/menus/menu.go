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
	"sacco/server/parser"
	"sacco/utils"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
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

func NewMenus(devMode *bool) *Menus {
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
		return m.devConsole(data)
	}
	m.FunctionsMap["memberLoansSummary"] = func(data map[string]any) string {
		return m.memberLoansSummary(data)
	}
	m.FunctionsMap["login"] = func(data map[string]any) string {
		return m.login(data)
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

	if session.SessionToken == nil {
		return m.login(map[string]any{
			"phoneNumber":       phoneNumber,
			"cacheFolder":       cacheFolder,
			"session":           session,
			"preferredLanguage": preferredLanguage,
			"preferencesFolder": preferencesFolder,
			"text":              text,
		})
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

	var phoneNumber string
	var cacheFolder string

	if data != nil {
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

func (m *Menus) devConsole(data map[string]any) string {
	var session *parser.Session
	var response, content, text, title string

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

	if session != nil {
		if text == "99" {
			session.CurrentMenu = "console"
		} else if session.CurrentMenu == "console" && regexp.MustCompile(`^\d+$`).MatchString(text) {
			session.CurrentMenu = fmt.Sprintf("%s.%s", session.CurrentMenu, text)
			text = ""
		}

		switch session.CurrentMenu {
		case "console.1":
			title = "WorkflowsMapping"

			if session.WorkflowsMapping != nil {
				data := map[string]any{}

				for key, wflow := range session.WorkflowsMapping {
					row := map[string]any{
						"data":           wflow.Data,
						"optionalFields": wflow.OptionalFields,
						"screenOrder":    wflow.ScreenOrder,
						"history":        wflow.History,
					}

					data[key] = row
				}

				payload, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					content = err.Error()
				} else {
					content = string(payload)
				}
			}
		case "console.2":
			title = "AddedModels"

			if session.AddedModels != nil {
				payload, err := json.MarshalIndent(session.AddedModels, "", "  ")
				if err != nil {
					content = err.Error()
				} else {
					content = string(payload)
				}
			}
		case "console.3":
			title = "ActiveData"

			if session.ActiveData != nil {
				payload, err := json.MarshalIndent(session.ActiveData, "", "  ")
				if err != nil {
					content = err.Error()
				} else {
					content = string(payload)
				}
			}
		case "console.4":
			title = "Data"

			if session.Data != nil {
				payload, err := json.MarshalIndent(session.Data, "", "  ")
				if err != nil {
					content = err.Error()
				} else {
					content = string(payload)
				}
			}
		case "console.5":
			title = "MemberId"

			if session.MemberId != nil {
				content = fmt.Sprint(*session.MemberId)
			}
		case "console.6":
			title = "SessionId"

			content = session.SessionId
		case "console.7":
			title = "PhoneNumber"

			content = session.PhoneNumber

		case "console.8":
			title = "SQL Query"

			result, err := DB.SQLQuery(text)
			if err != nil {
				content = fmt.Sprintf(`query: %s

response: %s`, text, err.Error())
			} else {
				payload, err := json.MarshalIndent(result, "", " ")
				if err != nil {
					content = fmt.Sprintf(`query: %s

response: %s`, text, err.Error())
				} else {
					content = fmt.Sprintf(`query: %s

response: %s`, text, payload)
				}
			}

		case "console.9":
			title = "Member By PhoneNumber"

			if text != "" {
				result, err := DB.MemberByPhoneNumber(text, nil, nil)
				if err != nil {
					content = fmt.Sprintf("%s\n", err.Error())
				} else {
					payload, err := json.MarshalIndent(result, "", " ")
					if err != nil {
						content = fmt.Sprintf("%s\n", err.Error())
					} else {
						content = fmt.Sprintf("%s\n", payload)
					}
				}
			} else {
				content = ""
			}

		default:
			session.CurrentMenu = "console"

			content = "Available Dumps:\n" +
				"1. WorkflowsMapping\n" +
				"2. AddedModels\n" +
				"3. ActiveData\n" +
				"4. Data\n" +
				"5. MemberId\n" +
				"6. SessionId\n" +
				"7. PhoneNumber\n" +
				"8. SQL Query\n" +
				"9. Member By PhoneNumber"
		}
	} else {
		content = "No active session provided"
	}

	response = "Dev Console\n\n" +
		title +
		fmt.Sprintf("\n\n%s\n", content) +
		"\n99. Cancel\n" +
		"00. Main Menu\n"

	return response
}

func (m *Menus) memberLoansSummary(data map[string]any) string {
	var response string = "Member Loans Summary\n\n00. Main Menu"

	_ = data

	return response
}

func (m *Menus) login(data map[string]any) string {
	var session *parser.Session
	var response string
	var phoneNumber, text, preferencesFolder, cacheFolder string

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

	if text == "" {
		response = "Login\n\nEnter username:\n"
	} else {
		if m.LastPrompt == "username" {
			m.Cache["username"] = text

			m.LastPrompt = "password"

			response = "Login\n\nEnter password:\n"
		} else {
			m.Cache["password"] = text

			text = ""

			result, err := DB.SQLQuery(fmt.Sprintf(`SELECT id, password, role FROM user WHERE username = "%v"`, m.Cache["username"]))
			if err != nil {
				response = err.Error()
			} else {
				if len(result) > 0 {
					passHash := fmt.Sprintf("%v", result[0]["password"])

					if utils.CheckPasswordHash(m.Cache["password"], passHash) {
						token := uuid.NewString()
						session.SessionToken = &token

						session.CurrentMenu = "main"

						return m.LoadMenu("main", session, phoneNumber, text, preferencesFolder, cacheFolder)
					} else {
						m.Cache = map[string]string{}
						m.LastPrompt = "username"

						response = "Login\n\nEnter username:\n"
					}
				} else {
					m.Cache = map[string]string{}
					m.LastPrompt = "username"

					response = "Login\n\nEnter username:\n"
				}
			}
		}
	}

	return response
}
