package utils_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sacco/utils"
	"testing"
)

func TestLoadYaml(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(".", "fixtures", "newMember.yml"))
	if err != nil {
		t.Fatal(err)
	}

	result, err := utils.LoadYaml(string(content))
	if err != nil {
		t.Fatal(err)
	}

	target := map[string]any{}

	refData, err := os.ReadFile(filepath.Join(".", "fixtures", "newMember.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(refData, &target)
	if err != nil {
		t.Fatal(err)
	}

	compareObjects := func(obj1, obj2 map[string]any) bool {
		if len(obj1) != len(obj2) {
			return false
		}

		for key, val1 := range obj1 {
			val2, exists := obj2[key]
			if !exists || fmt.Sprintf("%v", val1) != fmt.Sprintf("%v", val2) {
				return false
			}
		}

		return true
	}

	if !compareObjects(target, result) {
		t.Fatal("Test failed")
	}
}
