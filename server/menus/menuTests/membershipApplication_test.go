package menutests

import (
	"fmt"
	"sacco/server/menus"
	"sacco/server/parser"
	"sacco/utils"
	"testing"
)

func TestMembershipApplication(t *testing.T) {
	m := menus.NewMenus()

	phoneNumber := "0999888777"
	session := parser.NewSession(nil, nil, nil)

	text := ""

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
