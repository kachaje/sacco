package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type MemberOccupation struct {
	ID                   int64   `json:"id"`
	MemberId             int64   `json:"memberId"`
	EmployerName         string  `json:"employerName"`
	EmployerAddress      string  `json:"employerAddress"`
	EmployerPhone        string  `json:"employerPhone"`
	JobTitle             string  `json:"jobTitle"`
	PeriodEmployed       float64 `json:"periodEmployed"`
	GrossPay             float64 `json:"grossPay"`
	NetPay               float64 `json:"netPay"`
	HighestQualification string  `json:"highestQualification"`

	db *sql.DB
}

func NewMemberOccupation(db *sql.DB, memberId *int64) *MemberOccupation {
	m := &MemberOccupation{
		db: db,
	}

	if memberId != nil {
		m.MemberId = *memberId
	}

	return m
}

func (m *MemberOccupation) AddMemberOccupation(data map[string]any) (int64, error) {
	var id int64

	payload, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("memberOccupation.AddMemberOccupation.1: %s", err.Error())
	}

	err = json.Unmarshal(payload, m)
	if err != nil {
		return 0, fmt.Errorf("memberOccupation.AddMemberOccupation.2: %s", err.Error())
	}

	result, err := m.db.ExecContext(
		context.Background(),
		`INSERT INTO memberOccupation (
			memberId,
			employerName,
			employerAddress,
			employerPhone,
			jobTitle,
			periodEmployed,
			grossPay,
			netPay,
			highestQualification
		) VALUES (
		 	?, ?, ?, ?, ?, ?, ?, ?, ?
		)`,
		m.MemberId, m.EmployerName, m.EmployerAddress,
		m.EmployerPhone, m.JobTitle, m.PeriodEmployed,
		m.GrossPay, m.NetPay, m.HighestQualification,
	)
	if err != nil {
		return 0, fmt.Errorf("memberOccupation.AddMemberOccupation.3: %s", err.Error())
	}

	if id, err = result.LastInsertId(); err != nil {
		return 0, fmt.Errorf("memberOccupation.AddMemberOccupation.4: %s", err.Error())
	}

	return id, nil
}

func (m *MemberOccupation) UpdateMemberOccupation(data map[string]any, id int64) error {
	fields := []string{}
	values := []any{}

	for key, value := range data {
		fields = append(fields, fmt.Sprintf("%s = ?", key))
		values = append(values, value)
	}

	values = append(values, id)

	statement := fmt.Sprintf("UPDATE memberOccupation SET %s WHERE id=?", strings.Join(fields, ", "))

	_, err := m.db.Exec(statement, values...)
	if err != nil {
		return fmt.Errorf("memberOccupation.UpdateMemberOccupation.1: %s", err.Error())
	}

	return nil
}

func (m *MemberOccupation) loadRow(row any) (*MemberOccupation, bool, error) {
	var id int64
	var memberId int64
	var employerName,
		employerAddress,
		employerPhone,
		jobTitle,
		periodEmployed,
		grossPay,
		netPay,
		highestQualification any
	var err error

	val, ok := row.(*sql.Row)
	if ok {
		err = val.Scan(
			&id,
			&memberId,
			&employerName,
			&employerAddress,
			&employerPhone,
			&jobTitle,
			&periodEmployed,
			&grossPay,
			&netPay,
			&highestQualification,
		)
	} else {
		val, ok := row.(*sql.Rows)
		if ok {
			err = val.Scan(
				&id,
				&memberId,
				&employerName,
				&employerAddress,
				&employerPhone,
				&jobTitle,
				&periodEmployed,
				&grossPay,
				&netPay,
				&highestQualification,
			)
		}
	}
	if err != nil {
		return nil, false, fmt.Errorf("memberOccupation.loadRow.1: %s", err.Error())
	}

	memberOccupation := MemberOccupation{
		ID:       id,
		MemberId: memberId,
	}

	atLeastOneFieldAdded := false

	if employerName != nil {
		value := fmt.Sprintf("%v", employerName)
		if value != "" {
			atLeastOneFieldAdded = true
			memberOccupation.EmployerName = fmt.Sprintf("%v", employerName)
		}
	}
	if employerAddress != nil {
		value := fmt.Sprintf("%v", employerAddress)
		if value != "" {
			atLeastOneFieldAdded = true
			memberOccupation.EmployerAddress = fmt.Sprintf("%v", employerAddress)
		}
	}
	if employerPhone != nil {
		value := fmt.Sprintf("%v", employerPhone)
		if value != "" {
			atLeastOneFieldAdded = true
			memberOccupation.EmployerPhone = fmt.Sprintf("%v", employerPhone)
		}
	}
	if jobTitle != nil {
		value := fmt.Sprintf("%v", jobTitle)
		if value != "" {
			atLeastOneFieldAdded = true
			memberOccupation.JobTitle = fmt.Sprintf("%v", jobTitle)
		}
	}
	if periodEmployed != nil {
		value := periodEmployed.(float64)
		if value != 0 {
			atLeastOneFieldAdded = true
			memberOccupation.PeriodEmployed = value
		}
	}
	if grossPay != nil {
		value := grossPay.(float64)
		if value != 0 {
			atLeastOneFieldAdded = true
			memberOccupation.GrossPay = value
		}
	}
	if netPay != nil {
		value := netPay.(float64)
		if value != 0 {
			atLeastOneFieldAdded = true
			memberOccupation.NetPay = value
		}
	}
	if highestQualification != nil {
		value := fmt.Sprintf("%v", highestQualification)
		if value != "" {
			atLeastOneFieldAdded = true
			memberOccupation.HighestQualification = fmt.Sprintf("%v", highestQualification)
		}
	}

	return &memberOccupation, atLeastOneFieldAdded, nil
}

func (m *MemberOccupation) FetchMemberOccupation(id int64) (*MemberOccupation, error) {
	row := m.db.QueryRow(`SELECT 
		id,
		memberId,
		employerName,
		employerAddress,
		employerPhone,
		jobTitle,
		periodEmployed,
		grossPay,
		netPay,
		highestQualification
	FROM memberOccupation WHERE id=?`, id)

	memberOccupation, atLeastOneFieldAdded, err := m.loadRow(row)
	if err != nil {
		return nil, fmt.Errorf("memberOccupation.FetchMemberOccupation.1: %s", err.Error())
	}

	if !atLeastOneFieldAdded {
		return nil, nil
	}

	return memberOccupation, nil
}

func (m *MemberOccupation) FilterBy(whereStatement string) ([]MemberOccupation, error) {
	results := []MemberOccupation{}

	rows, err := m.db.QueryContext(
		context.Background(),
		fmt.Sprintf(`SELECT
			id,
			memberId,
			employerName,
			employerAddress,
			employerPhone,
			jobTitle,
			periodEmployed,
			grossPay,
			netPay,
			highestQualification
		FROM memberOccupation %s`, whereStatement))
	if err != nil {
		return nil, fmt.Errorf("memberOccupation.FetchBy.1: %s", err.Error())
	}

	for rows.Next() {
		memberOccupation, atLeastOneFieldAdded, err := m.loadRow(rows)
		if err != nil {
			return nil, fmt.Errorf("memberOccupation.FetchBy.1: %s", err.Error())
		}

		if !atLeastOneFieldAdded {
			continue
		}

		results = append(results, *memberOccupation)
	}

	return results, nil
}
