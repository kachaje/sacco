package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type Member struct {
	ID                int
	FirstName         string
	LastName          string
	OtherName         string
	Gender            string
	Title             string
	MaritalStatus     string
	DateOfBirth       string
	NationalId        string
	UtilityBillType   string
	UtilityBillNumber string

	db *sql.DB
}

func NewMember(db *sql.DB) *Member {
	return &Member{
		db: db,
	}
}

func (m *Member) AddMember(data map[string]any) (int, error) {
	var id int

	payload, err := json.Marshal(data)
	if err != nil {
		return id, err
	}

	err = json.Unmarshal(payload, m)
	if err != nil {
		return id, err
	}

	result, err := m.db.ExecContext(
		context.Background(),
		`INSERT INTO member (
			firstName,
			lastName,
			otherName,
			gender,
			title,
			maritalStatus,
			dateOfBirth,
			nationalId,
			utilityBillType,
			utilityBillNumber
		) VALUES (
		 	?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)`,
		m.FirstName, m.LastName, m.OtherName,
		m.Gender, m.Title, m.MaritalStatus,
		m.DateOfBirth, m.NationalId, m.UtilityBillType,
		m.UtilityBillNumber,
	)
	if err != nil {
		return id, err
	}

	fmt.Println(result)

	return id, nil
}
