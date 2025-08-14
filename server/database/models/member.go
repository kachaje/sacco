package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type Member struct {
	ID                int64
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

func (m *Member) FetchMember(id int64) (*Member, error) {

	row := m.db.QueryRow(`SELECT 
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
	FROM member WHERE id=?`, id)

	var firstName,
		lastName,
		otherName,
		gender,
		title,
		maritalStatus,
		dateOfBirth,
		nationalId,
		utilityBillType,
		utilityBillNumber any

	err := row.Scan(&firstName, &lastName, &otherName,
		&gender, &title, &maritalStatus,
		&dateOfBirth, &nationalId, &utilityBillType,
		&utilityBillNumber)
	if err != nil {
		return nil, err
	}

	member := &Member{
		ID: id,
	}

	if firstName != nil {
		member.FirstName = fmt.Sprintf("%v", firstName)
	}
	if lastName != nil {
		member.LastName = fmt.Sprintf("%v", lastName)
	}
	if otherName != nil {
		member.OtherName = fmt.Sprintf("%v", otherName)
	}
	if gender != nil {
		member.Gender = fmt.Sprintf("%v", gender)
	}
	if title != nil {
		member.Title = fmt.Sprintf("%v", title)
	}
	if maritalStatus != nil {
		member.MaritalStatus = fmt.Sprintf("%v", maritalStatus)
	}
	if dateOfBirth != nil {
		member.DateOfBirth = fmt.Sprintf("%v", dateOfBirth)
	}
	if nationalId != nil {
		member.NationalId = fmt.Sprintf("%v", nationalId)
	}
	if utilityBillType != nil {
		member.UtilityBillType = fmt.Sprintf("%v", utilityBillType)
	}
	if utilityBillNumber != nil {
		member.UtilityBillNumber = fmt.Sprintf("%v", utilityBillNumber)
	}

	return member, nil
}
