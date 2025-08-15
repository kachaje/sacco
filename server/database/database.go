package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

type Database struct {
	DbName string
	DB     *sql.DB
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
			gender TEXT,
			title TEXT,
			maritalStatus TEXT,
			dateOfBirth TEXT,
			nationalId TEXT,
			utilityBillType TEXT,
			utilityBillNumber TEXT,
			fileNumber TEXT,
			oldFileNumber TEXT
		);
		CREATE TABLE IF NOT EXISTS memberContact (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			postalAddress TEXT,
			residentialAddress TEXT,
			phoneNumber TEXT,
			homeVillage TEXT,
			homeTA TEXT,
			homeDistrict TEXT
		);
		CREATE TABLE IF NOT EXISTS memberNominee (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			nextOfKinName TEXT,
			nextOfKinPhone TEXT,
			nextOfKinAddress TEXT
		);
		CREATE TABLE IF NOT EXISTS memberOccupation (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			employerName TEXT,
			netPay REAL,
			jobTitle TEXT,
			employerAddress TEXT,
			highestQualification TEXT
		);
		CREATE TABLE IF NOT EXISTS memberBeneficiary (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			memberId INTEGER NOT NULL,
			name TEXT,
			percentage REAL,
			contact TEXT
		);
	`
	_, err := d.DB.Exec(sqlStmt)
	if err != nil {
		return err
	}
	return nil
}
