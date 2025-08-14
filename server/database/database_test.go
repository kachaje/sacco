package database_test

import (
	"os"
	"sacco/server/database"
	"slices"
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
			os.Remove("test.db")
		}
	}()

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
