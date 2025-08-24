package database_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sacco/server/database"
	filehandling "sacco/server/fileHandling"
	"sacco/utils"
	"slices"
	"testing"
)

func TestDatabase(t *testing.T) {
	db := database.NewDatabase(":memory:")
	defer db.Close()

	rows, err := db.DB.QueryContext(t.Context(), `SELECT name FROM sqlite_schema WHERE type="table"`)
	if err != nil {
		t.Fatal(err)
	}

	tables := []string{}

	for rows.Next() {
		var result string

		if err := rows.Scan(&result); err != nil {
			t.Fatal(err)
		}

		tables = append(tables, result)
	}

	for _, table := range []string{"member", "sqlite_sequence", "memberContact", "memberNominee", "memberOccupation", "memberBeneficiary"} {
		if !slices.Contains(tables, table) {
			t.Fatalf("Test failed. Missing: %s", table)
		}
	}
}
func TestDatabaseMemberBeneficiaries(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer func() {
		db.Close()

		if _, err := os.Stat("memberBeneficiary.json"); !os.IsNotExist(err) {
			os.Remove("memberBeneficiary.json")
		}
	}()

	content, err := os.ReadFile(filepath.Join(".", "models", "fixtures", "member.sql"))
	if err != nil {
		t.Fatal(err)
	}

	sqlStatement := string(content)

	_, err = db.DB.Exec(sqlStatement)
	if err != nil {
		t.Fatal(err)
	}

	phoneNumber := "09999999999"

	result, err := db.MemberByPhoneNumber(phoneNumber, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	memberBeneficiary := map[string]any{}

	val, ok := result["memberBeneficiary"].([]map[string]any)
	if ok {
		for i, row := range val {
			for key, value := range row {
				keyLabel := fmt.Sprintf("%s%d", key, i+1)

				memberBeneficiary[keyLabel] = value
			}
		}
	}

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(memberBeneficiary, "", "  ")

		fmt.Println(string(payload))
	}

	update := map[string]any{
		"contact1":    "P.O. Box 1234",
		"id1":         1,
		"memberId1":   1,
		"name1":       "Benefator 1",
		"percentage1": 35,
		"contact2":    "P.O. Box 5678",
		"id2":         2,
		"memberId2":   1,
		"name2":       "Benefator 2",
		"percentage2": 25,
	}

	model := "memberBeneficiary"

	err = filehandling.SaveData(update, &model, nil, nil, nil, db.GenericsSaveData, nil, memberBeneficiary)
	if err != nil {
		t.Fatal(err)
	}

	result, err = db.MemberByPhoneNumber(phoneNumber, nil, []string{})
	if err != nil {
		t.Fatal(err)
	}

	{
		memberBeneficiary := []map[string]any{}

		rows, ok := result["memberBeneficiary"].([]map[string]any)
		if ok {
			memberBeneficiary = append(memberBeneficiary, rows...)
		}

		if os.Getenv("DEBUG") == "true" {
			payload, _ := json.MarshalIndent(memberBeneficiary, "", "  ")

			fmt.Println(string(payload))
		}

		if len(memberBeneficiary) != 2 {
			t.Fatalf("Test failed. Expected: 2; Actual: %v", len(memberBeneficiary))
		}
	}
}

func TestGenericsSaveData(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	data := map[string]any{
		"memberId":        1,
		"yearsInBusiness": 1,
		"businessNature":  "Vendor",
		"businessName":    "Vendors Galore",
		"tradingArea":     "Mtandire",
	}

	mid, err := db.GenericsSaveData(data, "memberBusiness", 0)
	if err != nil {
		t.Fatal(err)
	}

	if mid == nil {
		t.Fatal("Test failed. Got nil id")
	}

	{
		result, err := db.GenericModels["memberBusiness"].FetchById(*mid)
		if err != nil {
			t.Fatal(err)
		}

		if result == nil {
			t.Fatal("Test failed. Got nil result")
		}

		for key, value := range data {
			if result[key] == nil {
				t.Fatal("Test failed")
			}

			if fmt.Sprintf("%v", result[key]) != fmt.Sprintf("%v", value) {
				t.Fatalf("Test failed. Expected: %v; Actual: %v", value, result[key])
			}
		}
	}
}

func TestGenericModel(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	data := map[string]any{
		"memberId":        1,
		"yearsInBusiness": 1,
		"businessNature":  "Vendor",
		"businessName":    "Vendors Galore",
		"tradingArea":     "Mtandire",
	}

	mid, err := db.GenericModels["memberBusiness"].AddRecord(data)
	if err != nil {
		t.Fatal(err)
	}

	if mid == nil {
		t.Fatal("Test failed. Got nil id")
	}

	{
		result, err := db.GenericModels["memberBusiness"].FetchById(*mid)
		if err != nil {
			t.Fatal(err)
		}

		if result == nil {
			t.Fatal("Test failed. Got nil result")
		}

		for key, value := range data {
			if result[key] == nil {
				t.Fatal("Test failed")
			}

			if fmt.Sprintf("%v", result[key]) != fmt.Sprintf("%v", value) {
				t.Fatalf("Test failed. Expected: %v; Actual: %v", value, result[key])
			}
		}
	}

	{
		result, err := db.GenericModels["memberBusiness"].FilterBy(`WHERE businessNature="Vendor"`)
		if err != nil {
			t.Fatal(err)
		}

		if len(result) <= 0 {
			t.Fatal("Test failed. Got nil result")
		}

		for key, value := range data {
			if result[0][key] == nil {
				t.Fatal("Test failed")
			}

			if fmt.Sprintf("%v", result[0][key]) != fmt.Sprintf("%v", value) {
				t.Fatalf("Test failed. Expected: %v; Actual: %v", value, result[0][key])
			}
		}
	}

	{
		err = db.GenericModels["memberBusiness"].UpdateRecord(map[string]any{
			"businessNature": "Taxi",
		}, *mid)
		if err != nil {
			t.Fatal(err)
		}

		result, err := db.GenericModels["memberBusiness"].FetchById(*mid)
		if err != nil {
			t.Fatal(err)
		}

		if result == nil {
			t.Fatal("Test failed. Got nil result")
		}

		for key, value := range data {
			if result[key] == nil {
				t.Fatal("Test failed")
			}

			if key == "businessNature" {
				if result[key].(string) != "Taxi" {
					t.Fatalf("Test failed. Expected: Taxi; Actual: %v", result[key])
				}
			} else {
				if fmt.Sprintf("%v", result[key]) != fmt.Sprintf("%v", value) {
					t.Fatalf("Test failed. Expected: %v; Actual: %v", value, result[key])
				}
			}
		}
	}
}

func TestMemberByPhoneNumber(t *testing.T) {
	phoneNumber := "0999888777"

	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer func() {
		db.Close()
	}()

	target := map[string]any{
		"dateOfBirth":   "1999-09-01",
		"firstName":     "Mary",
		"gender":        "Female",
		"id":            1,
		"lastName":      "Banda",
		"maritalStatus": "Single",
		"memberBeneficiary": []map[string]any{
			{
				"contact":    "0888777444",
				"id":         1,
				"memberId":   1,
				"name":       "John Phiri",
				"percentage": 10,
			},
			{
				"contact":    "07746635653",
				"id":         2,
				"memberId":   1,
				"name":       "Jean Banda",
				"percentage": 5,
			},
		},
		"memberContact": map[string]any{
			"homeDistrict":       "Lilongwe",
			"homeTA":             "Kalolo",
			"homeVillage":        "Kalulu",
			"id":                 1,
			"memberId":           1,
			"postalAddress":      "P.O. Box 1",
			"residentialAddress": "Area 49",
		},
		"memberNominee": map[string]any{
			"id":          1,
			"memberId":    1,
			"address":     "P.O. Box 1",
			"name":        "John Phiri",
			"phoneNumber": "0888444666",
		},
		"memberOccupation": map[string]any{
			"employerAddress":      "Kanengo",
			"employerName":         "SOBO",
			"employerPhone":        "01282373737",
			"grossPay":             100000,
			"highestQualification": "Secondary",
			"id":                   1,
			"jobTitle":             "Driver",
			"memberId":             1,
			"netPay":               90000,
			"periodEmployed":       36,
		},
		"nationalId":        "DHFYR8475",
		"phoneNumber":       "0999888777",
		"title":             "Miss",
		"utilityBillNumber": "29383746",
		"utilityBillType":   "ESCOM",
	}

	member := map[string]any{
		"dateOfBirth":       "1999-09-01",
		"phoneNumber":       phoneNumber,
		"fileNumber":        "",
		"firstName":         "Mary",
		"gender":            "Female",
		"id":                1,
		"lastName":          "Banda",
		"maritalStatus":     "Single",
		"nationalId":        "DHFYR8475",
		"oldFileNumber":     "",
		"otherName":         "",
		"title":             "Miss",
		"utilityBillNumber": "29383746",
		"utilityBillType":   "ESCOM",
	}

	memberContact := map[string]any{
		"homeDistrict":       "Lilongwe",
		"homeTA":             "Kalolo",
		"homeVillage":        "Kalulu",
		"phoneNumber":        "0999888777",
		"postalAddress":      "P.O. Box 1",
		"residentialAddress": "Area 49",
	}

	memberNominee := map[string]any{
		"address":     "P.O. Box 1",
		"name":        "John Phiri",
		"phoneNumber": "0888444666",
	}

	memberOccupation := map[string]any{
		"employerAddress":      "Kanengo",
		"employerName":         "SOBO",
		"employerPhone":        "01282373737",
		"grossPay":             100000,
		"highestQualification": "Secondary",
		"jobTitle":             "Driver",
		"netPay":               90000,
		"periodEmployed":       36,
	}

	memberBeneficiary := []map[string]any{
		{
			"contact":    "0888777444",
			"name":       "John Phiri",
			"percentage": 10,
		},
		{
			"contact":    "07746635653",
			"name":       "Jean Banda",
			"percentage": 5,
		},
	}

	id, err := db.GenericModels["member"].AddRecord(member)
	if err != nil {
		t.Fatal(err)
	}

	memberContact["memberId"] = *id
	memberNominee["memberId"] = *id
	memberOccupation["memberId"] = *id
	memberBeneficiary[0]["memberId"] = *id
	memberBeneficiary[1]["memberId"] = *id

	_, err = db.GenericModels["memberContact"].AddRecord(memberContact)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GenericModels["memberNominee"].AddRecord(memberNominee)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.GenericModels["memberOccupation"].AddRecord(memberOccupation)
	if err != nil {
		t.Fatal(err)
	}

	for i := range memberBeneficiary {
		_, err = db.GenericModels["memberBeneficiary"].AddRecord(memberBeneficiary[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	result, err := db.MemberByPhoneNumber(phoneNumber, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	payloadResult, _ := json.MarshalIndent(result, "", "  ")
	payloadTarget, _ := json.MarshalIndent(target, "", "  ")

	if utils.CleanScript(payloadResult) != utils.CleanScript(payloadTarget) {
		t.Fatal("Test failed")
	}
}
