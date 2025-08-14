package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type MemberContact struct {
	ID                 int64
	MemberId           int64
	PostalAddress      string
	ResidentialAddress string
	PhoneNumber        string
	HomeVillage        string
	HomeTA             string
	HomeDistrict       string

	db *sql.DB
}

func NewMemberContact(db *sql.DB) *MemberContact {
	return &MemberContact{
		db: db,
	}
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

	result, err := m.db.ExecContext(
		context.Background(),
		`INSERT INTO memberContact (
			memberId,
			postalAddress,
			residentialAddress,
			phoneNumber,
			homeVillage,
			homeTA,
			homeDistrict
		) VALUES (
		 	?, ?, ?, ?, ?, ?, ?
		)`,
		m.MemberId, m.PostalAddress, m.ResidentialAddress,
		m.PhoneNumber, m.HomeVillage, m.HomeTA, m.HomeDistrict,
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

	_, err := m.db.Exec(statement, values...)
	if err != nil {
		return err
	}

	return nil
}

func (m *MemberContact) FetchMemberContact(id int64) (*MemberContact, error) {

	row := m.db.QueryRow(`SELECT 
		memberId,
		postalAddress,
		residentialAddress,
		phoneNumber,
		homeVillage,
		homeTA,
		homeDistrict
	FROM memberContact WHERE id=?`, id)

	var memberId int64
	var postalAddress,
		residentialAddress,
		phoneNumber,
		homeVillage,
		homeTA,
		homeDistrict any

	err := row.Scan(
		&memberId,
		&postalAddress,
		&residentialAddress,
		&phoneNumber,
		&homeVillage,
		&homeTA,
		&homeDistrict,
	)
	if err != nil {
		return nil, err
	}

	memberContact := &MemberContact{
		ID:       id,
		MemberId: memberId,
	}

	if postalAddress != nil {
		memberContact.PostalAddress = fmt.Sprintf("%v", postalAddress)
	}
	if residentialAddress != nil {
		memberContact.ResidentialAddress = fmt.Sprintf("%v", residentialAddress)
	}
	if phoneNumber != nil {
		memberContact.PhoneNumber = fmt.Sprintf("%v", phoneNumber)
	}
	if homeVillage != nil {
		memberContact.HomeVillage = fmt.Sprintf("%v", homeVillage)
	}
	if homeTA != nil {
		memberContact.HomeTA = fmt.Sprintf("%v", homeTA)
	}
	if homeDistrict != nil {
		memberContact.HomeDistrict = fmt.Sprintf("%v", homeDistrict)
	}

	return memberContact, nil
}

func (m *MemberContact) FilterBy(whereStatement string) ([]MemberContact, error) {
	results := []MemberContact{}

	rows, err := m.db.QueryContext(context.Background(), fmt.Sprintf("SELECT * FROM memberContact %s", whereStatement))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		var memberId int64
		var postalAddress,
			residentialAddress,
			phoneNumber,
			homeVillage,
			homeTA,
			homeDistrict any

		err := rows.Scan(
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
			return nil, err
		}

		memberContact := MemberContact{
			ID:       id,
			MemberId: memberId,
		}

		if postalAddress != nil {
			memberContact.PostalAddress = fmt.Sprintf("%v", postalAddress)
		}
		if residentialAddress != nil {
			memberContact.ResidentialAddress = fmt.Sprintf("%v", residentialAddress)
		}
		if phoneNumber != nil {
			memberContact.PhoneNumber = fmt.Sprintf("%v", phoneNumber)
		}
		if homeVillage != nil {
			memberContact.HomeVillage = fmt.Sprintf("%v", homeVillage)
		}
		if homeTA != nil {
			memberContact.HomeTA = fmt.Sprintf("%v", homeTA)
		}
		if homeDistrict != nil {
			memberContact.HomeDistrict = fmt.Sprintf("%v", homeDistrict)
		}

		results = append(results, memberContact)
	}

	return results, nil
}
