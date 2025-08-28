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

func TestLoadMemberCache(t *testing.T) {
	session := parser.NewSession(nil, nil, nil)

	err := session.LoadCacheData("0999888777", filepath.Join("..", "database", "models", "fixtures", "cache"))
	if err != nil {
		t.Fatal(err)
	}

	if !session.AddedModels["memberContact"] {
		t.Fatalf("Test failed. Expected: true; Actual: %v", session.AddedModels["memberContact"])
	}

	if !session.AddedModels["memberNominee"] {
		t.Fatalf("Test failed. Expected: true; Actual: %v", session.AddedModels["memberNominee"])
	}

	if !session.AddedModels["memberBeneficiary"] {
		t.Fatalf("Test failed. Expected: true; Actual: %v", session.AddedModels["memberBeneficiary"])
	}

	if !session.AddedModels["memberOccupation"] {
		t.Fatalf("Test failed. Expected: true; Actual: %v", session.AddedModels["memberOccupation"])
	}
}
