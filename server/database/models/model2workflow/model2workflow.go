package model2workflow

import (
	"encoding/json"
	"fmt"
	"os"
	"sacco/utils"
)

func Main(model, sourceFile, destinationFile string) error {
	content, err := os.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	data, err := utils.LoadYaml(string(content))
	if err != nil {
		return err
	}

	payload, _ := json.MarshalIndent(data[model], "", "  ")

	fmt.Println(string(payload))

	return nil
}
