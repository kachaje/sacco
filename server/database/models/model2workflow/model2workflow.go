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

	sourceData, err := utils.LoadYaml(string(content))
	if err != nil {
		return err
	}

	data := map[string]any{
		"model": model,
		"formSummary": map[string]any{
			"type": "quitScreen",
		},
	}

	j := 0
	lastTag := ""
	for _, row := range sourceData[model].([]any) {
		if val, ok := row.(map[string]any); ok {
			for key, rawValue := range val {
				if value, ok := rawValue.(map[string]any); ok {
					tag := fmt.Sprintf("enter%s", utils.CapitalizeFirstLetter(key))

					if data["initialScreen"] == nil && value["hidden"] == nil {
						data["initialScreen"] = tag
					}

					data[tag] = map[string]any{
						"inputIdentifier": key,
					}

					if value["hidden"] == nil {
						j++

						data[tag].(map[string]any)["order"] = j
						data[tag].(map[string]any)["type"] = "inputScreen"
						data[tag].(map[string]any)["nextScreen"] = "formSummary"

						if lastTag != "" {
							data[lastTag].(map[string]any)["nextScreen"] = tag
						}

						lastTag = tag
					} else {
						data[tag].(map[string]any)["hidden"] = true
						data[tag].(map[string]any)["type"] = "hiddenField"
					}

					_ = value
				}
			}
		}
	}

	payload, _ := json.MarshalIndent(data, "", "  ")

	fmt.Println(string(payload))

	return nil
}
