package utils

import "regexp"

func CleanScript(content []byte) string {
	stage1 := regexp.MustCompile(`\n|\r`).ReplaceAllLiteralString(string(content), " ")

	return regexp.MustCompile(`\s+`).ReplaceAllLiteralString(stage1, " ")
}
