package menus_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sacco/server/menus"
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
