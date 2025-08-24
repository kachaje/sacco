package model2workflow_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sacco/server/database/models/model2workflow"
	"sacco/utils"
	"testing"
)

func TestModel2Workflow(t *testing.T) {
	workingFolder := filepath.Join(".", "tmpM2W")

	model := "member"
	srcFile := filepath.Join(".", "fixtures", "models.yml")
	dstFile := filepath.Join(workingFolder, fmt.Sprintf("%s.yml", model))

	err := os.MkdirAll(workingFolder, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(workingFolder)
	}()

	content, err := os.ReadFile(srcFile)
	if err != nil {
		t.Fatal(err)
	}

	data, err := utils.LoadYaml(string(content))
	if err != nil {
		t.Fatal(err)
	}

	result, err := model2workflow.Main(model, dstFile, data)
	if err != nil {
		t.Fatal(err)
	}

	target, err := os.ReadFile(filepath.Join(".", "fixtures", "member.yml"))
	if err != nil {
		t.Fatal(err)
	}

	if utils.CleanString(*result) != utils.CleanScript(target) {
		t.Fatal("Test failed")
	}
}
