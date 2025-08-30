package menutests

import (
	"fmt"
	"reflect"
	"sacco/server/menus"
	menufuncs "sacco/server/menus/menuFuncs"
	"sacco/server/parser"
	"sacco/utils"
	"testing"
)

func TestMemberBeneficiaries(t *testing.T) {
	var data map[string]any

	refData := map[string]any{
		"contact1":    "P.O. Box 1",
		"contact2":    "P.O. Box 2",
		"contact3":    "P.O. Box 3",
		"contact4":    "P.O. Box 4",
		"name1":       "John Phiri",
		"name2":       "Peter Banda",
		"name3":       "Mirriam Jere",
		"name4":       "Bornface Harawa",
		"percentage1": "10",
		"percentage2": "8",
		"percentage3": "6",
		"percentage4": "4",
	}

	phoneNumber := "0999888776"

	session := parser.NewSession(nil, nil, nil)

	menufuncs.Sessions[phoneNumber] = session

	callbackFn := func(a any, s1, s2, s3 *string, f func(map[string]any, string, int) (*int64, error), m1 map[string]*parser.Session, m2 map[string]any) error {
		if val, ok := a.(map[string]any); ok {
			data = val
		}

		localData := []map[string]any{}

		for i := range 4 {
			idLabel := fmt.Sprintf("id%d", i+1)
			nameLabel := fmt.Sprintf("name%d", i+1)
			contactLabel := fmt.Sprintf("contact%d", i+1)
			percentageLabel := fmt.Sprintf("percentage%d", i+1)

			if data[nameLabel] == nil {
				break
			}

			row := map[string]any{}

			if data[idLabel] != nil {
				row["id"] = data[idLabel]
			}
			if data[nameLabel] != nil {
				row["name"] = data[nameLabel]
			}
			if data[contactLabel] != nil {
				row["contact"] = data[contactLabel]
			}
			if data[percentageLabel] != nil {
				row["percentage"] = data[percentageLabel]
			}

			localData = append(localData, row)
		}

		session.UpdateActiveData(map[string]any{
			"memberBeneficiary": localData,
		}, 0)

		session.UpdateSessionFlags()

		return nil
	}

	for model, data := range workflowsData {
		session.WorkflowsMapping[model] = parser.NewWorkflow(data, callbackFn, nil, &phoneNumber, nil, nil, nil, menufuncs.Sessions, nil)
	}

	text := ""

	demo := true

	m := menus.NewMenus(nil, &demo)

	result := m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target := `
CON Welcome to Kaso SACCO
1. Membership Application
2. Loans
3. Check Balance
4. Banking Details
5. Preferred Language
6. Administration
7. Exit
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

	text = "5"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Name: 

02. Skip
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "John Phiri"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Percentage: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "10"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Contact: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "P.O. Box 1"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Name: 

00. Main Menu
02. Skip
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "Peter Banda"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Percentage: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "8"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Contact: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "P.O. Box 2"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Name: 

00. Main Menu
02. Skip
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "Mirriam Jere"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Percentage: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "6"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Contact: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "P.O. Box 3"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Name: 

00. Main Menu
02. Skip
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "Bornface Harawa"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Percentage: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "4"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Contact: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "P.O. Box 4"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Summary
- Name: John Phiri
- Percentage: 10
- Contact: P.O. Box 1
- Name: Peter Banda
- Percentage: 8
- Contact: P.O. Box 2
- Name: Mirriam Jere
- Percentage: 6
- Contact: P.O. Box 3
- Name: Bornface Harawa
- Percentage: 4
- Contact: P.O. Box 4

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
3. Contact Details 
4. Next of Kin Details 
5. Beneficiaries (*)
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

	text = "5"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Name: (John Phiri)

01. Keep
02. Skip
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "01"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Percentage: (10)

00. Main Menu
01. Keep
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "01"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Contact: (P.O. Box 1)

00. Main Menu
01. Keep
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "01"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Name: (Peter Banda)

00. Main Menu
01. Keep
02. Skip
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "01"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Percentage: (8)

00. Main Menu
01. Keep
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "20"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Contact: (P.O. Box 2)

00. Main Menu
01. Keep
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "P.O. Box 348589, Lilongwe"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Name: (Mirriam Jere)

00. Main Menu
01. Keep
02. Skip
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "01"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Percentage: (6)

00. Main Menu
01. Keep
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "01"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Contact: (P.O. Box 3)

00. Main Menu
01. Keep
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "01"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Name: (Bornface Harawa)

00. Main Menu
01. Keep
02. Skip
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "01"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Percentage: (4)

00. Main Menu
01. Keep
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "01"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Contact: (P.O. Box 4)

00. Main Menu
01. Keep
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "01"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "", "")

	target = `
Summary
- Name: John Phiri
- Percentage: 10
- Contact: P.O. Box 1
- Name: Peter Banda
- Percentage: 20
- Contact: P.O. Box 348589, Lilongwe
- Name: Mirriam Jere
- Percentage: 6
- Contact: P.O. Box 3
- Name: Bornface Harawa
- Percentage: 4
- Contact: P.O. Box 4

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
3. Contact Details 
4. Next of Kin Details 
5. Beneficiaries (*)
6. View Member Details

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	refData = map[string]any{
		"contact1":    "P.O. Box 1",
		"contact2":    "P.O. Box 348589, Lilongwe",
		"contact3":    "P.O. Box 3",
		"contact4":    "P.O. Box 4",
		"name1":       "John Phiri",
		"name2":       "Peter Banda",
		"name3":       "Mirriam Jere",
		"name4":       "Bornface Harawa",
		"percentage1": "10",
		"percentage2": "20",
		"percentage3": "6",
		"percentage4": "4",
	}

	if !reflect.DeepEqual(data, refData) {
		t.Fatal("Test failed")
	}
}
