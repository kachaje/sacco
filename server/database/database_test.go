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

func TestDatabaseAddMember(t *testing.T) {
	db := database.NewDatabase(":memory:")
	defer db.Close()

	target, err := os.ReadFile(filepath.Join(".", "models", "fixtures", "member.json"))
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]any{}

	err = json.Unmarshal(target, &data)
	if err != nil {
		t.Fatal(err)
	}

	var beneficiaries []map[string]any
	var contactDetails map[string]any
	var nominee map[string]any
	var occupationDetails map[string]any

	{
		val, ok := data["beneficiaries"].([]any)
		if ok {
			for i := range val {
				vl, ok := val[i].(map[string]any)
				if ok {
					delete(vl, "id")
					beneficiaries = append(beneficiaries, vl)
				}
			}
			delete(data, "beneficiaries")
		} else {
			t.Fatal("Test failed. Failed to convert map")
		}
	}
	{
		val, ok := data["contactDetails"].(map[string]any)
		if ok {
			delete(val, "id")
			contactDetails = val
			delete(data, "contactDetails")
		} else {
			t.Fatal("Test failed. Failed to convert map")
		}
	}
	{
		val, ok := data["nomineeDetails"].(map[string]any)
		if ok {
			delete(val, "id")
			nominee = val
			delete(data, "nomineeDetails")
		} else {
			t.Fatal("Test failed. Failed to convert map")
		}
	}
	{
		val, ok := data["occupationDetails"].(map[string]any)
		if ok {
			delete(val, "id")
			occupationDetails = val
			delete(data, "occupationDetails")
		} else {
			t.Fatal("Test failed. Failed to convert map")
		}
	}
	delete(data, "id")

	id, err := db.AddMember(data, contactDetails, nominee, occupationDetails, beneficiaries, nil)
	if err != nil {
		t.Fatal(err)
	}

	result, err := db.Member.MemberDetails(*id)
	if err != nil {
		t.Fatal(err)
	}

	delete(result, "memberIdNumber")
	delete(result, "shortMemberId")
	delete(result, "dateJoined")

	payload, _ := json.MarshalIndent(result, "", "  ")

	if utils.CleanScript(payload) != utils.CleanScript(target) {
		t.Fatal("Test failed")
	}
}

func TestMemberByPhoneNumber(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	content, err := os.ReadFile(filepath.Join(".", "models", "fixtures", "member.sql"))
	if err != nil {
		t.Fatal(err)
	}

	sqlStatement := string(content)

	_, err = db.DB.Exec(sqlStatement)
	if err != nil {
		t.Fatal(err)
	}

	result, err := db.MemberByPhoneNumber("09999999999")
	if err != nil {
		t.Fatal(err)
	}

	delete(result, "memberIdNumber")
	delete(result, "shortMemberId")
	delete(result, "dateJoined")

	payload, _ := json.MarshalIndent(result, "", "  ")

	target, err := os.ReadFile(filepath.Join(".", "models", "fixtures", "member.json"))
	if err != nil {
		t.Fatal(err)
	}

	if utils.CleanScript(payload) != utils.CleanScript(target) {
		t.Fatal("Test failed")
	}
}

func TestMemberBeneficiaries(t *testing.T) {
	t.Skip()

	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	content, err := os.ReadFile(filepath.Join(".", "models", "fixtures", "member.sql"))
	if err != nil {
		t.Fatal(err)
	}

	sqlStatement := string(content)

	_, err = db.DB.Exec(sqlStatement)
	if err != nil {
		t.Fatal(err)
	}

	result, err := db.MemberByPhoneNumber("09999999999")
	if err != nil {
		t.Fatal(err)
	}

	beneficiaries := map[string]any{}

	val, ok := result["beneficiaries"].([]any)
	if ok {
		for i, row := range val {
			v, ok := row.(map[string]any)
			if ok {
				for key, value := range v {
					keyLabel := fmt.Sprintf("%s%d", key, i+1)

					beneficiaries[keyLabel] = value
				}
			}
		}
	}

	if os.Getenv("DEBUG") == "true" {
		payload, _ := json.MarshalIndent(beneficiaries, "", "  ")

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

	model := "beneficiaries"

	err = filehandling.SaveData(update, &model, nil, nil, nil, nil, db.GenericsSaveData, nil, beneficiaries)
	if err != nil {
		t.Fatal(err)
	}

	result, err = db.MemberByPhoneNumber("09999999999")
	if err != nil {
		t.Fatal(err)
	}

	{
		beneficiaries := []map[string]any{}

		val, ok = result["beneficiaries"].([]any)
		if ok {
			for _, row := range val {
				v, ok := row.(map[string]any)
				if ok {
					beneficiaries = append(beneficiaries, v)
				}
			}
		}

		if os.Getenv("DEBUG") == "true" {
			payload, _ := json.MarshalIndent(beneficiaries, "", "  ")

			fmt.Println(string(payload))
		}

		if len(beneficiaries) != 2 {
			t.Fatalf("Test failed. Expected: 2; Actual: %v", len(beneficiaries))
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

func TestMemberDetailsByPhoneNumber(t *testing.T) {
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
			"id":             1,
			"memberId":       1,
			"nomineeAddress": "P.O. Box 1",
			"nomineeName":    "John Phiri",
			"nomineePhone":   "0888444666",
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
		"nomineeAddress": "P.O. Box 1",
		"nomineeName":    "John Phiri",
		"nomineePhone":   "0888444666",
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

	result, err := db.MemberDetailsByPhoneNumber(phoneNumber, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	payloadResult, _ := json.MarshalIndent(result, "", "  ")
	payloadTarget, _ := json.MarshalIndent(target, "", "  ")

	if utils.CleanScript(payloadResult) != utils.CleanScript(payloadTarget) {
		t.Fatal("Test failed")
	}
}
