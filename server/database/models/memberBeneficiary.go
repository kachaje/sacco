package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

	result, err := m.db.ExecContext(
		context.Background(),
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

	_, err := m.db.Exec(statement, values...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MemberBeneficiary) FetchMemberBeneficiary(id int64) (*MemberBeneficiary, error) {

	row := m.db.QueryRow(`SELECT 
		memberId,
		name,
		percentage,
		contact
	FROM memberBeneficiary WHERE id=?`, id)

	var memberId int64
	var name,
		percentage,
		contact any

	err := row.Scan(
		&memberId,
		&name,
		&percentage,
		&contact,
	)
	if err != nil {
		return nil, err
	}

	memberBeneficiary := &MemberBeneficiary{
		ID:       id,
		MemberId: memberId,
	}

	if name != nil {
		memberBeneficiary.Name = fmt.Sprintf("%v", name)
	}
	if percentage != nil {
		memberBeneficiary.Percentage = percentage.(float64)
	}
	if contact != nil {
		memberBeneficiary.Contact = fmt.Sprintf("%v", contact)
	}

	return memberBeneficiary, nil
}

func (m *MemberBeneficiary) FilterBy(whereStatement string) ([]MemberBeneficiary, error) {
	results := []MemberBeneficiary{}

	rows, err := m.db.QueryContext(context.Background(), fmt.Sprintf("SELECT * FROM memberBeneficiary %s", whereStatement))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		var memberId int64
		var name,
			percentage,
			contact any

		err := rows.Scan(
			&id,
			&memberId,
			&name,
			&percentage,
			&contact,
		)
		if err != nil {
			return nil, err
		}

		memberBeneficiary := MemberBeneficiary{
			ID:       id,
			MemberId: memberId,
		}

		if name != nil {
			memberBeneficiary.Name = fmt.Sprintf("%v", name)
		}
		if percentage != nil {
			memberBeneficiary.Percentage = percentage.(float64)
		}
		if contact != nil {
			memberBeneficiary.Contact = fmt.Sprintf("%v", contact)
		}

		results = append(results, memberBeneficiary)
	}

	return results, nil
}
