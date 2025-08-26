package menutests

import (
	"fmt"
	"reflect"
	"sacco/server/menus"
	"sacco/server/parser"
	"sacco/utils"
	"testing"
)

func TestMemberContact(t *testing.T) {
	var data map[string]any

	refData := map[string]any{
		"employerAddress":        "Kanengo",
		"employerName":           "SOBO",
		"employerPhone":          "01789987",
		"grossPay":               "100000",
		"highestQualification":   "Secondary",
		"jobTitle":               "Driver",
		"netPay":                 "90000",
		"periodEmployedInMonths": "36",
	}

	phoneNumber := "0999888776"

	session := parser.NewSession(nil, nil, nil)

	menus.Sessions[phoneNumber] = session

	callbackFn := func(a any, s1, s2, s3, s4 *string, f func(map[string]any, string, int) (*int64, error), m1 map[string]*parser.Session, m2 map[string]any) error {
		if val, ok := a.(map[string]any); ok {
			data = val
		}

		fmt.Printf("%#v\n", data)

		session.UpdateActiveMemberData(map[string]any{
			"memberContact": data,
		}, 0)

		session.UpdateSessionFlags()

		session.AddedModels["memberContact"] = true

		return nil
	}

	for model, data := range workflowsData {
		session.WorkflowsMapping[model] = parser.NewWorkflow(data, callbackFn, nil, &phoneNumber, nil, nil, nil, nil, menus.Sessions, nil)
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

	text = "2"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Employer Name: 

99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "SOBO"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Gross Pay: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "100000"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Net Pay: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "90000"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Job Title: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "Driver"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Employer Address: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "Kanengo"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Employer Phone: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "01789987"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Period Employed In Months: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "36"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Highest Qualification: 
1. Tertiary
2. Secondary
3. Primary
4. None

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "2"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Summary
- Employer Name: SOBO
- Gross Pay: 100000
- Net Pay: 90000
- Job Title: Driver
- Employer Address: Kanengo
- Employer Phone: 01789987
- Period Employed In Months: 36
- Highest Qualification: Secondary

0. Submit
00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "0"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
CON Choose Activity
1. Member Details 
2. Occupation Details 
3. Contact Details (*)
4. Next of Kin Details 
5. Beneficiaries 
6. View Member Details

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	if data == nil {
		t.Fatal("Test failed")
	}

	if !reflect.DeepEqual(data, refData) {
		t.Fatal("Test failed")
	}
}
