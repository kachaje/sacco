package parser_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sacco/server/parser"
	"testing"
)

func TestUpdateSessionFlags(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("..", "database", "models", "fixtures", "member.json"))
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]any{}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	session := parser.NewSession(nil, nil, nil)
	session.UpdateActiveData(data, 0)

	err = session.UpdateSessionFlags()
	if err != nil {
		t.Fatal(err)
	}
}
