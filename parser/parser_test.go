package parser_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sacco/parser"
	"testing"
)

var content []byte
var err error
var data map[string]any

func init() {
	content, err = os.ReadFile(filepath.Join(".", "fixtures", "sample.json"))
	if err != nil {
		panic(err)
	}

	data = map[string]any{}

	err = json.Unmarshal(content, &data)
	if err != nil {
		panic(err)
	}
}

func TestGetNode(t *testing.T) {
	wf := parser.NewWorkflow(data)

	result := wf.GetNode("enterLanguage")

	if result == nil {
		t.Fatal("Test failed")
	}

	for _, key := range []string{"type", "text", "options", "inputIdentifier", "nextScreen"} {
		if result[key] == nil {
			t.Fatalf("Test failed on key %s", key)
		}
	}
}

func TestInputIncluded(t *testing.T) {
	options := []map[string]any{
		{
			"position": 1,
			"label": map[string]any{
				"all": "English",
			},
		},
		{
			"position": 2,
			"label": map[string]any{
				"all": "Chichewa",
			},
		},
	}

	wf := parser.NewWorkflow(data)

	result := wf.InputIncluded("3", options)

	if result {
		t.Fatalf("Test failed. Expected: false; Actual: %v", result)
	}

	result = wf.InputIncluded("1", options)

	if !result {
		t.Fatalf("Test failed. Expected: true; Actual: %v", result)
	}
}

func TestNodeOptions(t *testing.T) {
	wf := parser.NewWorkflow(data)

	result := wf.NodeOptions("enterLanguage")

	if len(result) != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", len(result))
	}

	for i, entry := range []string{"1. English", "2. Chichewa"} {
		if result[i] != entry {
			t.Fatalf("Test failed. Expected: %s; Actual: %s", entry, result[i])
		}
	}

	wf.CurrentLanguage = "ny"

	result = wf.NodeOptions("enterMaritalStatus")

	if len(result) != 4 {
		t.Fatalf("Test failed. Expected: 4; Actual: %v", len(result))
	}

	for i, entry := range []string{"1. Inde", "2. Ayi", "3. Woferedwa", "4. Osudzulidwa"} {
		if result[i] != entry {
			t.Fatalf("Test failed. Expected: %s; Actual: %s", entry, result[i])
		}
	}

	wf.CurrentLanguage = "en"

	result = wf.NodeOptions("enterMaritalStatus")

	if len(result) != 4 {
		t.Fatalf("Test failed. Expected: 4; Actual: %v", len(result))
	}

	for i, entry := range []string{"1. Married", "2. Single", "3. Widowed", "4. Divorced"} {
		if result[i] != entry {
			t.Fatalf("Test failed. Expected: %s; Actual: %s", entry, result[i])
		}
	}
}

func TestNextNode(t *testing.T) {
	wf := parser.NewWorkflow(data)

	result := wf.NextNode("")

	if result == nil {
		t.Fatal("Test failed")
	}

	for _, key := range []string{"type", "text", "options", "inputIdentifier", "nextScreen"} {
		if result[key] == nil {
			t.Fatalf("Test failed on key %s", key)
		}
	}

	if wf.CurrentScreen != "enterLanguage" {
		t.Fatalf("Test failed. Expected: 'enterLanguage'; Actual: '%v'", wf.CurrentScreen)
	}

	if wf.PreviousScreen != "initialScreen" {
		t.Fatalf("Test failed. Expected: 'initialScreen'; Actual: '%v'", wf.PreviousScreen)
	}

	result = wf.NextNode("3")

	for _, key := range []string{"type", "text", "options", "inputIdentifier", "nextScreen"} {
		if result[key] == nil {
			t.Fatalf("Test failed on key %s", key)
		}
	}

	if wf.CurrentScreen != "enterLanguage" {
		t.Fatalf("Test failed. Expected: 'enterLanguage'; Actual: '%v'", wf.CurrentScreen)
	}

	if wf.PreviousScreen != "initialScreen" {
		t.Fatalf("Test failed. Expected: 'initialScreen'; Actual: '%v'", wf.PreviousScreen)
	}

	fmt.Printf("%#v\n", result)
}
