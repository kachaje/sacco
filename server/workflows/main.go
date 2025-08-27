package main

import (
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"sacco/server/database/models/model2workflow"
	"sacco/utils"
	"strings"
)

func main() {
	workingFolder := filepath.Join(".")

	_, err := os.Stat(workingFolder)
	if os.IsNotExist(err) {
		err = os.MkdirAll(workingFolder, 0755)
		if err != nil {
			log.Panic(err)
		}
	} else if err != nil {
		log.Panic(err)
	}

	content, err := os.ReadFile(filepath.Join("..", "database", "models.yml"))
	if err != nil {
		log.Panic(err)
	}

	data, err := utils.LoadYaml(string(content))
	if err != nil {
		log.Panic(err)
	}

	relationships := map[string]any{}

	for model := range data {
		targetFile := filepath.Join(workingFolder, fmt.Sprintf("%s.yml", model))

		_, row, err := model2workflow.Main(model, targetFile, data)
		if err != nil {
			log.Panic(err)
		}

		if len(row) > 0 {
			relationships[model] = row
		}
	}

	script := []string{}

	for key, value := range relationships {
		model := utils.CapitalizeFirstLetter(key)

		if val, ok := value.(map[string][]string); ok {
			if len(val["hasMany"]) > 0 {
				rows := []string{}

				for _, v := range val["hasMany"] {
					rows = append(rows, fmt.Sprintf(`"%s"`, v))
				}

				row := strings.TrimSpace(fmt.Sprintf(`%sArrayChildren = []string{
				%s,
				}`, model, strings.Join(rows, ",\n")))

				if len(row) > 0 {
					script = append(script, row)
				}
			}
			if len(val["hasOne"]) > 0 {
				rows := []string{}

				for _, v := range val["hasOne"] {
					rows = append(rows, fmt.Sprintf(`"%s"`, v))
				}

				row := strings.TrimSpace(fmt.Sprintf(`%sSingleChildren = []string{
				%s,
				}
				`, model, strings.Join(rows, ",\n")))

				if len(row) > 0 {
					script = append(script, row)
				}
			}
		}
	}

	targetName := filepath.Join("..", "database", "models.go")

	content, err = format.Source(fmt.Appendf(nil, `package database
	
	var (
	%s
	)
	`, strings.Join(script, "\n")))
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(targetName, content, 0644)
	if err != nil {
		panic(err)
	}
}
