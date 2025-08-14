package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
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

func (m *Member) AddMember(data map[string]any) (int64, error) {
	var id int64

	payload, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(payload, m)
	if err != nil {
		return 0, err
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
		return 0, err
	}

	if id, err = result.LastInsertId(); err != nil {
		return 0, err
	}

	return id, nil
}

func (m *Member) UpdateMember(data map[string]any, id int64) error {
	fields := []string{}
	values := []any{}

	for key, value := range data {
		fields = append(fields, fmt.Sprintf("%s = ?", key))
		values = append(values, value)
	}

	values = append(values, id)

	statement := fmt.Sprintf("UPDATE member SET %s WHERE id=?", strings.Join(fields, ", "))

	_, err := m.db.Exec(statement, values...)
	if err != nil {
		return err
	}

	return nil
}
