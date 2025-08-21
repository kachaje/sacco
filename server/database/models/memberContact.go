package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type MemberContact struct {
	ID                 *int64 `json:"id"`
	MemberId           *int64 `json:"memberId"`
	PostalAddress      string `json:"postalAddress"`
	ResidentialAddress string `json:"residentialAddress"`
	HomeVillage        string `json:"homeVillage"`
	HomeTA             string `json:"homeTA"`
	HomeDistrict       string `json:"homeDistrict"`

	db *sql.DB
}

func NewMemberContact(db *sql.DB, memberId *int64) *MemberContact {
	m := &MemberContact{
		db: db,
	}

	if memberId != nil {
		m.MemberId = memberId
	}

	return m
}

func (m *MemberContact) AddMemberContact(data map[string]any) (int64, error) {
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
		`INSERT INTO memberContact (
			memberId,
			postalAddress,
			residentialAddress,
			homeVillage,
			homeTA,
			homeDistrict
		) VALUES (
		 	?, ?, ?, ?, ?, ?
		)`,
		*m.MemberId, m.PostalAddress, m.ResidentialAddress,
		m.HomeVillage, m.HomeTA, m.HomeDistrict,
	)
	if err != nil {
		return 0, err
	}

	if id, err = result.LastInsertId(); err != nil {
		return 0, err
	}

	return id, nil
}

func (m *MemberContact) UpdateMemberContact(data map[string]any, id int64) error {
	fields := []string{}
	values := []any{}

	for key, value := range data {
		fields = append(fields, fmt.Sprintf("%s = ?", key))
		values = append(values, value)
	}

	values = append(values, id)

	statement := fmt.Sprintf("UPDATE memberContact SET %s WHERE id=?", strings.Join(fields, ", "))

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

func (m *MemberContact) loadRow(row any) (*MemberContact, bool, error) {
	var id int64
	var memberId int64
	var postalAddress,
		residentialAddress,
		homeVillage,
		homeTA,
		homeDistrict any
	var err error

	val, ok := row.(*sql.Row)
	if ok {
		err = val.Scan(
			&id,
			&memberId,
			&postalAddress,
			&residentialAddress,
			&homeVillage,
			&homeTA,
			&homeDistrict,
		)
	} else {
		val, ok := row.(*sql.Rows)
		if ok {
			err = val.Scan(
				&id,
				&memberId,
				&postalAddress,
				&residentialAddress,
				&homeVillage,
				&homeTA,
				&homeDistrict,
			)
		}
	}
	if err != nil {
		return nil, false, fmt.Errorf("memberContact.loadRow.1: %s", err.Error())
	}

	record := MemberContact{
		ID:       &id,
		MemberId: &memberId,
	}

	atLeastOneFieldAdded := false

	if postalAddress != nil {
		value := fmt.Sprintf("%v", postalAddress)
		if value != "" {
			atLeastOneFieldAdded = true
			record.PostalAddress = value
		}
	}
	if residentialAddress != nil {
		value := fmt.Sprintf("%v", residentialAddress)
		if value != "" {
			atLeastOneFieldAdded = true
			record.ResidentialAddress = value
		}
	}
	if homeVillage != nil {
		value := fmt.Sprintf("%v", homeVillage)
		if value != "" {
			atLeastOneFieldAdded = true
			record.HomeVillage = value
		}
	}
	if homeTA != nil {
		value := fmt.Sprintf("%v", homeTA)
		if value != "" {
			atLeastOneFieldAdded = true
			record.HomeTA = value
		}
	}
	if homeDistrict != nil {
		value := fmt.Sprintf("%v", homeDistrict)
		if value != "" {
			atLeastOneFieldAdded = true
			record.HomeDistrict = value
		}
	}

	return &record, atLeastOneFieldAdded, nil
}

func (m *MemberContact) FetchMemberContact(id int64) (*MemberContact, error) {

	row := m.db.QueryRow(`SELECT 
		id,
		memberId,
		postalAddress,
		residentialAddress,
		homeVillage,
		homeTA,
		homeDistrict
	FROM memberContact WHERE id=? AND active=1`, id)

	record, found, err := m.loadRow(row)
	if err != nil {
		return nil, fmt.Errorf("memberContact.FetchMemberContact.1: %s", err.Error())
	}

	if !found {
		return nil, nil
	}

	return record, nil
}

func (m *MemberContact) FilterBy(whereStatement string) ([]MemberContact, error) {
	results := []MemberContact{}

	if !regexp.MustCompile("active").MatchString(whereStatement) {
		whereStatement = fmt.Sprintf("%s AND active=1", whereStatement)
	}

	rows, err := m.db.QueryContext(
		context.Background(),
		fmt.Sprintf(`SELECT 
				id, 
				memberId,
				postalAddress,
				residentialAddress,
				homeVillage,
				homeTA,
				homeDistrict 
			FROM memberContact %s`,
			whereStatement,
		))
	if err != nil {
		return nil, fmt.Errorf("memberContact.FilterBy.1: %s", err.Error())
	}

	for rows.Next() {
		record, found, err := m.loadRow(rows)
		if err != nil {
			return nil, fmt.Errorf("memberContact.FetchMemberContact.1: %s", err.Error())
		}

		if !found {
			continue
		}

		results = append(results, *record)
	}

	return results, nil
}
