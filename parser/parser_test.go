package parser_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sacco/parser"
	"testing"
)

var content []byte
var err error

func init() {
	content, err = os.ReadFile(filepath.Join(".", "fixtures", "sample.json"))
	if err != nil {
		panic(err)
	}
}

func TestGetNode(t *testing.T) {
	data := map[string]any{}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

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
