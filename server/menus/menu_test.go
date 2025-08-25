package menus_test

import (
	"sacco/server/menus"
	"sacco/utils"
	"testing"
)

func TestLoadMenus(t *testing.T) {
	m := menus.NewMenus()

	result := m.LoadMenu("main", nil, "", "", "", "")

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
}
