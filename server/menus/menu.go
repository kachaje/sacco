package menus

import (
	"embed"
	"encoding/json"
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
	ActiveMenus  map[string]any
	Titles       map[string]string
	Workflows    map[string]string
	Functions    map[string]string
	FunctionsMap map[string]func()
}

func NewMenus() *Menus {
	m := &Menus{
		ActiveMenus:  map[string]any{},
		Titles:       map[string]string{},
		Workflows:    map[string]string{},
		Functions:    map[string]string{},
		FunctionsMap: map[string]func(){},
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
							}
						}
						if val["function"] != nil {
							if v, ok := val["function"].(string); ok {
								m.Functions[id] = v
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

	fmt.Println(session.CurrentMenu, text, keys)

	payload, _ := json.MarshalIndent(m.Workflows, "", "  ")

	fmt.Println(string(payload))

	if false {
		payload, _ = json.MarshalIndent(m.Functions, "", "  ")

		fmt.Println(string(payload))
	}

	if slices.Contains(keys, text) {
		target := text
		text = "000"

		if session != nil {
			session.CurrentMenu = kv[target]
		}

		return m.LoadMenu(kv[target], session, phoneNumber, text, preferencesFolder, cacheFolder)
	} else if m.Workflows[session.CurrentMenu] != "" {

		fmt.Println("##########", m.Workflows[session.CurrentMenu])

	} else {
		response = fmt.Sprintf("CON %s\n%s", m.Titles[menuName], strings.Join(values, ""))
	}

	return response
}
