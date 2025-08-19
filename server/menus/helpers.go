package menus

import (
	"fmt"
	"reflect"
	"strconv"
)

type DiffResult struct {
	Added   map[string]any
	Removed map[string]any
	Changed map[string]any
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

			switch level {
			case "beneficiaries":
				v, ok := data[level].([]any)
				if ok {
					var j float64 = 1
					for i := range v {
						vd, ok := v[i].(map[string]any)
						if ok {
							for _, field := range []string{"name", "percentage", "contact"} {
								vf, ok := fieldData[field].(map[string]any)
								if ok {
									keyLabel := fmt.Sprintf("%s%d", field, i+1)

									result[key].(map[string]any)[keyLabel] = map[string]any{
										"order": j,
										"label": fmt.Sprintf("%v", vf["label"]),
									}

									if vd[field] != nil {
										result[key].(map[string]any)[keyLabel].(map[string]any)["value"] = vd[field]
									}

									j++
								}
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

	return result
}
