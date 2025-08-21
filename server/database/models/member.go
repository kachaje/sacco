package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Member struct {
	ID                int64  `json:"id"`
	FirstName         string `json:"firstName"`
	LastName          string `json:"lastName"`
	OtherName         string `json:"otherName"`
	Gender            string `json:"gender"`
	Title             string `json:"title"`
	MaritalStatus     string `json:"maritalStatus"`
	DateOfBirth       string `json:"dateOfBirth"`
	NationalId        string `json:"nationalId"`
	UtilityBillType   string `json:"utilityBillType"`
	UtilityBillNumber string `json:"utilityBillNumber"`
	FileNumber        string `json:"fileNumber"`
	OldFileNumber     string `json:"oldFileNumber"`
	PhoneNumber       string `json:"phoneNumber"`
	MemberIdNumber    string `json:"memberIdNumber"`
	ShortMemberId     string `json:"shortMemberId"`
	DateJoined        string `json:"dateJoined"`

	Beneficiaries     []MemberBeneficiary `json:"beneficiaries"`
	ContactDetails    *MemberContact      `json:"contactDetails"`
	Nominee           *MemberNominee      `json:"nomineeDetails"`
	OccupationDetails *MemberOccupation   `json:"occupationDetails"`

	validFields []string
	db          *sql.DB
}

func NewMember(db *sql.DB) *Member {
	return &Member{
		db: db,
		validFields: []string{
			"id", "firstName", "lastName", "otherName", "gender",
			"title", "maritalStatus", "dateOfBirth", "nationalId",
			"utilityBillType", "utilityBillNumber", "fileNumber",
			"oldFileNumber", "phoneNumber",
			"shortMemberId", "memberIdNumber", "dateJoined",
		},
	}
}

func (m *Member) MemberDetails(memberId int64) (map[string]any, error) {
	fullRecord := map[string]any{}

	member, err := m.FetchMember(memberId)
	if err != nil {
		return nil, fmt.Errorf("member.MemberDetails.1: %s", err.Error())
	}

	filter := fmt.Sprintf(`WHERE memberId = %d`, memberId)

	c := NewMemberContact(m.db, &memberId)

	contactDetails, err := c.FilterBy(filter)
	if err != nil {
		return nil, fmt.Errorf("member.MemberDetails.2: %s", err.Error())
	}

	if len(contactDetails) > 0 {
		member.ContactDetails = &contactDetails[0]
	}

	n := NewMemberNominee(m.db, &memberId)

	nominee, err := n.FilterBy(filter)
	if err != nil {
		return nil, fmt.Errorf("member.MemberDetails.3: %s", err.Error())
	}

	if len(nominee) > 0 {
		member.Nominee = &nominee[0]
	}

	o := NewMemberOccupation(m.db, &memberId)

	occupation, err := o.FilterBy(filter)
	if err != nil {
		return nil, fmt.Errorf("member.MemberDetails.4: %s", err.Error())
	}

	if len(occupation) > 0 {
		member.OccupationDetails = &occupation[0]
	}

	b := NewMemberBeneficiary(m.db, &memberId)

	beneficiaries, err := b.FilterBy(filter)
	if err != nil {
		return nil, fmt.Errorf("member.MemberDetails.5: %s", err.Error())
	}

	if len(beneficiaries) > 0 {
		member.Beneficiaries = beneficiaries
	}

	payload, err := json.Marshal(member)
	if err != nil {
		return nil, fmt.Errorf("member.MemberDetails.6: %s", err.Error())
	}

	err = json.Unmarshal(payload, &fullRecord)
	if err != nil {
		return nil, fmt.Errorf("member.MemberDetails.7: %s", err.Error())
	}

	return fullRecord, nil
}

func (m *Member) AddMember(data map[string]any) (int64, error) {
	var id int64

	if data["memberIdNumber"] == nil {
		memberIdNumber := strings.ToUpper(
			regexp.MustCompile(`[^A-Za-z0-9]`).
				ReplaceAllLiteralString(uuid.NewString(), ""),
		)

		data["memberIdNumber"] = memberIdNumber
		data["shortMemberId"] = memberIdNumber[:8]
	}

	data["dateJoined"] = time.Now().Format("2006-01-02")

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
			phoneNumber,
			memberIdNumber,
			shortMemberId,
			dateJoined
		) VALUES (
		 	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)`,
		m.FirstName, m.LastName, m.OtherName,
		m.Gender, m.Title, m.MaritalStatus,
		m.DateOfBirth, m.NationalId, m.UtilityBillType,
		m.UtilityBillNumber, m.FileNumber, m.OldFileNumber,
		m.PhoneNumber, m.MemberIdNumber,
		m.ShortMemberId, m.DateJoined,
	)
	if err != nil {
		return 0, fmt.Errorf("member.AddMember.1: %s", err.Error())
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
		if slices.Contains(m.validFields, key) {
			fields = append(fields, fmt.Sprintf("%s = ?", key))
			values = append(values, value)
		}
	}

	values = append(values, id)

	statement := fmt.Sprintf("UPDATE member SET %s WHERE id=?", strings.Join(fields, ", "))

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

func (m *Member) loadRow(row any) (*Member, bool, error) {
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
		phoneNumber,
		memberIdNumber,
		shortMemberId,
		dateJoined any
	var err error

	val, ok := row.(*sql.Row)
	if ok {
		err = val.Scan(
			&id,
			&firstName,
			&lastName,
			&otherName,
			&gender,
			&title,
			&maritalStatus,
			&dateOfBirth,
			&nationalId,
			&utilityBillType,
			&utilityBillNumber,
			&fileNumber,
			&oldFileNumber,
			&phoneNumber,
			&memberIdNumber,
			&shortMemberId,
			&dateJoined,
		)
	} else {
		val, ok := row.(*sql.Rows)
		if ok {
			err = val.Scan(
				&id,
				&firstName,
				&lastName,
				&otherName,
				&gender,
				&title,
				&maritalStatus,
				&dateOfBirth,
				&nationalId,
				&utilityBillType,
				&utilityBillNumber,
				&fileNumber,
				&oldFileNumber,
				&phoneNumber,
				&memberIdNumber,
				&shortMemberId,
				&dateJoined,
			)
		}
	}
	if err != nil {
		return nil, false, fmt.Errorf("member.loadRow.1: %s", err.Error())
	}

	record := Member{
		ID: id,
	}

	atLeastOneFieldAdded := false

	if firstName != nil {
		value := fmt.Sprintf("%v", firstName)
		if value != "" {
			atLeastOneFieldAdded = true
			record.FirstName = value
		}
	}
	if lastName != nil {
		value := fmt.Sprintf("%v", lastName)
		if value != "" {
			atLeastOneFieldAdded = true
			record.LastName = value
		}
	}
	if otherName != nil {
		value := fmt.Sprintf("%v", otherName)
		if value != "" {
			atLeastOneFieldAdded = true
			record.OtherName = value
		}
	}
	if gender != nil {
		value := fmt.Sprintf("%v", gender)
		if value != "" {
			atLeastOneFieldAdded = true
			record.Gender = value
		}
	}
	if title != nil {
		value := fmt.Sprintf("%v", title)
		if value != "" {
			atLeastOneFieldAdded = true
			record.Title = value
		}
	}
	if maritalStatus != nil {
		value := fmt.Sprintf("%v", maritalStatus)
		if value != "" {
			atLeastOneFieldAdded = true
			record.MaritalStatus = value
		}
	}
	if dateOfBirth != nil {
		value := fmt.Sprintf("%v", dateOfBirth)
		if value != "" {
			atLeastOneFieldAdded = true
			record.DateOfBirth = value
		}
	}
	if nationalId != nil {
		value := fmt.Sprintf("%v", nationalId)
		if value != "" {
			atLeastOneFieldAdded = true
			record.NationalId = value
		}
	}
	if utilityBillType != nil {
		value := fmt.Sprintf("%v", utilityBillType)
		if value != "" {
			atLeastOneFieldAdded = true
			record.UtilityBillType = value
		}
	}
	if utilityBillNumber != nil {
		value := fmt.Sprintf("%v", utilityBillNumber)
		if value != "" {
			atLeastOneFieldAdded = true
			record.UtilityBillNumber = value
		}
	}
	if fileNumber != nil {
		value := fmt.Sprintf("%v", fileNumber)
		if value != "" {
			atLeastOneFieldAdded = true
			record.FileNumber = value
		}
	}
	if oldFileNumber != nil {
		value := fmt.Sprintf("%v", oldFileNumber)
		if value != "" {
			atLeastOneFieldAdded = true
			record.OldFileNumber = value
		}
	}
	if phoneNumber != nil {
		value := fmt.Sprintf("%v", phoneNumber)
		if value != "" {
			atLeastOneFieldAdded = true
			record.PhoneNumber = value
		}
	}
	if memberIdNumber != nil {
		value := fmt.Sprintf("%v", memberIdNumber)
		if value != "" {
			atLeastOneFieldAdded = true
			record.MemberIdNumber = value
		}
	}
	if shortMemberId != nil {
		value := fmt.Sprintf("%v", shortMemberId)
		if value != "" {
			atLeastOneFieldAdded = true
			record.ShortMemberId = value
		}
	}
	if dateJoined != nil {
		value := fmt.Sprintf("%v", dateJoined)
		if value != "" {
			atLeastOneFieldAdded = true
			record.DateJoined = value
		}
	}

	return &record, atLeastOneFieldAdded, nil
}

func (m *Member) FetchMember(id int64) (*Member, error) {

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
		phoneNumber,
		memberIdNumber,
		shortMemberId,
		dateJoined
	FROM member WHERE id=? AND active=1`, id)

	record, found, err := m.loadRow(row)
	if err != nil {
		return nil, fmt.Errorf("member.FetchMember.1: %s", err.Error())
	}

	if !found {
		return nil, nil
	}

	return record, nil
}

func (m *Member) FilterBy(whereStatement string) ([]Member, error) {
	results := []Member{}

	if !regexp.MustCompile("active").MatchString(whereStatement) {
		whereStatement = fmt.Sprintf("%s AND active=1", whereStatement)
	}

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
				phoneNumber,
				memberIdNumber,
				shortMemberId,
				dateJoined
			FROM member %s`,
			whereStatement,
		))
	if err != nil {
		return nil, fmt.Errorf("member.FilterBy.1: %s", err.Error())
	}

	for rows.Next() {
		record, found, err := m.loadRow(rows)
		if err != nil {
			return nil, fmt.Errorf("member.FilterBy.1: %s", err.Error())
		}

		if !found {
			return nil, nil
		}

		results = append(results, *record)
	}

	return results, nil
}

func (m *Member) FetchMemberByPhoneNumber(phoneNumber string) (*Member, error) {
	retries := 0

RETRY:
	time.Sleep(time.Duration(retries) * time.Second)

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
		phoneNumber,
		memberIdNumber,
		shortMemberId,
		dateJoined
	FROM member WHERE phoneNumber=? AND active=1`, phoneNumber)

	record, found, err := m.loadRow(row)
	if err != nil {
		if regexp.MustCompile("SQL logic error: no such table").MatchString(err.Error()) {
			if retries < 3 {
				retries++

				log.Printf("member.FetchMemberByPhoneNumber.retry: %d\n", retries)

				goto RETRY
			}
		}
		return nil, fmt.Errorf("member.FetchMemberByPhoneNumber.1: %s", err.Error())
	}

	if !found {
		return nil, nil
	}

	return record, nil
}
