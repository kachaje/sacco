package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type Member struct {
	ID                 int64  `json:"id"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	OtherName          string `json:"otherName"`
	Gender             string `json:"gender"`
	Title              string `json:"title"`
	MaritalStatus      string `json:"maritalStatus"`
	DateOfBirth        string `json:"dateOfBirth"`
	NationalId         string `json:"nationalId"`
	UtilityBillType    string `json:"utilityBillType"`
	UtilityBillNumber  string `json:"utilityBillNumber"`
	FileNumber         string `json:"fileNumber"`
	OldFileNumber      string `json:"oldFileNumber"`
	DefaultPhoneNumber string `json:"defaultPhoneNumber"`

	Beneficiaries     []MemberBeneficiary `json:"beneficiaries"`
	ContactDetails    *MemberContact      `json:"contactDetails"`
	Nominee           *MemberNominee      `json:"nominee"`
	OccupationDetails *MemberOccupation   `json:"occupationDetails"`

	db *sql.DB
}

func NewMember(db *sql.DB) *Member {
	return &Member{
		db: db,
	}
}

func (m *Member) MemberDetails(memberId int64) (map[string]any, error) {
	fullRecord := map[string]any{}

	member, err := m.FetchMember(memberId)
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf(`WHERE memberId = %d`, memberId)

	c := NewMemberContact(m.db, &memberId)

	contactDetails, err := c.FilterBy(filter)
	if err != nil {
		return nil, err
	}

	if len(contactDetails) > 0 {
		member.ContactDetails = &contactDetails[0]
	}

	n := NewMemberNominee(m.db, &memberId)

	nominee, err := n.FilterBy(filter)
	if err != nil {
		return nil, err
	}

	if len(nominee) > 0 {
		member.Nominee = &nominee[0]
	}

	o := NewMemberOccupation(m.db, &memberId)

	occupation, err := o.FilterBy(filter)
	if err != nil {
		return nil, err
	}

	if len(occupation) > 0 {
		member.OccupationDetails = &occupation[0]
	}

	b := NewMemberBeneficiary(m.db, &memberId)

	beneficiaries, err := b.FilterBy(filter)
	if err != nil {
		return nil, err
	}

	if len(beneficiaries) > 0 {
		member.Beneficiaries = beneficiaries
	}

	payload, err := json.Marshal(member)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(payload, &fullRecord)
	if err != nil {
		return nil, err
	}

	return fullRecord, nil
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
			utilityBillNumber,
			fileNumber,
			oldFileNumber,
			defaultPhoneNumber
		) VALUES (
		 	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)`,
		m.FirstName, m.LastName, m.OtherName,
		m.Gender, m.Title, m.MaritalStatus,
		m.DateOfBirth, m.NationalId, m.UtilityBillType,
		m.UtilityBillNumber, m.FileNumber, m.OldFileNumber,
		m.DefaultPhoneNumber,
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

	fmt.Println("##########", fields, values)

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
		utilityBillNumber,
		fileNumber,
		oldFileNumber,
		defaultPhoneNumber
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
		utilityBillNumber,
		fileNumber,
		oldFileNumber,
		defaultPhoneNumber any

	err := row.Scan(&firstName, &lastName, &otherName,
		&gender, &title, &maritalStatus,
		&dateOfBirth, &nationalId, &utilityBillType,
		&utilityBillNumber, &fileNumber, &oldFileNumber,
		&defaultPhoneNumber)
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
	if fileNumber != nil {
		member.FileNumber = fmt.Sprintf("%v", fileNumber)
	}
	if oldFileNumber != nil {
		member.OldFileNumber = fmt.Sprintf("%v", oldFileNumber)
	}
	if defaultPhoneNumber != nil {
		member.DefaultPhoneNumber = fmt.Sprintf("%v", defaultPhoneNumber)
	}

	return member, nil
}

func (m *Member) FilterBy(whereStatement string) ([]Member, error) {
	results := []Member{}

	rows, err := m.db.QueryContext(
		context.Background(),
		fmt.Sprintf(`SELECT 
				id,
				firstName,
				lastName,
				otherName,
				gender,
				title,
				maritalStatus,
				dateOfBirth,
				nationalId,
				utilityBillType,
				utilityBillNumber,
				fileNumber,
				oldFileNumber,
				defaultPhoneNumber
			FROM member %s`,
			whereStatement,
		))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		var firstName,
			lastName,
			otherName,
			gender,
			title,
			maritalStatus,
			dateOfBirth,
			nationalId,
			utilityBillType,
			utilityBillNumber,
			fileNumber,
			oldFileNumber,
			defaultPhoneNumber any

		err := rows.Scan(&id, &firstName, &lastName, &otherName,
			&gender, &title, &maritalStatus,
			&dateOfBirth, &nationalId, &utilityBillType,
			&utilityBillNumber, &fileNumber, &oldFileNumber,
			&defaultPhoneNumber)
		if err != nil {
			return nil, err
		}

		member := Member{
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
		if fileNumber != nil {
			member.FileNumber = fmt.Sprintf("%v", fileNumber)
		}
		if oldFileNumber != nil {
			member.OldFileNumber = fmt.Sprintf("%v", oldFileNumber)
		}
		if defaultPhoneNumber != nil {
			member.DefaultPhoneNumber = fmt.Sprintf("%v", defaultPhoneNumber)
		}

		results = append(results, member)
	}

	return results, nil
}

func (m *Member) FetchMemberByPhoneNumber(phoneNumber string) (*Member, error) {

	row := m.db.QueryRow(`SELECT 
		id,
		firstName,
		lastName,
		otherName,
		gender,
		title,
		maritalStatus,
		dateOfBirth,
		nationalId,
		utilityBillType,
		utilityBillNumber,
		fileNumber,
		oldFileNumber,
		defaultPhoneNumber
	FROM member WHERE defaultPhoneNumber=?`, phoneNumber)

	var id int64
	var firstName,
		lastName,
		otherName,
		gender,
		title,
		maritalStatus,
		dateOfBirth,
		nationalId,
		utilityBillType,
		utilityBillNumber,
		fileNumber,
		oldFileNumber,
		defaultPhoneNumber any

	err := row.Scan(&id, &firstName, &lastName, &otherName,
		&gender, &title, &maritalStatus,
		&dateOfBirth, &nationalId, &utilityBillType,
		&utilityBillNumber, &fileNumber, &oldFileNumber, &defaultPhoneNumber)
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
	if fileNumber != nil {
		member.FileNumber = fmt.Sprintf("%v", fileNumber)
	}
	if oldFileNumber != nil {
		member.OldFileNumber = fmt.Sprintf("%v", oldFileNumber)
	}
	if defaultPhoneNumber != nil {
		member.DefaultPhoneNumber = fmt.Sprintf("%v", defaultPhoneNumber)
	}

	return member, nil
}
