package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type MemberBeneficiary struct {
	ID         int64   `json:"id"`
	MemberId   int64   `json:"memberId"`
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`
	Contact    string  `json:"contact"`

	db *sql.DB
}

func NewMemberBeneficiary(db *sql.DB, memberId *int64) *MemberBeneficiary {
	m := &MemberBeneficiary{
		db: db,
	}

	if memberId != nil {
		m.MemberId = *memberId
	}

	return m
}

func (m *MemberBeneficiary) AddMemberBeneficiary(data map[string]any) (int64, error) {
	var id int64

	payload, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(payload, m)
	if err != nil {
		return 0, err
	}

	result, err := QueryWithRetry(
		m.db,
		context.Background(), 0,
		`INSERT INTO memberBeneficiary (
			memberId,
			name,
			percentage,
			contact
		) VALUES (
		 	?, ?, ?, ?
		)`,
		m.MemberId, m.Name,
		m.Percentage, m.Contact,
	)
	if err != nil {
		return 0, err
	}

	if id, err = result.LastInsertId(); err != nil {
		return 0, err
	}

	return id, nil
}

func (m *MemberBeneficiary) UpdateMemberBeneficiary(data map[string]any, id int64) error {
	fields := []string{}
	values := []any{}

	for key, value := range data {
		fields = append(fields, fmt.Sprintf("%s = ?", key))
		values = append(values, value)
	}

	values = append(values, id)

	statement := fmt.Sprintf("UPDATE memberBeneficiary SET %s WHERE id=?", strings.Join(fields, ", "))

	_, err := QueryWithRetry(
		m.db,
		context.Background(), 0,
		statement, values...,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *MemberBeneficiary) loadRow(row any) (*MemberBeneficiary, bool, error) {
	var id int64
	var memberId int64
	var name,
		percentage,
		contact any
	var err error

	val, ok := row.(*sql.Row)
	if ok {
		err = val.Scan(
			&id,
			&memberId,
			&name,
			&percentage,
			&contact,
		)
	} else {
		val, ok := row.(*sql.Rows)
		if ok {
			err = val.Scan(
				&id,
				&memberId,
				&name,
				&percentage,
				&contact,
			)
		}
	}
	if err != nil {
		return nil, false, fmt.Errorf("memberBeneficiary.loadRow.1: %s", err.Error())
	}

	record := MemberBeneficiary{
		ID:       id,
		MemberId: memberId,
	}

	atLeastOneFieldAdded := false

	if name != nil {
		value := fmt.Sprintf("%v", name)
		if value != "" {
			atLeastOneFieldAdded = true
			record.Name = value
		}
	}
	if percentage != nil {
		value := percentage.(float64)
		if value != 0 {
			atLeastOneFieldAdded = true
			record.Percentage = value
		}
	}
	if contact != nil {
		value := fmt.Sprintf("%v", contact)
		if value != "" {
			atLeastOneFieldAdded = true
			record.Contact = value
		}
	}

	return &record, atLeastOneFieldAdded, nil
}

func (m *MemberBeneficiary) FetchMemberBeneficiary(id int64) (*MemberBeneficiary, error) {

	row := m.db.QueryRow(`SELECT 
		id,
		memberId,
		name,
		percentage,
		contact
	FROM memberBeneficiary WHERE id=? AND active=1`, id)

	record, found, err := m.loadRow(row)
	if err != nil {
		return nil, fmt.Errorf("memberBeneficiary.FetchMemberBeneficiary.1: %s", err.Error())
	}

	if !found {
		return nil, nil
	}

	return record, nil
}

func (m *MemberBeneficiary) FilterBy(whereStatement string) ([]MemberBeneficiary, error) {
	results := []MemberBeneficiary{}

	if !regexp.MustCompile("active").MatchString(whereStatement) {
		whereStatement = fmt.Sprintf("%s AND active=1", whereStatement)
	}

	rows, err := m.db.QueryContext(
		context.Background(),
		fmt.Sprintf(`SELECT
			id,
			memberId,
			name,
			percentage,
			contact
		FROM memberBeneficiary %s`, whereStatement))
	if err != nil {
		return nil, fmt.Errorf("memberBeneficiary.FilterBy.1: %s", err.Error())
	}

	for rows.Next() {
		record, found, err := m.loadRow(rows)
		if err != nil {
			return nil, fmt.Errorf("memberBeneficiary.FetchMemberBeneficiary.1: %s", err.Error())
		}

		if !found {
			continue
		}

		results = append(results, *record)
	}

	return results, nil
}
