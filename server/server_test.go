package server_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sacco/server"
	"sacco/server/menus"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"server": server.Main,
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
