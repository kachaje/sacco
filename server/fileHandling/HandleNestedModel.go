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

func HandleNestedModel(data any, model, phoneNumber, cacheFolder *string,
	saveFunc func(map[string]any, string, int) (*int64, error), sessions map[string]*parser.Session, sessionFolder string) error {
	if modelData, ok := data.(map[string]any); ok {
		var id int64

		_ = id

		if phoneNumber != nil && *phoneNumber != "default" {
			if modelData["phoneNumber"] == nil {
				modelData["phoneNumber"] = *phoneNumber
			}
		}

		if sessions[*phoneNumber] != nil && model != nil {
			filename := filepath.Join(sessionFolder, fmt.Sprintf("%s.json", *model))

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
			CacheFile(filename, modelData, 0)
			defer func() {
				if transactionDone {
					os.Remove(filename)
				}
			}()

			someChildAdded := false

			capName := utils.CapitalizeFirstLetter(*model)

			groupSingleName := fmt.Sprintf("%sSingleChildren", capName)
			groupArrayName := fmt.Sprintf("%sArrayChildren", capName)

			models := []string{}

			if database.SingleChildren[groupSingleName] != nil {
				models = append(models, database.SingleChildren[groupSingleName]...)
			}

			if database.ArrayChildren[groupArrayName] != nil {
				models = append(models, database.ArrayChildren[groupArrayName]...)
			}

			for _, key := range models {
				if sessions[*phoneNumber].AddedModels[key] {
					someChildAdded = true
					break
				}
			}

			if saveFunc == nil {
				return fmt.Errorf("server.HandleNestedModel.%s:missing saveFunc", *model)
			}

			if database.ParentModels[*model] != nil {
				for _, value := range database.ParentModels[*model] {
					key := fmt.Sprintf("%sId", value)
					if sessions[*phoneNumber].GlobalIds[key] > 0 {
						modelData[key] = sessions[*phoneNumber].GlobalIds[key]
					}
				}
			}

			mid, err := saveFunc(modelData, *model, 0)
			if err != nil {
				log.Println(err)
				return err
			}

			sessions[*phoneNumber].GlobalIds[fmt.Sprintf("%sId", *model)] = *mid

			id = *mid

			modelData["id"] = id

			if someChildAdded {
				for _, childModel := range models {
					file := filepath.Join(sessionFolder, fmt.Sprintf("%s.json", childModel))

					_, err = os.Stat(file)
					if !os.IsNotExist(err) {
						fileArrayData := []map[string]any{}
						fileFlatData := map[string]any{}

						content, err := os.ReadFile(file)
						if err != nil {
							log.Printf("server.HandleNestedModel:%s\n", err.Error())
						} else {
							err = json.Unmarshal(content, &fileArrayData)
							if err != nil {
								err = json.Unmarshal(content, &fileFlatData)
								if err != nil {
									log.Printf("server.HandleNestedModel:%s\n", err.Error())
								}
							}
						}

						if len(fileArrayData) > 0 {
							for i := range fileArrayData {
								fileArrayData[i][fmt.Sprintf("%sId", *model)] = id

								mid, err = saveFunc(fileArrayData[i], childModel, 0)
								if err != nil {
									log.Printf("server.HandleNestedModel:%s\n", err.Error())
									continue
								}

								fileArrayData[i]["id"] = *mid
							}

							modelData[childModel] = fileArrayData

							if os.Getenv("DEBUG") == "true" {
								CacheFile(file, fileArrayData, 0)
							} else {
								os.Remove(file)
							}
						} else if len(fileFlatData) > 0 {
							fileFlatData[fmt.Sprintf("%sId", *model)] = id

							mid, err = saveFunc(fileFlatData, childModel, 0)
							if err != nil {
								log.Printf("server.HandleNestedModel:%s\n", err.Error())
								continue
							}

							fileFlatData["id"] = *mid

							modelData[childModel] = fileFlatData

							if os.Getenv("DEBUG") == "true" {
								CacheFile(file, fileFlatData, 0)
							} else {
								os.Remove(file)
							}
						}

						sessions[*phoneNumber].GlobalIds[fmt.Sprintf("%sId", childModel)] = *mid
					}
				}
			}

			transactionDone = true
		}

		sessions[*phoneNumber].UpdateActiveData(modelData, 0)

		sessions[*phoneNumber].RefreshSession()

		sessions[*phoneNumber].LoadCacheData(*phoneNumber, *cacheFolder)
	}

	return nil
}
