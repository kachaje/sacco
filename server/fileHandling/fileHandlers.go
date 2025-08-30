package filehandling

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sacco/server/database"
	"sacco/server/parser"
	"sacco/utils"
	"strconv"
)

func SaveModelData(data any, model, phoneNumber *string,
	saveFunc func(map[string]any, string, int) (*int64, error), sessions map[string]*parser.Session, refData map[string]any) error {
	if rawData, ok := data.(map[string]any); ok {
		dataRows := utils.UnpackData(rawData)

		if refData != nil {
			unpackedRefData := utils.UnpackData(refData)

			missingIds := utils.GetSkippedRefIds(dataRows, unpackedRefData)

			for _, row := range missingIds {
				row["active"] = 0

				dataRows = append(dataRows, row)
			}
		}

		for _, modelData := range dataRows {
			if model != nil {
				for _, key := range database.FloatFields {
					if modelData[key] != nil {
						nv, ok := modelData[key].(string)
						if ok {
							real, err := strconv.ParseFloat(nv, 64)
							if err == nil {
								modelData[key] = real
							}
						}
					}
				}

				if saveFunc == nil {
					return fmt.Errorf("server.SaveModelData.%s:missing saveFunc", *model)
				}

				if sessions[*phoneNumber] != nil {
					if sessions[*phoneNumber].GlobalIds == nil {
						sessions[*phoneNumber].GlobalIds = map[string]int64{}
					}
					if sessions[*phoneNumber].AddedModels == nil {
						sessions[*phoneNumber].AddedModels = map[string]bool{}
					}

					if database.ParentModels[*model] != nil {
						for _, value := range database.ParentModels[*model] {
							key := fmt.Sprintf("%sId", value)
							if sessions[*phoneNumber].GlobalIds[key] > 0 {
								modelData[key] = sessions[*phoneNumber].GlobalIds[key]
							}
						}
					}
				}

				mid, err := saveFunc(modelData, *model, 0)
				if err != nil {
					return err
				}

				var id int64

				if mid == nil && modelData["id"] != nil {
					val, err := strconv.ParseInt(fmt.Sprintf("%v", modelData["id"]), 10, 64)
					if err == nil {
						mid = &val
					}
				}
				if mid != nil {
					if sessions[*phoneNumber] != nil {
						if len(dataRows) == 1 {
							sessions[*phoneNumber].GlobalIds[fmt.Sprintf("%sId", *model)] = *mid
						}

						sessions[*phoneNumber].AddedModels[*model] = true
					}

					id = *mid

					modelData["id"] = id
				}
			}

			if sessions != nil && sessions[*phoneNumber] != nil {
				sessions[*phoneNumber].UpdateActiveData(modelData, 0)

				sessions[*phoneNumber].RefreshSession()
			}
		}
	}

	return nil
}

func SaveData(
	data any, model, phoneNumber, preferenceFolder *string,
	saveFunc func(
		map[string]any,
		string,
		int,
	) (*int64, error), sessions map[string]*parser.Session, refData map[string]any) error {
	switch *model {
	case "language":
		val, ok := data.(map[string]any)
		if ok {
			if val["language"] != nil && phoneNumber != nil {
				language, ok := val["language"].(string)
				if ok {
					SavePreference(*phoneNumber, "language", language, *preferenceFolder)
				}
			}
		}

	default:
		return SaveModelData(data, model, phoneNumber, saveFunc, sessions, refData)
	}

	return nil
}

func SavePreference(phoneNumber, key, value, preferencesFolder string) error {
	settingsFile := filepath.Join(preferencesFolder, phoneNumber)

	data := map[string]any{}

	_, err := os.Stat(settingsFile)
	if !os.IsNotExist(err) {
		content, err := os.ReadFile(settingsFile)
		if err != nil {
			log.Println(err)
			return err
		}

		err = json.Unmarshal(content, &data)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	data[key] = value

	payload, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
		return err
	}

	return os.WriteFile(settingsFile, payload, 0644)
}
