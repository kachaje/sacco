package menus_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sacco/server/menus"
	"sacco/utils"
	"strings"
	"testing"
)

func TestLoadTemplateData(t *testing.T) {
	data := map[string]any{}
	templateData := map[string]any{}
	targetData := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "database", "models", "fixtures", "member.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join(".", "templates", "member.template.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &templateData)
	if err != nil {
		t.Fatal(err)
	}

	delete(templateData, "1. OFFICIAL DETAILS")

	content, err = os.ReadFile(filepath.Join("..", "database", "models", "fixtures", "member.template.output.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &targetData)
	if err != nil {
		t.Fatal(err)
	}

	result := menus.LoadTemplateData(data, templateData)

	if !reflect.DeepEqual(targetData, result) {
		t.Fatal("Test failed")
	}
}

func TestTabulateData(t *testing.T) {
	data := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "database", "models", "fixtures", "member.template.output.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	content, err = os.ReadFile(filepath.Join("..", "database", "models", "fixtures", "member.txt"))
	if err != nil {
		t.Fatal(err)
	}

	target := string(content)

	result := menus.TabulateData(data)

	if os.Getenv("DEBUG") == "true" {
		fmt.Println(strings.Join(result, "\n"))

		os.WriteFile(filepath.Join("..", "database", "models", "fixtures", "member.txt"), []byte(strings.Join(result, "\n")), 0644)
	}

	if utils.CleanString(target) != utils.CleanString(strings.Join(result, "\n")) {
		t.Fatal("Test failed")
	}
}
