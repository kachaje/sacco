package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func flattenRecursive(m map[string]any, prefix string, flat map[string]any) {
	for key, value := range m {
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]any:
			flattenRecursive(v, newKey, flat)
		case []map[string]any, []any:
			arr := []map[string]any{}

			if vv, ok := v.([]any); ok {
				for _, vi := range vv {
					if vc, ok := vi.(map[string]any); ok {
						arr = append(arr, vc)
					}
				}
			} else if vv, ok := v.([]map[string]any); ok {
				arr = vv
			}

			for i, vc := range arr {
				newKey = fmt.Sprintf("%s.%v", prefix+"."+key, i)

				flattenRecursive(vc, newKey, flat)
				break
			}
		default:
			if strings.HasSuffix(newKey, ".id") {
				re := regexp.MustCompile(`([A-Z-a-z]+)\.*0*\.id$`)

				if re.MatchString(newKey) {
					model := re.FindAllStringSubmatch(newKey, -1)[0][1]

					flat[model+"Id"] = newKey
				}
			}
		}
	}
}

func FlattenMap(m map[string]any) map[string]any {
	flat := make(map[string]any)
	flattenRecursive(m, "", flat)
	return flat
}

func FlattenKeys(rawData any, seed map[string]any, parent *string) map[string]any {
	handleArrayValues := func(value any, seed map[string]any, parent *string) map[string]any {
		rows := []map[string]any{}

		if val, ok := value.([]map[string]any); ok {
			rows = val
		} else if val, ok := value.([]any); ok {
			for _, row := range val {
				if v, ok := row.(map[string]any); ok {
					rows = append(rows, v)
				}
			}
		}

		if parent != nil {
			for i, row := range rows {
				refKey := fmt.Sprintf("%s.%d", *parent, i)

				seed = FlattenKeys(row, seed, &refKey)
			}
		}

		return seed
	}

	var refKey string

	if data, ok := rawData.(map[string]any); ok {
		for key, value := range data {
			if value == nil {
				continue
			}

			if reflect.TypeOf(value).String() == "map[string]interface {}" {
				if val, ok := value.(map[string]any); ok {
					for k, v := range val {
						if v == nil {
							continue
						}

						if slices.Contains([]string{"[]map[string]interface {}", "[]interface {}", "map[string]interface {}"}, reflect.TypeOf(v).String()) {
							if parent != nil {
								refKey = fmt.Sprintf("%s.%s.%s", *parent, key, k)
							} else {
								refKey = fmt.Sprintf("%s.%s", key, k)
							}

							FlattenKeys(v, seed, &refKey)
						} else {
							if k == "id" {
								refKey = fmt.Sprintf("%sId", key)
							} else {
								refKey = fmt.Sprintf("%s.%s", key, k)
							}

							seed[refKey] = v
						}
					}
				}
			} else if slices.Contains([]string{"[]map[string]interface {}", "[]interface {}"}, reflect.TypeOf(value).String()) {
				if parent != nil {
					refKey = fmt.Sprintf("%s.%s", *parent, key)
				} else {
					refKey = key
				}

				handleArrayValues(value, seed, &refKey)
			} else {
				if parent != nil {
					refKey = fmt.Sprintf("%s.%s", *parent, key)
				} else {
					refKey = key
				}

				seed[refKey] = value
			}
		}
	} else if slices.Contains([]string{"[]map[string]interface {}", "[]interface {}"}, reflect.TypeOf(rawData).String()) && rawData != nil {
		handleArrayValues(rawData, seed, parent)
	}

	return seed
}

func DecodeKey(keyPath string, data map[string]any) (any, bool) {
	var currentMap any = data
	keys := strings.Split(keyPath, ".")

	for i, key := range keys {
		var value map[string]any

		if val, ok := currentMap.(map[string]any); ok {
			value = val
		} else if slices.Contains([]string{"[]map[string]interface {}", "[]interface {}"}, reflect.TypeOf(currentMap).String()) {
			target := []map[string]any{}

			if val, ok := currentMap.([]map[string]any); ok {
				target = val
			} else if val, ok := currentMap.([]any); ok {
				for _, child := range val {
					if v, ok := child.(map[string]any); ok {
						target = append(target, v)
					}
				}
			}

			index, err := strconv.Atoi(key)
			if err == nil {
				currentMap = target[index]
			}
			continue
		}

		if val, ok := value[key]; ok {
			if i == len(keys)-1 {
				return val, true
			}

			if nestedMap, isMap := val.(map[string]any); isMap {
				currentMap = nestedMap
			} else if slices.Contains([]string{"[]map[string]interface {}", "[]interface {}"}, reflect.TypeOf(val).String()) {
				target := []map[string]any{}

				if vt, ok := val.([]map[string]any); ok {
					target = vt
				} else if vt, ok := val.([]any); ok {
					for _, child := range vt {
						if v, ok := child.(map[string]any); ok {
							target = append(target, v)
						}
					}
				}

				currentMap = target
			} else {
				return nil, false
			}
		} else {
			if strings.HasSuffix(key, "Id") {
				rootKey := key[:len(key)-2]

				if val, ok := value[rootKey].(map[string]any)["id"]; ok {
					if i == len(keys)-1 {
						return val, true
					}
				}
			} else if reflect.TypeOf(value).String() == "map[string]interface {}" {
				continue
			}

			return nil, false
		}
	}

	return nil, false
}

func LoadKeys(rawData any, seed map[string]any, parent *string) map[string]any {
	if seed == nil {
		seed = map[string]any{}
	}

	handleArrayValues := func(value any, seed map[string]any, parent *string) map[string]any {
		rows := []map[string]any{}

		if val, ok := value.([]map[string]any); ok {
			rows = val
		} else if val, ok := value.([]any); ok {
			for _, row := range val {
				if v, ok := row.(map[string]any); ok {
					rows = append(rows, v)
				}
			}
		}

		if parent != nil {
			seed[*parent] = []map[string]any{}

			for _, row := range rows {
				result := LoadKeys(row, map[string]any{}, parent)

				if len(result) > 0 {
					seed[*parent] = append(seed[*parent].([]map[string]any), result)
				}
			}
		}

		return seed
	}

	if data, ok := rawData.(map[string]any); ok {
		for key, value := range data {
			if value == nil {
				continue
			}

			if key == "id" {
				if parent != nil {
					seed[fmt.Sprintf("%vId", *parent)] = fmt.Sprintf("%v", value)
				} else {
					seed[key] = fmt.Sprintf("%v", value)
				}
			} else if reflect.TypeOf(value).String() == "map[string]interface {}" {
				if val, ok := value.(map[string]any); ok {
					for k, v := range val {
						if v == nil {
							continue
						}

						if k == "id" {
							seed[fmt.Sprintf("%vId", key)] = fmt.Sprintf("%v", v)
						} else {
							seed = LoadKeys(v, seed, &k)
						}
					}
				}
			} else if slices.Contains([]string{"[]map[string]interface {}", "[]interface {}"}, reflect.TypeOf(value).String()) {
				seed = handleArrayValues(value, seed, &key)
			}
		}
	} else if slices.Contains([]string{"[]map[string]interface {}", "[]interface {}"}, reflect.TypeOf(rawData).String()) && rawData != nil {
		seed = handleArrayValues(rawData, seed, parent)
	}

	return seed
}
