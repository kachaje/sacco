package filehandling

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sacco/server/database"
	"sacco/server/parser"
	"sacco/utils"
	"slices"
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

				if saveFunc == nil {
					return fmt.Errorf("server.HandleNestedModel.%s:missing saveFunc", *model)
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
						childRows, err := CacheDataByModel(childModel, sessionFolder)
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

func CacheDataByModel(filterModel, sessionFolder string) ([]map[string]any, error) {
	result := []map[string]any{}

	err := filepath.WalkDir(sessionFolder, func(fullpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		filename := filepath.Base(fullpath)

		re := regexp.MustCompile(`\.[a-z0-9-]+\.json$`)

		if !re.MatchString(filename) {
			return nil
		}

		model := re.ReplaceAllLiteralString(filename, "")

		if model != filterModel {
			return nil
		}

		content, err := os.ReadFile(fullpath)
		if err != nil {
			return err
		}

		if data := map[string]any{}; json.Unmarshal(content, &data) == nil {
			result = append(result, map[string]any{
				"data":     data,
				"filename": filename,
			})
		} else if data := []map[string]any{}; json.Unmarshal(content, &data) == nil {
			rows := []map[string]any{}

			for _, row := range data {
				rows = append(rows, map[string]any{
					"data":     row,
					"filename": filename,
				})
			}

			result = append(result, rows...)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
