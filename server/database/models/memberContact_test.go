package models_test

import (
	"fmt"
	"os"
	"sacco/server/database"
	"sacco/server/database/models"
	"testing"
)

func TestAddMemberContact(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberContact(db.DB, nil)

	data := map[string]any{
		"memberId":           1,
		"postalAddress":      "TEXT",
		"residentialAddress": "TEXT",
		"phoneNumber":        "TEXT",
		"homeVillage":        "TEXT",
		"homeTA":             "TEXT",
		"homeDistrict":       "TEXT",
	}

	id, err := m.AddMemberContact(data)
	if err != nil {
		t.Fatal(err)
	}

	row := db.DB.QueryRow(`SELECT * FROM memberContact WHERE id=?`, id)

	var memberId int64
	var postalAddress,
		residentialAddress,
		phoneNumber,
		homeVillage,
		homeTA,
		homeDistrict string

	err = row.Scan(
		&id,
		&memberId,
		&postalAddress,
		&residentialAddress,
		&phoneNumber,
		&homeVillage,
		&homeTA,
		&homeDistrict,
	)
	if err != nil {
		t.Fatal(err)
	}

	if os.Getenv("DEBUG") == "true" {
		fmt.Println(
			id,
			memberId,
			postalAddress,
			residentialAddress,
			phoneNumber,
			homeVillage,
			homeTA,
			homeDistrict,
		)
	}
}

func TestUpdateMemberContact(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberContact(db.DB, nil)

	fields := []any{
		1,
		"Mudzi",
		"Mfumu",
		"Boma",
	}

	result, err := db.DB.Exec(`INSERT INTO memberContact (memberId, homeVillage, homeTA, homeDistrict) VALUES (?, ?, ?, ?)`, fields...)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]any{
		"homeVillage":  "Makwinja",
		"homeTA":       "Kalolo",
		"homeDistrict": "Lilongwe",
	}

	err = m.UpdateMemberContact(data, id)
	if err != nil {
		t.Fatal(err)
	}

	row := db.DB.QueryRow(`SELECT id, homeVillage, homeTA, homeDistrict FROM memberContact WHERE id=?`, id)

	var homeVillage, homeTA, homeDistrict string

	err = row.Scan(&id, &homeVillage, &homeTA, &homeDistrict)
	if err != nil {
		t.Fatal(err)
	}

	if homeVillage != data["homeVillage"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["homeVillage"], homeVillage)
	}

	if homeTA != data["homeTA"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["homeTA"], homeTA)
	}

	if homeDistrict != data["homeDistrict"].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", data["homeDistrict"], homeDistrict)
	}
}

func TestFetchMemberContact(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberContact(db.DB, nil)

	fields := []any{
		1,
		"Mudzi",
		"Mfumu",
		"Boma",
	}

	result, err := db.DB.Exec(`INSERT INTO memberContact (memberId, homeVillage, homeTA, homeDistrict) VALUES (?, ?, ?, ?)`, fields...)
	if err != nil {
		t.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	memberContact, err := m.FetchMemberContact(id)
	if err != nil {
		t.Fatal(err)
	}

	if memberContact.HomeVillage != fields[1].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[1], memberContact.HomeVillage)
	}

	if memberContact.HomeTA != fields[2].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[2], memberContact.HomeTA)
	}

	if memberContact.HomeDistrict != fields[3].(string) {
		t.Fatalf("Test failed. Expected: %s; Actual: %v", fields[3], memberContact.HomeDistrict)
	}
}

func TestFilterMemberContactBy(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

	m := models.NewMemberContact(db.DB, nil)

	fields := [][]any{
		{
			1,
			"Mary",
			"Banda",
			"Female",
		},
		{
			2,
			"John",
			"Bongwe",
			"Male",
		},
		{
			3,
			"Paul",
			"Bandawe",
			"Male",
		},
		{
			4,
			"Peter",
			"Banda",
			"Male",
		},
	}

	for i := range fields {
		_, err := db.DB.Exec(`INSERT INTO memberContact (memberId, homeVillage, homeTA, homeDistrict) VALUES (?, ?, ?, ?)`, fields[i]...)
		if err != nil {
			t.Fatal(err)
		}
	}

	results, err := m.FilterBy(`WHERE homeTA LIKE "Banda%" AND homeDistrict = "Male"`)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 2 {
		t.Fatalf("Test failed. Expected: 2; Actual: %v", len(results))
	}
}
