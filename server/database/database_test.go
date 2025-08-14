package database_test

import (
	"database/sql"
	"os"
	"sacco/server/database"
	"sacco/utils"
	"slices"
	"strings"
	"testing"
)

func TestDatabase(t *testing.T) {
	db := database.NewDatabase("test")
	defer func() {
		db.Close()

		_, err := os.Stat("test.db")
		if os.IsNotExist(err) {
			t.Fatal("Test failed")
		} else {
			// os.Remove("test.db")
		}
	}()

	row := db.DB.QueryRow(`SELECT sql FROM sqlite_schema WHERE name="member"`)

	target := utils.CleanString(`
CREATE TABLE member (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	firstName TEXT,
	lastName TEXT,
	otherName TEXT,
	gender TEXT,
	title TEXT
	maritalStatus TEXT,
	dateOfBirth TEXT,
	nationalId TEXT,
	utilityBillType TEXT,
	utilityBillNumber TEXT
)`)

	var result string

	if err := row.Scan(&result); err == sql.ErrNoRows {
		t.Fatal("Test failed")
	}

	result = utils.CleanString(result)

	if strings.Compare(result, target) != 0 {
		t.Fatal("Test failed")
	}

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
