package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sacco/server/database/models/model2workflow"
	"sacco/utils"
)

func main() {
	workingFolder := filepath.Join(".", "consolidated")

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

	for model := range data {
		targetFile := filepath.Join(workingFolder, fmt.Sprintf("%s.yml", model))

		_, err := model2workflow.Main(model, targetFile, data)
		if err != nil {
			log.Panic(err)
		}
	}
}
