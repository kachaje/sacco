package menutests

import (
	"io/fs"
	"log"
	"path/filepath"
	"reflect"
	"sacco/server/menus"
	menufuncs "sacco/server/menus/menuFuncs"
	"sacco/server/parser"
	"sacco/utils"
	"strings"
	"testing"
)

var workflowsData map[string]map[string]any

func init() {
	var err error

	workflowsData = map[string]map[string]any{}

	err = fs.WalkDir(menus.RawWorkflows, ".", func(file string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(file, ".yml") {
			return nil
		}

		content, err := menus.RawWorkflows.ReadFile(file)
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

func TestNewMembership(t *testing.T) {
	var data map[string]any

	refData := map[string]any{
		"dateOfBirth":       "1999-09-01",
		"firstName":         "Mary",
		"gender":            "Female",
		"lastName":          "Banda",
		"maritalStatus":     "Single",
		"nationalId":        "DHDG57636",
		"phoneNumber":       "0999888777",
		"title":             "Miss",
		"utilityBillNumber": "98383727",
		"utilityBillType":   "ESCOM",
	}

	phoneNumber := "0999888777"

	session := parser.NewSession(nil, nil, nil)

	role := "admin"

	session.SessionUserRole = &role

	session.CurrentPhoneNumber = "0999888777"

	session.GlobalIds = map[string]any{
		"memberId":     1,
		"memberLoanId": 1,
	}

	menufuncs.Sessions[phoneNumber] = session

	callbackFn := func(a any, s1, s2, s3 *string, f func(map[string]any, string, int) (*int64, error), m1 map[string]*parser.Session, m2 map[string]any) error {
		if val, ok := a.(map[string]any); ok {
			data = val
		}

		session.UpdateActiveData(data, 0)

		session.UpdateSessionFlags(nil)

		return nil
	}

	for model, data := range workflowsData {
		session.WorkflowsMapping[model] = parser.NewWorkflow(data, callbackFn, nil, &phoneNumber, nil, nil, nil, menufuncs.Sessions, nil)
	}

	text := ""

	demo := true

	m := menus.NewMenus(nil, &demo)

	result := m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target := `
CON Welcome to Kaso SACCO
1. Membership Application
2. Loans
3. Check Balance
4. Banking Details
5. Preferred Language
6. Administration
7. Exit
9. Set PhoneNumber
`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "1"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
CON Choose Activity
1. Member Details
2. Contact Details
3. Next of Kin Details
4. Beneficiaries
5. View Member Details

00. Main Menu
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "1"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
First Name: 

99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "Mary"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
Last Name: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "Banda"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
Other Name: 

00. Main Menu
02. Skip
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "02"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
Gender: 
1. Female
2. Male

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "1"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
Phone Number: (0999888777)

00. Main Menu
01. Keep
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "01"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
Title: 
1. Mr
2. Mrs
3. Miss
4. Dr
5. Prof
6. Rev
7. Other

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "3"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
Marital Status: 
1. Married
2. Single
3. Widowed
4. Divorced

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "2"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
Date Of Birth: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "1999-09-01"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
National Id: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "DHDG57636"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
Utility Bill Type: 
1. ESCOM
2. Water Board

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "1"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
Utility Bill Number: 

00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "98383727"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
Summary
- First Name: Mary
- Last Name: Banda
- Gender: Female
- Phone Number: 0999888777
- Title: Miss
- Marital Status: Single
- Date Of Birth: 1999-09-01
- National Id: DHDG57636
- Utility Bill Type: ESCOM
- Utility Bill Number: 98383727

0. Submit
00. Main Menu
98. Back
99. Cancel
	`

	if utils.CleanString(result) != utils.CleanString(target) {
		t.Fatal("Test failed")
	}

	text = "0"

	result = m.LoadMenu(session.CurrentMenu, session, phoneNumber, text, "")

	target = `
CON Choose Activity
1. Member Details (*)
2. Contact Details 
3. Next of Kin Details 
4. Beneficiaries 
5. View Member Details

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
