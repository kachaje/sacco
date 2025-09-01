package database_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sacco/server/database"
	"sacco/utils"
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

	os.WriteFile(filepath.Join(".", "models", "fixtures", "sample.json"), payload, 0644)

	target, err := os.ReadFile(filepath.Join(".", "models", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	if utils.CleanScript(payload) != utils.CleanScript(target) {
		t.Fatal("Test failed")
	}
}

func TestFullMemberRecord(t *testing.T) {
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

	_ = payload

	target, err := os.ReadFile(filepath.Join(".", "models", "fixtures", "sample.json"))
	if err != nil {
		t.Fatal(err)
	}

	if utils.CleanScript(payload) != utils.CleanScript(target) {
		t.Fatal("Test failed")
	}
}
