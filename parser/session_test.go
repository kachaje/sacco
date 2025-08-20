package parser_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sacco/parser"
	"testing"
)

func TestUpdateSessionFlags(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("..", "server", "database", "models", "fixtures", "member.json"))
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]any{}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	session := parser.NewSession(nil)
	session.ActiveMemberData = data

	err = session.UpdateSessionFlags()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadMemberCache(t *testing.T) {
	session := parser.NewSession(nil)

	err := session.LoadMemberCache("0999888777", filepath.Join("..", "server", "database", "models", "fixtures", "cache"))
	if err != nil {
		t.Fatal(err)
	}

	if !session.ContactsAdded {
		t.Fatalf("Test failed. Expected: true; Actual: %v", session.ContactsAdded)
	}

	if !session.NomineeAdded {
		t.Fatalf("Test failed. Expected: true; Actual: %v", session.NomineeAdded)
	}

	if !session.BeneficiariesAdded {
		t.Fatalf("Test failed. Expected: true; Actual: %v", session.BeneficiariesAdded)
	}

	if !session.OccupationAdded {
		t.Fatalf("Test failed. Expected: true; Actual: %v", session.OccupationAdded)
	}
}
