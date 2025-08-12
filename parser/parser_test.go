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

	fmt.Printf("%#v\n", result)
}
