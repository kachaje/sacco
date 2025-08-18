package server_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sacco/server"
	"sacco/server/menus"
	"sacco/wscli"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"server": server.Main,
		"wscli":  wscli.Main,
	})
}

func TestMemberApplication(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "testdata/memberApplication",
	})
}

func TestUpdateSessionFlags(t *testing.T) {
	content, err := os.ReadFile(filepath.Join(".", "database", "models", "fixtures", "member.json"))
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]any{}

	err = json.Unmarshal(content, &data)
	if err != nil {
		t.Fatal(err)
	}

	session := &menus.Session{}
	session.ActiveMemberData = data

	err = server.UpdateSessionFlags(session)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadMemberCache(t *testing.T) {
	session := &menus.Session{}

	err := server.LoadMemberCache(session, "0999888777", filepath.Join(".", "database", "models", "fixtures", "cache"))
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
