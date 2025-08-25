package menus

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
	"sacco/server/parser"
	"sacco/utils"
	"slices"
	"strings"
)

//go:embed menus/*
var menuFiles embed.FS

type Menus struct {
	ActiveMenus   map[string]any
	Titles        map[string]string
	Workflows     map[string]any
	Functions     map[string]string
	FunctionsMap  map[string]func()
	TargetKeys    map[string][]string
	LabelWorkflow map[string]any
}

func NewMenus() *Menus {
	m := &Menus{
		ActiveMenus:   map[string]any{},
		Titles:        map[string]string{},
		Workflows:     map[string]any{},
		Functions:     map[string]string{},
		FunctionsMap:  map[string]func(){},
		TargetKeys:    map[string][]string{},
		LabelWorkflow: map[string]any{},
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

			for key, row := range val {
				keys = append(keys, key)

				if val, ok := row.(map[string]any); ok {
					if val["id"] != nil && val["label"] != nil && val["label"].(map[string]any)["en"] != nil {
						id := fmt.Sprintf("%v", val["id"])
						label := fmt.Sprintf("%v", val["label"].(map[string]any)["en"])

						kv[key] = id

						value := fmt.Sprintf("%v. %v\n", key, label)

						values = append(values, value)

						if val["workflow"] != nil {
							if v, ok := val["workflow"].(string); ok {
								m.Workflows[id] = v

								m.LabelWorkflow[group].(map[string]any)[value] = map[string]any{
									"workflow": v,
									"id":       id,
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

	slices.Sort(values)

	index := utils.Index(values, "00. Main Menu\n")

	if index >= 0 {
		values = append(values[:index], values[index+1:]...)

		values = append(values, "\n00. Main Menu\n")
	}

	if slices.Contains(keys, text) {
		target := text
		text = "000"

		if session != nil {
			session.CurrentMenu = kv[target]
		}

		return m.LoadMenu(kv[target], session, phoneNumber, text, preferencesFolder, cacheFolder)
	} else if session != nil &&  m.Workflows[session.CurrentMenu] != nil {
		model := session.CurrentMenu
		workflow := fmt.Sprintf("%v", m.Workflows[model])

		if session.ActiveMemberData != nil {
			data := map[string]any{}

			if regexp.MustCompile(`^\d+$`).MatchString(phoneNumber) {
				data["phoneNumber"] = phoneNumber
			}

			if m.TargetKeys[model] != nil {
				targetKeys := m.TargetKeys[model]

				for key, value := range session.ActiveMemberData {
					if slices.Contains(targetKeys, key) {
						data[key] = fmt.Sprintf("%v", value)
					}
				}

				if session.WorkflowsMapping != nil &&
					session.WorkflowsMapping[workflow] != nil {
					for key, value := range session.WorkflowsMapping[workflow].Data {
						if slices.Contains(targetKeys, key) {
							data[key] = fmt.Sprintf("%v", value)
						}
					}

					session.WorkflowsMapping[workflow].Data = data
				}
			}
		}

		if session.WorkflowsMapping != nil &&
			session.WorkflowsMapping[workflow] != nil {
			response = session.WorkflowsMapping[workflow].NavNext(text)

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
			response = "NOT IMPLEMENTED YET"
		}

	} else {
		newValues := []string{}

		if m.LabelWorkflow[menuName] != nil && session != nil {
			for _, value := range values {
				if m.LabelWorkflow[menuName].(map[string]any)[value] != nil {
					workflow := fmt.Sprintf("%v", m.LabelWorkflow[menuName].(map[string]any)[value].(map[string]any)["workflow"])

					suffix := ""

					if session.AddedModels[workflow] {
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

		response = fmt.Sprintf("CON %s\n%s", m.Titles[menuName], strings.Join(newValues, ""))
	}

	return response
}
