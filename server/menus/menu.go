package menus

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
	menufuncs "sacco/server/menus/menuFuncs"
	"sacco/server/parser"
	"sacco/utils"
	"slices"
	"strings"
	"sync"
)

//go:embed menus/*
var menuFiles embed.FS

var funcsMap map[string]func(map[string]any) string

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

func init() {
	funcsMap = map[string]func(map[string]any) string{}

	fmt.Println(funcsMap)
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
		return menufuncs.DoExit(m.LoadMenu, data)
	}
	m.FunctionsMap["businessSummary"] = func(data map[string]any) string {
		return menufuncs.BusinessSummary(m.LoadMenu, data)
	}
	m.FunctionsMap["employmentSummary"] = func(data map[string]any) string {
		return menufuncs.EmploymentSummary(m.LoadMenu, data)
	}
	m.FunctionsMap["checkBalance"] = func(data map[string]any) string {
		return menufuncs.CheckBalance(m.LoadMenu, data)
	}
	m.FunctionsMap["bankingDetails"] = func(data map[string]any) string {
		return menufuncs.BankingDetails(m.LoadMenu, data)
	}
	m.FunctionsMap["viewMemberDetails"] = func(data map[string]any) string {
		return menufuncs.ViewMemberDetails(m.LoadMenu, data)
	}
	m.FunctionsMap["devConsole"] = func(data map[string]any) string {
		return menufuncs.DevConsole(m.LoadMenu, data)
	}
	m.FunctionsMap["memberLoansSummary"] = func(data map[string]any) string {
		return menufuncs.MemberLoansSummary(m.LoadMenu, data)
	}
	m.FunctionsMap["signIn"] = func(data map[string]any) string {
		return menufuncs.SignIn(m.LoadMenu, data)
	}
	m.FunctionsMap["listUsers"] = func(data map[string]any) string {
		return menufuncs.ListUsers(m.LoadMenu, data)
	}
	m.FunctionsMap["blockUser"] = func(data map[string]any) string {
		return menufuncs.BlockUser(m.LoadMenu, data)
	}
	m.FunctionsMap["editUser"] = func(data map[string]any) string {
		return menufuncs.EditUser(m.LoadMenu, data)
	}
	m.FunctionsMap["changePassword"] = func(data map[string]any) string {
		return menufuncs.ChangePassword(m.LoadMenu, data)
	}
	m.FunctionsMap["signUp"] = func(data map[string]any) string {
		return menufuncs.SignUp(m.LoadMenu, data)
	}
	m.FunctionsMap["landing"] = func(data map[string]any) string {
		return menufuncs.Landing(m.LoadMenu, data)
	}

	err := m.populateMenus()
	if err != nil {
		log.Panic(err)
	}

	return m
}

func (m *Menus) populateMenus() error {
	return fs.WalkDir(menuFiles, ".", func(file string, d fs.DirEntry, err error) error {
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
}

func (m *Menus) LoadMenu(menuName string, session *parser.Session, phoneNumber, text, preferencesFolder, cacheFolder string) string {
	var response string

	preferredLanguage := menufuncs.CheckPreferredLanguage(phoneNumber, preferencesFolder)

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
				m.LoadMenu,
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
				m.LoadMenu,
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
				m.LoadMenu,
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
