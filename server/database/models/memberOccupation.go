package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type MemberOccupation struct {
	ID                   int64
	MemberId             int64
	EmployerName         string
	NetPay               float64
	JobTitle             string
	EmployerAddress      string
	HighestQualification string

	db *sql.DB
}

func NewMemberOccupation(db *sql.DB) *MemberOccupation {
	return &MemberOccupation{
		db: db,
	}
}

func (m *MemberOccupation) AddMemberOccupation(data map[string]any) (int64, error) {
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
		`INSERT INTO memberOccupation (
			memberId,
			employerName,
			netPay,
			jobTitle,
			employerAddress,
			highestQualification
		) VALUES (
		 	?, ?, ?, ?, ?, ?
		)`,
		m.MemberId, m.EmployerName, m.NetPay,
		m.JobTitle, m.EmployerAddress,
		m.HighestQualification,
	)
	if err != nil {
		return 0, err
	}

	if id, err = result.LastInsertId(); err != nil {
		return 0, err
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
		return err
	}

	return nil
}

func (m *MemberOccupation) FetchMemberOccupation(id int64) (*MemberOccupation, error) {

	row := m.db.QueryRow(`SELECT 
		memberId,
		employerName,
		netPay,
		jobTitle,
		employerAddress,
		highestQualification
	FROM memberOccupation WHERE id=?`, id)

	var memberId int64
	var employerName,
		netPay,
		jobTitle,
		employerAddress,
		highestQualification any

	err := row.Scan(
		&memberId,
		&employerName,
		&netPay,
		&jobTitle,
		&employerAddress,
		&highestQualification,
	)
	if err != nil {
		return nil, err
	}

	memberOccupation := &MemberOccupation{
		ID:       id,
		MemberId: memberId,
	}

	if employerName != nil {
		memberOccupation.EmployerName = fmt.Sprintf("%v", employerName)
	}
	if netPay != nil {
		memberOccupation.NetPay = netPay.(float64)
	}
	if jobTitle != nil {
		memberOccupation.JobTitle = fmt.Sprintf("%v", jobTitle)
	}
	if employerAddress != nil {
		memberOccupation.EmployerAddress = fmt.Sprintf("%v", employerAddress)
	}
	if highestQualification != nil {
		memberOccupation.HighestQualification = fmt.Sprintf("%v", highestQualification)
	}

	return memberOccupation, nil
}

func (m *MemberOccupation) FilterBy(whereStatement string) ([]MemberOccupation, error) {
	results := []MemberOccupation{}

	rows, err := m.db.QueryContext(context.Background(), fmt.Sprintf("SELECT * FROM memberOccupation %s", whereStatement))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int64
		var memberId int64
		var employerName,
			netPay,
			jobTitle,
			employerAddress,
			highestQualification any

		err := rows.Scan(
			&id,
			&memberId,
			&employerName,
			&netPay,
			&jobTitle,
			&employerAddress,
			&highestQualification,
		)
		if err != nil {
			return nil, err
		}

		memberOccupation := MemberOccupation{
			ID:       id,
			MemberId: memberId,
		}

		if employerName != nil {
			memberOccupation.EmployerName = fmt.Sprintf("%v", employerName)
		}
		if netPay != nil {
			memberOccupation.NetPay = netPay.(float64)
		}
		if jobTitle != nil {
			memberOccupation.JobTitle = fmt.Sprintf("%v", jobTitle)
		}
		if employerAddress != nil {
			memberOccupation.EmployerAddress = fmt.Sprintf("%v", employerAddress)
		}
		if highestQualification != nil {
			memberOccupation.HighestQualification = fmt.Sprintf("%v", highestQualification)
		}

		results = append(results, memberOccupation)
	}

	return results, nil
}
