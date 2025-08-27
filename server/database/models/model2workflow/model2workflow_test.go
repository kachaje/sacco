package model2workflow_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sacco/server/database/models/model2workflow"
	"sacco/utils"
	"testing"
)

func TestModel2WorkflowBasic(t *testing.T) {
	workingFolder := filepath.Join(".", "tmpM2WBasic")

	model := "member"
	srcFile := filepath.Join("..", "..", "models.yml")
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

	result, _, err := model2workflow.Main(model, dstFile, data)
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

func TestModel2WorkflowComplex(t *testing.T) {
	workingFolder := filepath.Join(".", "tmpM2WComplex")

	model := "memberBeneficiary"
	srcFile := filepath.Join("..", "..", "models.yml")
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

	result, _, err := model2workflow.Main(model, dstFile, data)
	if err != nil {
		t.Fatal(err)
	}

	target, err := os.ReadFile(filepath.Join(".", "fixtures", "memberBeneficiary.yml"))
	if err != nil {
		t.Fatal(err)
	}

	if utils.CleanString(*result) != utils.CleanScript(target) {
		t.Fatal("Test failed")
	}
}
