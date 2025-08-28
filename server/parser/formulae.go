package parser

import (
	"regexp"
)

func GetTokens(query string) map[string]any {
	result := map[string]any{}

	re := regexp.MustCompile(`^([A-z]+)`)

	op := re.FindAllString(query, -1)[0]

	result["op"] = op

	re = regexp.MustCompile(`([A-Za-z]+)`)

	result["terms"] = []string{}

	for _, term := range re.FindAllString(query, -1) {
		if term != op {
			result["terms"] = append(result["terms"].([]string), term)
		}
	}

	return result
}
