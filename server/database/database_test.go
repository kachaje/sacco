package database_test

import (
	"encoding/json"
	"os"
	"path/filepath"
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
			contactDetails = val
			delete(data, "contactDetails")
		} else {
			t.Fatal("Test failed. Failed to convert map")
		}
	}
	{
		val, ok := data["nominee"].(map[string]any)
		if ok {
			nominee = val
			delete(data, "nominee")
		} else {
			t.Fatal("Test failed. Failed to convert map")
		}
	}
	{
		val, ok := data["occupationDetails"].(map[string]any)
		if ok {
			occupationDetails = val
			delete(data, "occupationDetails")
		} else {
			t.Fatal("Test failed. Failed to convert map")
		}
	}

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
