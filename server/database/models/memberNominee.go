package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type MemberNominee struct {
	ID               int64
	MemberId         int64
	NextOfKinName    string
	NextOfKinPhone   string
	NextOfKinAddress string

	db *sql.DB
}

func NewMemberNominee(db *sql.DB) *MemberNominee {
	return &MemberNominee{
		db: db,
	}
}

func (m *MemberNominee) AddMemberNominee(data map[string]any) (int64, error) {
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
		`INSERT INTO memberNominee (
			memberId,
			nextOfKinName,
			nextOfKinPhone,
			nextOfKinAddress
		) VALUES (
		 	?, ?, ?, ?
		)`,
		m.MemberId, m.NextOfKinName,
		m.NextOfKinPhone, m.NextOfKinAddress,
	)
	if err != nil {
		return 0, err
	}

	if id, err = result.LastInsertId(); err != nil {
		return 0, err
	}

	return id, nil
}

func (m *MemberNominee) UpdateMemberNominee(data map[string]any, id int64) error {
	fields := []string{}
	values := []any{}

	for key, value := range data {
		fields = append(fields, fmt.Sprintf("%s = ?", key))
		values = append(values, value)
	}

	values = append(values, id)

	statement := fmt.Sprintf("UPDATE memberNominee SET %s WHERE id=?", strings.Join(fields, ", "))

	_, err := m.db.Exec(statement, values...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MemberNominee) FetchMemberNominee(id int64) (*MemberNominee, error) {

	row := m.db.QueryRow(`SELECT 
		memberId,
		nextOfKinName,
		nextOfKinPhone,
		nextOfKinAddress
	FROM memberNominee WHERE id=?`, id)

	var memberId int64
	var nextOfKinName,
		nextOfKinPhone,
		nextOfKinAddress any

	err := row.Scan(
		&memberId,
		&nextOfKinName,
		&nextOfKinPhone,
		&nextOfKinAddress,
	)
	if err != nil {
		return nil, err
	}

	memberNominee := &MemberNominee{
		ID:       id,
		MemberId: memberId,
	}

	if nextOfKinName != nil {
		memberNominee.NextOfKinName = fmt.Sprintf("%v", nextOfKinName)
	}
	if nextOfKinPhone != nil {
		memberNominee.NextOfKinPhone = fmt.Sprintf("%v", nextOfKinPhone)
	}
	if nextOfKinAddress != nil {
		memberNominee.NextOfKinAddress = fmt.Sprintf("%v", nextOfKinAddress)
	}

	return memberNominee, nil
}

func (m *MemberNominee) FilterBy(whereStatement string) ([]MemberNominee, error) {
	results := []MemberNominee{}

	rows, err := m.db.QueryContext(context.Background(), fmt.Sprintf("SELECT * FROM memberNominee %s", whereStatement))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		var memberId int64
		var nextOfKinName,
			nextOfKinPhone,
			nextOfKinAddress any

		err := rows.Scan(
			&id,
			&memberId,
			&nextOfKinName,
			&nextOfKinPhone,
			&nextOfKinAddress,
		)
		if err != nil {
			return nil, err
		}

		memberNominee := MemberNominee{
			ID:       id,
			MemberId: memberId,
		}

		if nextOfKinName != nil {
			memberNominee.NextOfKinName = fmt.Sprintf("%v", nextOfKinName)
		}
		if nextOfKinPhone != nil {
			memberNominee.NextOfKinPhone = fmt.Sprintf("%v", nextOfKinPhone)
		}
		if nextOfKinAddress != nil {
			memberNominee.NextOfKinAddress = fmt.Sprintf("%v", nextOfKinAddress)
		}

		results = append(results, memberNominee)
	}

	return results, nil
}
