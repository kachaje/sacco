package server_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sacco/server"
	"sacco/server/menus"
	"testing"
)

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
