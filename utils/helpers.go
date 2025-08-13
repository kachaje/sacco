package utils

import (
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func CleanScript(content []byte) string {
	stage1 := regexp.MustCompile(`\n|\r`).ReplaceAllLiteralString(string(content), " ")

	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllLiteralString(stage1, " "))
}

func CleanString(content string) string {
	stage1 := regexp.MustCompile(`\n|\r`).ReplaceAllLiteralString(string(content), " ")

	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllLiteralString(stage1, " "))
}

func LoadYaml(yamlData string) (map[string]any, error) {
	var data map[string]any

	err := yaml.Unmarshal([]byte(yamlData), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
