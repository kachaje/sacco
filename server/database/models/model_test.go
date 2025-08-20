package models_test

import (
	"database/sql"
	"fmt"
	"os"
	"sacco/server/database/models"
	"testing"
)

var (
	tableName = "person"
)

func setupDb(dbname string) (*sql.DB, *models.Model, error) {
	db, err := sql.Open("sqlite", dbname)
	if err != nil {
		return nil, nil, err
	}

	sqlStmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		firstName TEXT,
		lastName TEXT,
		gender TEXT,
		height REAL,
		weight REAL
	);`, tableName)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, nil, err
	}

	fields := []string{"firstName", "lastName", "gender", "height", "weight"}
	fieldTypes := map[string]string{
		"id":        "int",
		"firstName": "string",
		"lastName":  "string",
		"gender":    "string",
		"height":    "float",
		"weight":    "weight",
	}

	model, err := models.NewModel(db, tableName, fields, fieldTypes)
	if err != nil {
		return nil, nil, err
	}

	return db, model, nil
}

func TestNewModel(t *testing.T) {
	dbname := "test.db"

	db, model, err := setupDb(dbname)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(dbname)
	}()

	_, _ = db, model
}

func TestAddRecord(t *testing.T) {
	dbname := "testAdd.db"

	db, model, err := setupDb(dbname)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(dbname)
	}()

	data := map[string]any{
		"firstName": "Mary",
		"lastName":  "Banda",
		"gender":    "Female",
		"height":    168.0,
		"weight":    62.0,
	}

	mid, err := model.AddRecord(data)
	if err != nil {
		t.Fatal(err)
	}

	if mid == nil {
		t.Fatal("Test failed. Got nil id")
	}

	row := db.QueryRow(fmt.Sprintf(`SELECT
		id,
		firstName,
		lastName,
		gender,
		height,
		weight
	FROM %s WHERE id=?`, tableName), *mid)

	var id int64
	var weight, height float64
	var firstName,
		lastName,
		gender string

	err = row.Scan(&id,
		&firstName,
		&lastName,
		&gender,
		&height,
		&weight,
	)
	if err != nil {
		t.Fatal(err)
	}

	if firstName != data["firstName"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["firstName"], firstName)
	}
	if lastName != data["lastName"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["lastName"], lastName)
	}
	if gender != data["gender"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["gender"], gender)
	}
	if height != data["height"].(float64) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", data["height"], height)
	}
	if weight != data["weight"].(float64) {
		t.Fatalf("Test failed. Expected: %v; Actual: %v", data["weight"], weight)
	}
}
