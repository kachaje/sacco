package model2workflow

import (
	"fmt"
	"os"
	"sacco/utils"
)

func Main(model, destinationFile string, sourceData map[string]any) (*string, error) {
	data := map[string]any{
		"model": model,
		"formSummary": map[string]any{
			"type": "quitScreen",
		},
	}

	j := 0
	lastTag := ""

	if rawData, ok := sourceData[model].(map[string]any); ok {
		for _, row := range rawData["fields"].([]any) {
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

							text := utils.IdentifierToLabel(key)

							data[tag].(map[string]any)["text"] = map[string]any{
								"en": text,
							}

							data[tag].(map[string]any)["order"] = j
							data[tag].(map[string]any)["type"] = "inputScreen"
							data[tag].(map[string]any)["nextScreen"] = "formSummary"

							if value["optional"] != nil {
								data[tag].(map[string]any)["optional"] = true
							}

							if value["numericField"] != nil {
								data[tag].(map[string]any)["validationRule"] = "^\\d+\\.*\\d*$"
							}

							if value["validationRule"] != nil {
								data[tag].(map[string]any)["validationRule"] = value["validationRule"].(string)
							}

							if value["terminateBlockOnEmpty"] != nil {
								data[tag].(map[string]any)["terminateBlockOnEmpty"] = true
							}

							if value["adminOnly"] != nil {
								data[tag].(map[string]any)["adminOnly"] = true
							}

							if value["options"] != nil {
								if opts, ok := value["options"].([]any); ok {
									options := []any{}

									for i, opt := range opts {
										option := map[string]any{
											"position": i + 1,
											"label": map[string]any{
												"en": opt,
											},
										}

										options = append(options, option)
									}

									data[tag].(map[string]any)["options"] = options
								}
							}

							if lastTag != "" {
								data[lastTag].(map[string]any)["nextScreen"] = tag
							}

							lastTag = tag
						} else {
							data[tag].(map[string]any)["hidden"] = true
							data[tag].(map[string]any)["type"] = "hiddenField"
						}
					}
				}
			}
		}
	}

	yamlString, err := utils.DumpYaml(data)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(destinationFile, []byte(*yamlString), 0644)
	if err != nil {
		return nil, err
	}

	return yamlString, nil
}
