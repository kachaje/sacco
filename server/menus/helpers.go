package menus

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
)

type DiffResult struct {
	Added   map[string]any
	Removed map[string]any
	Changed map[string]any
}

func CheckPreferredLanguage(phoneNumber, preferencesFolder string) *string {
	settingsFile := filepath.Join(preferencesFolder, phoneNumber)

	_, err := os.Stat(settingsFile)
	if !os.IsNotExist(err) {
		content, err := os.ReadFile(settingsFile)
		if err != nil {
			log.Println(err)
			return nil
		}

		data := map[string]any{}

		err = json.Unmarshal(content, &data)
		if err != nil {
			log.Println(err)
			return nil
		}

		var preferredLanguage string

		if data["language"] != nil {
			val, ok := data["language"].(string)
			if ok {
				preferredLanguage = val
			}
		}

		return &preferredLanguage
	}

	return nil
}

func GetMapDiff(map1, map2 map[string]any) DiffResult {
	diff := DiffResult{
		Added:   make(map[string]any),
		Removed: make(map[string]any),
		Changed: make(map[string]any),
	}

	for key, val1 := range map1 {
		if val2, ok := map2[key]; !ok {
			diff.Removed[key] = val1
		} else {
			if nestedMap1, isMap1 := val1.(map[string]any); isMap1 {
				if nestedMap2, isMap2 := val2.(map[string]any); isMap2 {
					nestedDiff := GetMapDiff(nestedMap1, nestedMap2)
					if len(nestedDiff.Added) > 0 || len(nestedDiff.Removed) > 0 || len(nestedDiff.Changed) > 0 {
						diff.Changed[key] = nestedDiff
					}
				} else {
					diff.Changed[key] = map[string]any{
						"old":     val1,
						"new":     val2,
						"oldType": reflect.TypeOf(val1).String(),
						"newType": reflect.TypeOf(val2).String(),
					}
				}
			} else if !reflect.DeepEqual(val1, val2) {
				diff.Changed[key] = map[string]any{
					"old":     val1,
					"new":     val2,
					"oldType": reflect.TypeOf(val1).String(),
					"newType": reflect.TypeOf(val2).String(),
				}
			}
		}
	}

	for key, val2 := range map2 {
		if _, ok := map1[key]; !ok {
			diff.Added[key] = val2
		}
	}

	return diff
}

func LoadTemplateData(data map[string]any, template map[string]any) map[string]any {
	result := map[string]any{}

	for key, value := range template {
		result[key] = map[string]any{}

		var level string

		val, ok := value.(map[string]any)
		if ok {
			if val["level"] != nil {
				level = fmt.Sprintf("%v", val["level"])
			}

			fieldData := map[string]any{}

			v, ok := val["data"].(map[string]any)
			if ok {
				fieldData = v
			}

			loadBeneficiary := func(vd map[string]any, i int, j *float64) {
				for _, field := range []string{"name", "percentage", "contact"} {
					vf, ok := fieldData[field].(map[string]any)
					if ok {
						keyLabel := fmt.Sprintf("%s%d", field, i+1)

						result[key].(map[string]any)[keyLabel] = map[string]any{
							"order": *j,
							"label": fmt.Sprintf("%v", vf["label"]),
						}

						if vd[field] != nil {
							result[key].(map[string]any)[keyLabel].(map[string]any)["value"] = vd[field]
						}

						*j++
					}
				}
			}

			switch level {
			case "memberBeneficiary":
				v, ok := data[level].([]any)
				if ok {
					var j float64 = 1
					for i := range v {
						vd, ok := v[i].(map[string]any)
						if ok {
							loadBeneficiary(vd, i, &j)
						}
					}
				} else {
					v, ok := data[level].([]map[string]any)
					if ok {
						var j float64 = 1
						for i := range v {
							vd := v[i]
							if ok {
								loadBeneficiary(vd, i, &j)
							}
						}
					}
				}
			default:
				for field, kids := range fieldData {
					vf, ok := kids.(map[string]any)
					if ok {
						order, err := strconv.ParseFloat(fmt.Sprintf("%v", vf["order"]), 64)
						if err == nil {
							result[key].(map[string]any)[field] = map[string]any{
								"order": order,
								"label": fmt.Sprintf("%v", vf["label"]),
							}

							if level == "root" {
								if data[field] != nil && fmt.Sprintf("%v", data[field]) != "" {
									result[key].(map[string]any)[field].(map[string]any)["value"] = data[field]
								}
							} else {
								if data[level] != nil {
									v, ok := data[level].(map[string]any)
									if ok {
										if v[field] != nil {
											result[key].(map[string]any)[field].(map[string]any)["value"] = v[field]
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return result
}

func TabulateData(data map[string]any) []string {
	result := []string{}

	keys := []string{}

	for key := range data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		result = append(result, key)

		keysMap := map[float64]string{}
		floatKeys := []float64{}

		row, ok := data[key].(map[string]any)
		if ok {
			childData := map[string]map[string]any{}

			for k, v := range row {
				val, ok := v.(map[string]any)
				if ok && val["order"] != nil {
					childData[k] = val

					order, err := strconv.ParseFloat(fmt.Sprintf("%v", val["order"]), 64)
					if err == nil {
						keysMap[order] = k
						floatKeys = append(floatKeys, order)
					}
				}
			}

			sort.Float64s(floatKeys)

			if key == "E. BENEFICIARIES DETAILS" {
				row1 := "--- --------------------- --------- ------------"
				row2 := "No | Name of Beneficiary | Percent | Contact"

				result = append(result, row1)
				result = append(result, row2)
				result = append(result, row1)

				for i := range 4 {
					index := i + 1

					nameLabel := fmt.Sprintf("name%d", index)
					percentageLabel := fmt.Sprintf("percentage%d", index)
					contactLabel := fmt.Sprintf("contact%d", index)

					if childData[nameLabel] == nil {
						break
					}

					var name string
					var percentage float64
					var contact string

					name = fmt.Sprintf("%v", childData[nameLabel]["value"])

					if childData[percentageLabel] != nil {
						v, err := strconv.ParseFloat(fmt.Sprintf("%v", childData[percentageLabel]["value"]), 64)
						if err == nil {
							percentage = v
						}
					}
					if childData[contactLabel] != nil {
						contact = fmt.Sprintf("%v", childData[contactLabel]["value"])
					}

					entry := fmt.Sprintf("%-3d| %-19s | %7.1f | %s", index, name, percentage, contact)

					result = append(result, entry)
				}
			} else {
				for _, order := range floatKeys {
					var label string
					var value string

					childKey := keysMap[order]

					if childData[childKey]["label"] != nil {
						label = fmt.Sprintf("%v:", childData[childKey]["label"])
					}
					if childData[childKey]["value"] != nil {
						value = fmt.Sprintf("%v", childData[childKey]["value"])
					}

					entry := fmt.Sprintf("   %-25s| %s", label, value)

					result = append(result, entry)
				}
			}

			result = append(result, "")
		}
	}

	if false {
		payload, _ := json.MarshalIndent(data, "", "  ")

		fmt.Println(string(payload))
	}

	return result
}
