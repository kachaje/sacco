package database_test

import (
	"encoding/json"
	"fmt"
	"sacco/server/database"
	"testing"

	_ "embed"
)

//go:embed models/fixtures/sample.sql
var sampleScript string

func TestFullRecord(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	_, err := db.DB.Exec(sampleScript)
	if err != nil {
		t.Fatal(err)
	}

	phoneNumber := "09999999999"

	result, err := db.FullRecord(phoneNumber)
	if err != nil {
		t.Fatal(err)
	}

	payload, _ := json.MarshalIndent(result, "", "  ")

	fmt.Println(string(payload))
}
