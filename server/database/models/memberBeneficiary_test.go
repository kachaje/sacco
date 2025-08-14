package models_test

import (
	"fmt"
	"os"
	"sacco/server/database"
	"sacco/server/database/models"
	"testing"
)

func TestAddMemberBeneficiary(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberBeneficiary(db.DB)

	data := map[string]any{
		"memberId":   1,
		"name":       "TEXT",
		"percentage": 5,
		"contact":    "TEXT",
	}

	id, err := m.AddMemberBeneficiary(data)
	if err != nil {
		t.Fatal(err)
	}

	row := db.DB.QueryRow(`SELECT * FROM memberBeneficiary WHERE id=?`, id)

	var memberId int64
	var percentage float64
	var name, contact string

	err = row.Scan(
		&id,
		&memberId,
		&name,
		&percentage,
		&contact,
	)
	if err != nil {
		t.Fatal(err)
	}

	if os.Getenv("DEBUG") == "true" {
		fmt.Println(
			id,
			memberId,
			name,
			percentage,
			contact,
		)
	}
}

func TestUpdateMemberBeneficiary(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberBeneficiary(db.DB)

	fields := []any{
		1,
		"Sample",
		5.0,
		"Degree",
	}

	result, err := db.DB.Exec(`INSERT INTO memberBeneficiary (memberId, name, percentage, contact) VALUES (?, ?, ?, ?)`, fields...)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]any{
		"name":       "Sobo",
		"percentage": 10.0,
		"contact":    "Diploma",
	}

	err = m.UpdateMemberBeneficiary(data, id)
	if err != nil {
		t.Fatal(err)
	}

	row := db.DB.QueryRow(`SELECT id, name, percentage, contact FROM memberBeneficiary WHERE id=?`, id)

	var percentage float64
	var name, contact string

	err = row.Scan(&id, &name, &percentage, &contact)
	if err != nil {
		t.Fatal(err)
	}

	if name != data["name"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["name"], name)
	}

	if percentage != data["percentage"].(float64) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["percentage"], percentage)
	}

	if contact != data["contact"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["contact"], contact)
	}
}

func TestFetchMemberBeneficiary(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberBeneficiary(db.DB)

	fields := []any{
		1,
		"Sample",
		8.0,
		"Degree",
	}

	result, err := db.DB.Exec(`INSERT INTO memberBeneficiary (memberId, name, percentage, contact) VALUES (?, ?, ?, ?)`, fields...)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	memberBeneficiary, err := m.FetchMemberBeneficiary(id)
	if err != nil {
		t.Fatal(err)
	}

	if memberBeneficiary.Name != fields[1].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[1], memberBeneficiary.Name)
	}

	if memberBeneficiary.Percentage != fields[2].(float64) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[2], memberBeneficiary.Percentage)
	}

	if memberBeneficiary.Contact != fields[3].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[3], memberBeneficiary.Contact)
	}
}

func TestFilterMemberBeneficiaryBy(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberBeneficiary(db.DB)

	fields := [][]any{
		{
			1,
			"Sample1",
			5.0,
			"Degree",
		},
		{
			2,
			"Sample2",
			5.0,
			"Degree",
		},
		{
			3,
			"Sample3",
			6.0,
			"Masters",
		},
		{
			4,
			"Sample4",
			6.0,
			"Diploma",
		},
	}

	for i := range fields {
		_, err := db.DB.Exec(`INSERT INTO memberBeneficiary (memberId, name, percentage, contact) VALUES (?, ?, ?, ?)`, fields[i]...)
		if err != nil {
			t.Fatal(err)
		}
	}

	results, err := m.FilterBy(`WHERE contact = "Degree"`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", len(results))
	}
}
