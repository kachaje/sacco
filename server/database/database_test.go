package database_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sacco/server"
	"sacco/server/database"
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

	payload, _ := json.MarshalIndent(result, "", "  ")

	if utils.CleanScript(payload) != utils.CleanScript(target) {
		t.Fatal("Test failed")
	}
}

func TestMemberByDefaultPhoneNumber(t *testing.T) {
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

	result, err := db.MemberByDefaultPhoneNumber("09999999999")
	if err != nil {
		t.Fatal(err)
	}

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

	result, err := db.MemberByDefaultPhoneNumber("09999999999")
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

	err = server.SaveData(update, &model, nil, nil, nil, nil, db.AddMember, nil, beneficiaries)
	if err != nil {
		t.Fatal(err)
	}

	result, err = db.MemberByDefaultPhoneNumber("09999999999")
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
