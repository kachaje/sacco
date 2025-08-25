package menutests

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sacco/server"
	"sacco/server/menus"
	"sacco/server/parser"
	"sacco/utils"
	"strings"
	"testing"
)

var workflowsData map[string]map[string]any

func init() {
	var err error

	workflowsData = map[string]map[string]any{}

	err = fs.WalkDir(server.RawWorkflows, ".", func(file string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(file, ".yml") {
			return nil
		}

		content, err := server.RawWorkflows.ReadFile(file)
		if err != nil {
			return err
		}

		data, err := utils.LoadYaml(string(content))
		if err != nil {
			log.Fatal(err)
		}

		model := strings.Split(filepath.Base(file), ".")[0]

		workflowsData[model] = data

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func TestMembershipApplication(t *testing.T) {
	phoneNumber := "0999888777"

	session := parser.NewSession(nil, nil, nil)

	menus.Sessions[phoneNumber] = session

	for model, data := range workflowsData {
		session.WorkflowsMapping[model] = parser.NewWorkflow(data, nil, nil, &phoneNumber, nil, nil, nil, nil, menus.Sessions, nil)
	}

	text := ""

	m := menus.NewMenus()

	result := m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target := `
CON Welcome to Kaso SACCO
1. Membership Application
2. Loans
3. Check Balance
4. Banking Details
5. Preferred Language
6. Exit
`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "1"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
CON Choose Activity
1. Member Details
2. Occupation Details
3. Contact Details
4. Next of Kin Details
5. Beneficiaries
6. View Member Details

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "1"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	fmt.Println(result)
}
