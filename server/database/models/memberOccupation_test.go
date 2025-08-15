package models_test

import (
	"fmt"
	"os"
	"sacco/server/database"
	"sacco/server/database/models"
	"testing"
)

func TestAddMemberOccupation(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberOccupation(db.DB, nil)

	data := map[string]any{
		"memberId":             1,
		"employerName":         "TEXT",
		"netPay":               100,
		"jobTitle":             "TEXT",
		"employerAddress":      "TEXT",
		"highestQualification": "TEXT",
	}

	id, err := m.AddMemberOccupation(data)
	if err != nil {
		t.Fatal(err)
	}

	row := db.DB.QueryRow(`SELECT * FROM memberOccupation WHERE id=?`, id)

	var memberId int64
	var netPay float64
	var employerName,
		jobTitle,
		employerAddress,
		highestQualification string

	err = row.Scan(
		&id,
		&memberId,
		&employerName,
		&netPay,
		&jobTitle,
		&employerAddress,
		&highestQualification,
	)
	if err != nil {
		t.Fatal(err)
	}

	if os.Getenv("DEBUG") == "true" {
		fmt.Println(
			id,
			memberId,
			employerName,
			netPay,
			jobTitle,
			employerAddress,
			highestQualification,
		)
	}
}

func TestUpdateMemberOccupation(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberOccupation(db.DB, nil)

	fields := []any{
		1,
		"Sample",
		"Boss",
		"Degree",
	}

	result, err := db.DB.Exec(`INSERT INTO memberOccupation (memberId, employerName, jobTitle, highestQualification) VALUES (?, ?, ?, ?)`, fields...)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]any{
		"employerName":         "Sobo",
		"jobTitle":             "Supervisor",
		"highestQualification": "Diploma",
	}

	err = m.UpdateMemberOccupation(data, id)
	if err != nil {
		t.Fatal(err)
	}

	row := db.DB.QueryRow(`SELECT id, employerName, jobTitle, highestQualification FROM memberOccupation WHERE id=?`, id)

	var employerName, jobTitle, highestQualification string

	err = row.Scan(&id, &employerName, &jobTitle, &highestQualification)
	if err != nil {
		t.Fatal(err)
	}

	if employerName != data["employerName"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["employerName"], employerName)
	}

	if jobTitle != data["jobTitle"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["jobTitle"], jobTitle)
	}

	if highestQualification != data["highestQualification"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["highestQualification"], highestQualification)
	}
}

func TestFetchMemberOccupation(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberOccupation(db.DB, nil)

	fields := []any{
		1,
		"Sample",
		"Boss",
		"Degree",
	}

	result, err := db.DB.Exec(`INSERT INTO memberOccupation (memberId, employerName, jobTitle, highestQualification) VALUES (?, ?, ?, ?)`, fields...)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	memberOccupation, err := m.FetchMemberOccupation(id)
	if err != nil {
		t.Fatal(err)
	}

	if memberOccupation.EmployerName != fields[1].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[1], memberOccupation.EmployerName)
	}

	if memberOccupation.JobTitle != fields[2].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[2], memberOccupation.JobTitle)
	}

	if memberOccupation.HighestQualification != fields[3].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[3], memberOccupation.HighestQualification)
	}
}

func TestFilterMemberOccupationBy(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberOccupation(db.DB, nil)

	fields := [][]any{
		{
			1,
			"Sample1",
			"Boss1",
			"Degree",
		},
		{
			2,
			"Sample2",
			"Boss2",
			"Degree",
		},
		{
			3,
			"Sample3",
			"Boss3",
			"Masters",
		},
		{
			4,
			"Sample4",
			"Boss4",
			"Diploma",
		},
	}

	for i := range fields {
		_, err := db.DB.Exec(`INSERT INTO memberOccupation (memberId, employerName, jobTitle, highestQualification) VALUES (?, ?, ?, ?)`, fields[i]...)
		if err != nil {
			t.Fatal(err)
		}
	}

	results, err := m.FilterBy(`WHERE jobTitle LIKE "Boss%" AND highestQualification = "Degree"`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", len(results))
	}
}
