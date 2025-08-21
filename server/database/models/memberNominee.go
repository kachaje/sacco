package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type MemberNominee struct {
	ID               int64  `json:"id"`
	MemberId         int64  `json:"memberId"`
	NextOfKinName    string `json:"nomineeName"`
	NextOfKinPhone   string `json:"nomineePhone"`
	NextOfKinAddress string `json:"nomineeAddress"`

	db *sql.DB
}

func NewMemberNominee(db *sql.DB, memberId *int64) *MemberNominee {
	m := &MemberNominee{
		db: db,
	}

	if memberId != nil {
		m.MemberId = *memberId
	}

	return m
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

	result, err := QueryWithRetry(
		m.db,
		context.Background(), 0,
		`INSERT INTO memberNominee (
			memberId,
			nomineeName,
			nomineePhone,
			nomineeAddress
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

func (m *MemberNominee) loadRow(row any) (*MemberNominee, bool, error) {
	var id int64
	var memberId int64
	var nomineeName,
		nomineePhone,
		nomineeAddress any
	var err error

	val, ok := row.(*sql.Row)
	if ok {
		err = val.Scan(
			&id,
			&memberId,
			&nomineeName,
			&nomineePhone,
			&nomineeAddress,
		)
	} else {
		val, ok := row.(*sql.Rows)
		if ok {
			err = val.Scan(
				&id,
				&memberId,
				&nomineeName,
				&nomineePhone,
				&nomineeAddress,
			)
		}
	}
	if err != nil {
		return nil, false, fmt.Errorf("memberNominee.loadRow.1: %s", err.Error())
	}

	record := MemberNominee{
		ID:       id,
		MemberId: memberId,
	}

	atLeastOneFieldAdded := false

	if nomineeName != nil {
		value := fmt.Sprintf("%v", nomineeName)
		if value != "" {
			atLeastOneFieldAdded = true
			record.NextOfKinName = value
		}
	}
	if nomineePhone != nil {
		value := fmt.Sprintf("%v", nomineePhone)
		if value != "" {
			atLeastOneFieldAdded = true
			record.NextOfKinPhone = value
		}
	}
	if nomineeAddress != nil {
		value := fmt.Sprintf("%v", nomineeAddress)
		if value != "" {
			atLeastOneFieldAdded = true
			record.NextOfKinAddress = value
		}
	}

	return &record, atLeastOneFieldAdded, nil
}

func (m *MemberNominee) FetchMemberNominee(id int64) (*MemberNominee, error) {
	row := m.db.QueryRow(`SELECT 
		id,
		memberId,
		nomineeName,
		nomineePhone,
		nomineeAddress
	FROM memberNominee WHERE id=? AND active=1`, id)

	record, found, err := m.loadRow(row)
	if err != nil {
		return nil, fmt.Errorf("memberNominee.FetchMemberNominee.1: %s", err.Error())
	}

	if !found {
		return nil, nil
	}

	return record, nil
}

func (m *MemberNominee) FilterBy(whereStatement string) ([]MemberNominee, error) {
	results := []MemberNominee{}

	if !regexp.MustCompile("active").MatchString(whereStatement) {
		whereStatement = fmt.Sprintf("%s AND active=1", whereStatement)
	}

	rows, err := m.db.QueryContext(
		context.Background(),
		fmt.Sprintf(`SELECT
			id,
			memberId,
			nomineeName,
			nomineePhone,
			nomineeAddress
		FROM memberNominee %s`, whereStatement))
	if err != nil {
		return nil, fmt.Errorf("memberNominee.FilterBy.1: %s", err.Error())
	}

	for rows.Next() {
		record, found, err := m.loadRow(rows)
		if err != nil {
			return nil, fmt.Errorf("memberNominee.FetchMemberNominee.1: %s", err.Error())
		}

		if !found {
			return nil, nil
		}

		results = append(results, *record)
	}

	return results, nil
}
