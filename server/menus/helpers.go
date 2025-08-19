package menus

import (
	"fmt"
	"strconv"
)

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
					j := 1
					for i := range v {
						vd, ok := v[i].(map[string]any)
						if ok {
							for field, kids := range fieldData {
								vf, ok := kids.(map[string]any)
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
						order, err := strconv.ParseInt(fmt.Sprintf("%v", vf["order"]), 10, 64)
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
