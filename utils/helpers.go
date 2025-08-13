package utils

import (
	"regexp"
	"strings"
)

func CleanScript(content []byte) string {
	stage1 := regexp.MustCompile(`\n|\r`).ReplaceAllLiteralString(string(content), " ")

	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllLiteralString(stage1, " "))
}

func CleanString(content string) string {
	stage1 := regexp.MustCompile(`\n|\r`).ReplaceAllLiteralString(string(content), " ")

	return strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllLiteralString(stage1, " "))
}
