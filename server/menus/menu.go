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
	"strings"
)

//go:embed menus/*
var menuFiles embed.FS

var menuFilesData map[string]any

type Menus struct {
	ActiveMenus map[string]any
}

func NewMenus() *Menus {
	m := &Menus{
		ActiveMenus: map[string]any{},
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

		m.ActiveMenus[group] = data

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return m
}

func (m *Menus) LoadMenu(menuName string, session *parser.Session, phoneNumber, text, preferencesFolder, cacheFolder string) string {
	var response string

	payload, err := json.MarshalIndent(m.ActiveMenus, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(payload))

	return response
}
