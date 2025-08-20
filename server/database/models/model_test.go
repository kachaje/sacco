package models_test

import (
	"database/sql"
	"os"
	"sacco/server/database/models"
	"testing"
)

func setupDb(dbname string) (*sql.DB, *models.Model, error) {
	db, err := sql.Open("sqlite", dbname)
	if err != nil {
		return nil, nil, err
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS person (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		firstName TEXT,
		lastName TEXT,
		gender TEXT,
		height REAL,
		weight REAL
	);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, nil, err
	}

	fields := []string{"firstName", "lastName", "gender", "height", "weight"}
	fieldTypes := map[string]string{
		"firstName": "string",
		"lastName":  "string",
		"gender":    "string",
		"height":    "float",
		"weight":    "weight",
	}

	model := models.NewModel(db, fields, fieldTypes)

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
