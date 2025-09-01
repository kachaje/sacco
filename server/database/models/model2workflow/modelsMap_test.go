package model2workflow_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sacco/server/parser"
	"testing"
)

func TestModelsMap(t *testing.T) {
	sample := map[string]any{}

	content, err := os.ReadFile(filepath.Join("..", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &sample)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]string{}

	content, err = os.ReadFile(filepath.Join(".", "modelsMap.json"))
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	results := map[string]any{}

	session := parser.NewSession(nil, nil, nil)

	for key, value := range data {
		result, ok := session.DecodeKey(value, sample)
		if ok {
			results[key] = result
		} else {
			fmt.Println("missed:", key)
		}
	}

	payload, _ := json.MarshalIndent(results, "", "  ")

	fmt.Println(string(payload))
}
