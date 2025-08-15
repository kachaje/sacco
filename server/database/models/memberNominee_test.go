package models_test

import (
	"fmt"
	"os"
	"sacco/server/database"
	"sacco/server/database/models"
	"testing"
)

func TestAddMemberNominee(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberNominee(db.DB, nil)

	data := map[string]any{
		"memberId":         1,
		"nextOfKinName":    "TEXT",
		"nextOfKinPhone":   "TEXT",
		"nextOfKinAddress": "TEXT",
	}

	id, err := m.AddMemberNominee(data)
	if err != nil {
		t.Fatal(err)
	}

	row := db.DB.QueryRow(`SELECT * FROM memberNominee WHERE id=?`, id)

	var memberId int64
	var nextOfKinName,
		nextOfKinPhone,
		nextOfKinAddress string

	err = row.Scan(
		&id,
		&memberId,
		&nextOfKinName,
		&nextOfKinPhone,
		&nextOfKinAddress,
	)
	if err != nil {
		t.Fatal(err)
	}

	if os.Getenv("DEBUG") == "true" {
		fmt.Println(
			id,
			memberId,
			nextOfKinName,
			nextOfKinPhone,
			nextOfKinAddress,
		)
	}
}

func TestUpdateMemberNominee(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberNominee(db.DB, nil)

	fields := []any{
		1,
		"Sample",
		"Boss",
		"Degree",
	}

	result, err := db.DB.Exec(`INSERT INTO memberNominee (memberId, nextOfKinName, nextOfKinPhone, nextOfKinAddress) VALUES (?, ?, ?, ?)`, fields...)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]any{
		"nextOfKinName":    "Sobo",
		"nextOfKinPhone":   "Supervisor",
		"nextOfKinAddress": "Diploma",
	}

	err = m.UpdateMemberNominee(data, id)
	if err != nil {
		t.Fatal(err)
	}

	row := db.DB.QueryRow(`SELECT id, nextOfKinName, nextOfKinPhone, nextOfKinAddress FROM memberNominee WHERE id=?`, id)

	var nextOfKinName, nextOfKinPhone, nextOfKinAddress string

	err = row.Scan(&id, &nextOfKinName, &nextOfKinPhone, &nextOfKinAddress)
	if err != nil {
		t.Fatal(err)
	}

	if nextOfKinName != data["nextOfKinName"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["nextOfKinName"], nextOfKinName)
	}

	if nextOfKinPhone != data["nextOfKinPhone"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["nextOfKinPhone"], nextOfKinPhone)
	}

	if nextOfKinAddress != data["nextOfKinAddress"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["nextOfKinAddress"], nextOfKinAddress)
	}
}

func TestFetchMemberNominee(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberNominee(db.DB, nil)

	fields := []any{
		1,
		"Sample",
		"Boss",
		"Degree",
	}

	result, err := db.DB.Exec(`INSERT INTO memberNominee (memberId, nextOfKinName, nextOfKinPhone, nextOfKinAddress) VALUES (?, ?, ?, ?)`, fields...)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	memberNominee, err := m.FetchMemberNominee(id)
	if err != nil {
		t.Fatal(err)
	}

	if memberNominee.NextOfKinName != fields[1].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[1], memberNominee.NextOfKinName)
	}

	if memberNominee.NextOfKinPhone != fields[2].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[2], memberNominee.NextOfKinPhone)
	}

	if memberNominee.NextOfKinAddress != fields[3].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[3], memberNominee.NextOfKinAddress)
	}
}

func TestFilterMemberNomineeBy(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberNominee(db.DB, nil)

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
		_, err := db.DB.Exec(`INSERT INTO memberNominee (memberId, nextOfKinName, nextOfKinPhone, nextOfKinAddress) VALUES (?, ?, ?, ?)`, fields[i]...)
		if err != nil {
			t.Fatal(err)
		}
	}

	results, err := m.FilterBy(`WHERE nextOfKinPhone LIKE "Boss%" AND nextOfKinAddress = "Degree"`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", len(results))
	}
}
