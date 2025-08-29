package filehandling

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sacco/server/database"
	"sacco/server/parser"
	"sacco/utils"
	"strconv"

	"github.com/google/uuid"
)

func UnpackData(data map[string]any) []map[string]any {
	result := []map[string]any{}
	rows := map[string]map[string]any{}

	for key, value := range data {
		re := regexp.MustCompile(`^(.+)(\d+)$`)

		if re.MatchString(key) {
			parts := re.FindAllStringSubmatch(key, -1)[0]

			field := parts[1]
			index := parts[2]

			if rows[index] == nil {
				rows[index] = map[string]any{}
			}

			rows[index][field] = value
		} else {
			rows["1"] = data
			break
		}
	}

	for _, row := range rows {
		result = append(result, row)
	}

	return result
}

func GetSkippedRefIds(data, refData []map[string]any) []map[string]any {
	result := []map[string]any{}

	for _, row := range refData {
		if row["id"] != nil {
			id := fmt.Sprintf("%v", row["id"])
			found := false

			for _, child := range data {
				if child["id"] != nil && id == fmt.Sprintf("%v", child["id"]) {
					found = true
					break
				}
			}
			if !found {
				result = append(result, row)
			}
		}
	}

	return result
}

func HandleNestedModel(data any, model, phoneNumber, cacheFolder *string,
	saveFunc func(map[string]any, string, int) (*int64, error), sessions map[string]*parser.Session, sessionFolder string, refData map[string]any) error {
	if rawData, ok := data.(map[string]any); ok {
		dataRows := UnpackData(rawData)

		if refData != nil {
			unpackedRefData := UnpackData(refData)

			missingIds := GetSkippedRefIds(dataRows, unpackedRefData)

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
				CacheFile(filename, modelData, 0)
				defer func() {
					if transactionDone || flag.Lookup("test.v") != nil {
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

				if saveFunc == nil {
					return fmt.Errorf("server.HandleNestedModel.%s:missing saveFunc", *model)
				}

				if sessions[*phoneNumber] != nil {
					for _, key := range models {
						if sessions[*phoneNumber].AddedModels[key] {
							someChildAdded = true
							break
						}
					}

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
				}

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

							if sessions[*phoneNumber] != nil {
								sessions[*phoneNumber].GlobalIds[fmt.Sprintf("%sId", childModel)] = *mid
							}
						}
					}
				}

				transactionDone = true
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
