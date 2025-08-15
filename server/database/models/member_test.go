package models_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sacco/server/database"
	"sacco/server/database/models"
	"testing"
)

func TestAddMember(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

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
		"fileNumber":        "TEXT",
		"oldFileNumber":     "TEXT",
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
		utilityBillNumber,
		fileNumber,
		oldFileNumber string

	err = row.Scan(&id, &firstName, &lastName, &otherName,
		&gender, &title, &maritalStatus,
		&dateOfBirth, &nationalId, &utilityBillType,
		&utilityBillNumber, &fileNumber, &oldFileNumber)
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
			fileNumber,
			oldFileNumber,
		)
	}
}

func TestUpdateMember(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMember(db.DB)

	fields := []any{
		"Mary",
		"Banda",
		"Female",
	}

	result, err := db.DB.Exec(`INSERT INTO member (firstName, lastName, gender) VALUES (?, ?, ?)`, fields...)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]any{
		"firstName": "John",
		"lastName":  "Bandawe",
		"gender":    "Male",
	}

	err = m.UpdateMember(data, id)
	if err != nil {
		t.Fatal(err)
	}

	row := db.DB.QueryRow(`SELECT id, firstName, lastName, gender FROM member WHERE id=?`, id)

	var firstName,
		lastName,
		gender string

	err = row.Scan(&id, &firstName, &lastName, &gender)
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
}

func TestFetchMember(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMember(db.DB)

	fields := []any{
		"Mary",
		"Banda",
		"Female",
	}

	result, err := db.DB.Exec(`INSERT INTO member (firstName, lastName, gender) VALUES (?, ?, ?)`, fields...)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	member, err := m.FetchMember(id)
	if err != nil {
		t.Fatal(err)
	}

	if member.FirstName != fields[0].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[0], member.FirstName)
	}

	if member.LastName != fields[1].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[1], member.LastName)
	}

	if member.Gender != fields[2].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[2], member.Gender)
	}
}

func TestMemberFilterBy(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMember(db.DB)

	fields := [][]any{
		{
			"Mary",
			"Banda",
			"Female",
		},
		{
			"John",
			"Bongwe",
			"Male",
		},
		{
			"Paul",
			"Bandawe",
			"Male",
		},
		{
			"Peter",
			"Banda",
			"Male",
		},
	}

	for i := range fields {
		_, err := db.DB.Exec(`INSERT INTO member (firstName, lastName, gender) VALUES (?, ?, ?)`, fields[i]...)
		if err != nil {
			t.Fatal(err)
		}
	}

	results, err := m.FilterBy(`WHERE lastName LIKE "Banda%" AND gender = "Male"`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", len(results))
	}
}

func TestMemberDetails(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMember(db.DB)

	content, err := os.ReadFile(filepath.Join(".", "fixtures", "member.sql"))
	if err != nil {
		t.Fatal(err)
	}

	sqlStatement := string(content)

	_, err = db.DB.Exec(sqlStatement)
	if err != nil {
		t.Fatal(err)
	}

	var id int64 = 10

	result, err := m.MemberDetails(id)
	if err != nil {
		t.Fatal(err)
	}

	payload, _ := json.MarshalIndent(result, "", "  ")

	fmt.Println(string(payload))
}
