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

func setupDb() (*database.Database, error) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)

	_, err := db.DB.Exec(sampleScript)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestLoadModelChildren(t *testing.T) {
	db, err := setupDb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	result, err := db.LoadModelChildren("member", 1)
	if err != nil {
		t.Fatal(err)
	}

	payload, _ := json.MarshalIndent(result, "", "  ")

	fmt.Println(string(payload))
}

func TestFullRecord(t *testing.T) {
	db, err := setupDb()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	phoneNumber := "09999999999"

	result, err := db.FullMemberRecord(phoneNumber)
	if err != nil {
		t.Fatal(err)
	}

	payload, _ := json.MarshalIndent(result, "", "  ")

	fmt.Println(string(payload))
}
