package database

import (
	"database/sql"
	"fmt"
	"log"
	"sacco/server/database/models"
	"strings"

	_ "modernc.org/sqlite"
)

type Database struct {
	DbName            string
	DB                *sql.DB
	Member            *models.Member
	MemberContact     *models.MemberContact
	MemberBeneficiary *models.MemberBeneficiary
	MemberOccupation  *models.MemberOccupation
	MemberNominee     *models.MemberNominee
}

func NewDatabase(dbname string) *Database {
	if dbname != ":memory:" && !strings.HasSuffix(dbname, ".db") {
		dbname = fmt.Sprintf("%s.db", dbname)
	}

	db, err := sql.Open("sqlite", dbname)
	if err != nil {
		log.Fatal(err)
	}

	instance := &Database{
		DbName: dbname,
		DB:     db,
	}

	err = instance.initDb()
	if err != nil {
		log.Fatal(err)
	}

	instance.Member = models.NewMember(db)
	instance.MemberContact = models.NewMemberContact(db, nil)
	instance.MemberBeneficiary = models.NewMemberBeneficiary(db, nil)
	instance.MemberOccupation = models.NewMemberOccupation(db, nil)
	instance.MemberNominee = models.NewMemberNominee(db, nil)

	return instance
}

func (d *Database) Close() {
	d.DB.Close()
}

func (d *Database) initDb() error {
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS member (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			firstName TEXT,
			lastName TEXT,
			otherName TEXT,
			gender TEXT CHECK (gender IN ('Male', 'Female')),
			title TEXT,
			maritalStatus TEXT,
			dateOfBirth TEXT,
			nationalId TEXT,
			utilityBillType TEXT,
			utilityBillNumber TEXT,
			fileNumber TEXT,
			oldFileNumber TEXT,
			defaultPhoneNumber TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS memberContact (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			postalAddress TEXT,
			residentialAddress TEXT,
			phoneNumber TEXT,
			homeVillage TEXT,
			homeTA TEXT,
			homeDistrict TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS memberNominee (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			nomineeName TEXT,
			nomineePhone TEXT,
			nomineeAddress TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS memberOccupation (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			employerName TEXT,
			employerAddress TEXT,
			employerPhone TEXT,
			jobTitle TEXT,
			periodEmployed REAL,
			grossPay REAL,
			netPay REAL,
			highestQualification TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS memberBeneficiary (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			name TEXT,
			percentage REAL,
			contact TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS memberBusiness (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			numberOfBusinessYears REAL,
			typeOfBusiness TEXT,
			nameOfBusiness TEXT,
			tradingArea TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS memberLastYearBusinessHistory (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberBusinessId INTEGER NOT NULL,
			totalIncome REAL,
			totalCostOfGoods REAL,
			employeesWages REAL,
			ownSalary REAL,
			transport REAL,
			loanInterest REAL,
			utilities REAL,
			rentals REAL,
			otherCosts REAL,
			totalCosts REAL,
			netProfitLoss REAL,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS memberNextYearBusinessProjection (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberBusinessId INTEGER NOT NULL,
			totalIncome REAL,
			totalCostOfGoods REAL,
			employeesWages REAL,
			ownSalary REAL,
			transport REAL,
			loanInterest REAL,
			utilities REAL,
			rentals REAL,
			otherCosts REAL,
			totalCosts REAL,
			netProfitLoss REAL,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS share (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			numberOfShares REAL,
			pricePerShare REAL,
			sharesType TEXT NOT NULL CHECK (sharesType IN ('FIXED', 'REDEEMABLE')),
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS loan (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			employeeNumber TEXT,
			loanAmount REAL,
			repaymentPeriod REAL,
			loanPurpose TEXT,
			loanStatus TEXT NOT NULL CHECK (loanStatus IN ('PENDING', 'APPROVED', 'REJECTED')),
			loanType TEXT NOT NULL CHECK (loanType IN ('PERSONAL', 'BUSINESS', 'AGRICULTURAL', 'EMERGENCY')),
			amountRecommended REAL,
			approvedBy TEXT,
			approvalDate TEXT,
			amountApproved REAL,
			verifiedBy TEXT,
			dateVerified TEXT,
			denialOrPartialReason TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS loanLiability (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			loanId INTEGER NOT NULL,
			description TEXT,
			value REAL,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS loanSecurity (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			loanId INTEGER NOT NULL,
			description TEXT,
			value REAL,
			serialNumber TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS loanWitness (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			loanId INTEGER NOT NULL,
			name TEXT,
			telephone TEXT,
			address TEXT,
			date TEXT,
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := d.DB.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) AddMember(
	memberData, contactData,
	nomineeData, occupationData map[string]any,
	beneficiariesData []map[string]any,
	existingMemberId *int64,
) (*int64, error) {
	var memberId int64
	var err error

	if existingMemberId == nil {
		if memberData != nil {
			id, err := d.Member.AddMember(memberData)
			if err != nil {
				return nil, err
			}

			memberId = id
		}
	} else {
		memberId = *existingMemberId
	}

	if contactData != nil {
		contactData["memberId"] = memberId

		_, err = d.MemberContact.AddMemberContact(contactData)
		if err != nil {
			return nil, err
		}
	}

	if nomineeData != nil {
		nomineeData["memberId"] = memberId

		_, err = d.MemberNominee.AddMemberNominee(nomineeData)
		if err != nil {
			return nil, err
		}
	}

	if occupationData != nil {
		occupationData["memberId"] = memberId

		_, err = d.MemberOccupation.AddMemberOccupation(occupationData)
		if err != nil {
			return nil, err
		}
	}

	for _, beneficiaryData := range beneficiariesData {
		beneficiaryData["memberId"] = memberId

		_, err = d.MemberBeneficiary.AddMemberBeneficiary(beneficiaryData)
		if err != nil {
			return nil, err
		}
	}

	return &memberId, nil
}

func (d *Database) MemberByDefaultPhoneNumber(phoneNumber string) (map[string]any, error) {
	member, err := d.Member.FetchMemberByPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}

	fullRecord, err := d.Member.MemberDetails(member.ID)
	if err != nil {
		return nil, err
	}

	return fullRecord, nil
}
