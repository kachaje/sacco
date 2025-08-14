package models_test

import (
	"fmt"
	"os"
	"sacco/server/database"
	"sacco/server/database/models"
	"testing"
)

func TestAddMember(t *testing.T) {
	dbname := "addMember"
	db := database.NewDatabase(dbname)
	defer func() {
		db.Close()

		dbFile := fmt.Sprintf("%s.db", dbname)

		_, err := os.Stat(dbFile)
		if os.IsNotExist(err) {
			t.Fatal("Test failed")
		} else {
			os.Remove(dbFile)
		}
	}()

	m := models.NewMember(db.DB)

	data := map[string]any{
		"firstName":         "TEXT",
		"lastName":          "TEXT",
		"otherName":         "TEXT",
		"gender":            "TEXT",
		"title":             "TEXT",
		"maritalStatus":     "TEXT",
		"dateOfBirth":       "TEXT",
		"nationalId":        "TEXT",
		"utilityBillType":   "TEXT",
		"utilityBillNumber": "TEXT",
	}

	id, err := m.AddMember(data)
	if err != nil {
		t.Fatal(err)
	}

	row := db.DB.QueryRow(`SELECT * FROM member WHERE id=?`, id)

	var firstName,
		lastName,
		otherName,
		gender,
		title,
		maritalStatus,
		dateOfBirth,
		nationalId,
		utilityBillType,
		utilityBillNumber string

	err = row.Scan(&id, &firstName, &lastName, &otherName,
		&gender, &title, &maritalStatus,
		&dateOfBirth, &nationalId, &utilityBillType,
		&utilityBillNumber)
	if err != nil {
		t.Fatal(err)
	}

	if os.Getenv("DEBUG") == "true" {
		fmt.Println(
			id,
			firstName,
			lastName,
			otherName,
			gender,
			title,
			maritalStatus,
			dateOfBirth,
			nationalId,
			utilityBillType,
			utilityBillNumber,
		)
	}
}
