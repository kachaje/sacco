package model2workflow_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sacco/server/database/models/model2workflow"
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
		if false {
			os.RemoveAll(workingFolder)
		}
	}()

	err = model2workflow.Main(model, srcFile, dstFile)
	if err != nil {
		t.Fatal(err)
	}
}
