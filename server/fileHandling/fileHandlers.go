package filehandling

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sacco/server/database"
	"sacco/server/parser"
	"sacco/utils"
	"slices"
	"strconv"

	"github.com/google/uuid"
)

func SaveModelData(data any, model, phoneNumber, cacheFolder *string,
	saveFunc func(map[string]any, string, int) (*int64, error), sessions map[string]*parser.Session, sessionFolder string, refData map[string]any) error {
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
				filename := filepath.Join(sessionFolder, fmt.Sprintf("%s.%s.json", *model, uuid.NewString()))

				transactionDone := false

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

				// By default cache the data first in case we lose database connection
				utils.CacheFile(filename, modelData, 0)
				defer func() {
					if transactionDone || flag.Lookup("test.v") != nil {
						os.Remove(filename)
					}
				}()

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
					log.Println(err)
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
						sessions[*phoneNumber].GlobalIds[fmt.Sprintf("%sId", *model)] = *mid

						sessions[*phoneNumber].AddedModels[*model] = true
					}

					id = *mid

					modelData["id"] = id

					capName := utils.CapitalizeFirstLetter(*model)

					groupSingleName := fmt.Sprintf("%sSingleChildren", capName)
					groupArrayName := fmt.Sprintf("%sArrayChildren", capName)

					models := []string{}
					arrayModels := []string{}
					singleModels := []string{}

					if database.SingleChildren[groupSingleName] != nil {
						models = append(models, database.SingleChildren[groupSingleName]...)

						singleModels = append(singleModels, database.SingleChildren[groupSingleName]...)
					}

					if database.ArrayChildren[groupArrayName] != nil {
						models = append(models, database.ArrayChildren[groupArrayName]...)

						arrayModels = append(arrayModels, database.ArrayChildren[groupArrayName]...)
					}

					for _, childModel := range models {
						childRows, err := utils.CacheDataByModel(childModel, sessionFolder)
						if err != nil {
							log.Println(err)
							continue
						}

						for _, row := range childRows {
							childData, ok := row["data"].(map[string]any)
							if ok {
								filename, ok := row["filename"].(string)
								if ok {
									parentId := fmt.Sprintf("%sId", *model)

									childData[parentId] = id

									lid, err := saveFunc(childData, childModel, 0)
									if err != nil {
										log.Println(err)
										continue
									}

									childKey := fmt.Sprintf("%sId", childModel)
									sessions[*phoneNumber].GlobalIds[childKey] = *lid

									if sessions[*phoneNumber].ActiveData == nil {
										sessions[*phoneNumber].ActiveData = map[string]any{}
									}

									if slices.Contains(arrayModels, childModel) {
										if sessions[*phoneNumber].ActiveData[childModel] == nil {
											sessions[*phoneNumber].ActiveData[childModel] = []map[string]any{}
										}

										sessions[*phoneNumber].ActiveData[childModel] = append(sessions[*phoneNumber].ActiveData[childModel].([]map[string]any), childData)
									} else if slices.Contains(singleModels, childModel) {
										sessions[*phoneNumber].ActiveData[childModel] = childData
									}

									if os.Getenv("DEBUG") != "true" {
										os.Remove(filepath.Join(sessionFolder, filename))
									}
								}
							}
						}
					}

					transactionDone = true
				}
			}

			if sessions != nil && sessions[*phoneNumber] != nil {
				sessions[*phoneNumber].UpdateActiveData(modelData, 0)

				sessions[*phoneNumber].RefreshSession()

				sessions[*phoneNumber].LoadCacheData(*phoneNumber, *cacheFolder)
			}
		}
	}

	return nil
}

func SaveData(
	data any, model, phoneNumber, cacheFolder, preferenceFolder *string,
	saveFunc func(
		map[string]any,
		string,
		int,
	) (*int64, error), sessions map[string]*parser.Session, refData map[string]any) error {
	var sessionFolder string

	if cacheFolder != nil {
		sessionFolder = filepath.Join(*cacheFolder, *phoneNumber)

		_, err := os.Stat(sessionFolder)
		if os.IsNotExist(err) {
			os.MkdirAll(sessionFolder, 0755)
		}
	}

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
		return SaveModelData(data, model, phoneNumber, cacheFolder, saveFunc, sessions, sessionFolder, refData)
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
