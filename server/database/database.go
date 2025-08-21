package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"sacco/server/database/models"
	"sacco/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaStatement string

//go:embed models.yml
var modelTemplates string

var modelTemplatesData map[string]any

type Database struct {
	DbName            string
	DB                *sql.DB
	Member            *models.Member
	MemberContact     *models.MemberContact
	MemberBeneficiary *models.MemberBeneficiary
	MemberOccupation  *models.MemberOccupation
	MemberNominee     *models.MemberNominee
	GenericModels     map[string]*models.Model

	Mu *sync.Mutex
}

func init() {
	var err error

	modelTemplatesData, err = utils.LoadYaml(modelTemplates)
	if err != nil {
		log.Fatal(err)
	}
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
		DbName:        dbname,
		DB:            db,
		GenericModels: map[string]*models.Model{},
		Mu:            &sync.Mutex{},
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

	for table, value := range modelTemplatesData {
		val, ok := value.([]any)
		if ok {
			fields := []string{}

			for _, v := range val {
				fields = append(fields, v.(string))
			}

			model, err := models.NewModel(instance.DB, table, fields)
			if err != nil {
				log.Printf("server.database.NewDatabase: %s", err.Error())
				continue
			}

			instance.GenericModels[table] = model
		}
	}

	return instance
}

func (d *Database) Close() {
	d.DB.Close()
}

func (d *Database) initDb() error {
	_, err := d.DB.Exec(schemaStatement)
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

	if memberData["id"] != nil {
		val := fmt.Sprintf("%v", memberData["id"])

		v, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			existingMemberId = &v
		}
	}

	if existingMemberId == nil {
		if memberData != nil {
			id, err := d.Member.AddMember(memberData)
			if err != nil {
				return nil, fmt.Errorf("database.AddMember.1: %s", err.Error())
			}

			memberId = id
		}
	} else {
		memberId = *existingMemberId

		if memberData != nil {
			err = d.Member.UpdateMember(memberData, *existingMemberId)
			if err != nil {
				return nil, fmt.Errorf("database.AddMember.2: %s", err.Error())
			}
		}
	}

	if contactData != nil {
		contactData["memberId"] = memberId

		if contactData["id"] != nil {
			id, ok := contactData["id"].(int64)
			if ok {
				err = d.MemberContact.UpdateMemberContact(contactData, id)
				if err != nil {
					return nil, fmt.Errorf("database.AddMember.4a: %s", err.Error())
				}
			}
		} else {
			_, err = d.MemberContact.AddMemberContact(contactData)
			if err != nil {
				return nil, fmt.Errorf("database.AddMember.3: %s", err.Error())
			}
		}
	}

	if nomineeData != nil {
		nomineeData["memberId"] = memberId

		if nomineeData["id"] != nil {
			id, ok := nomineeData["id"].(int64)
			if ok {
				err = d.MemberNominee.UpdateMemberNominee(nomineeData, id)
				if err != nil {
					return nil, fmt.Errorf("database.AddMember.4a: %s", err.Error())
				}
			}
		} else {
			_, err = d.MemberNominee.AddMemberNominee(nomineeData)
			if err != nil {
				return nil, fmt.Errorf("database.AddMember.4b: %s", err.Error())
			}
		}
	}

	if occupationData != nil {
		occupationData["memberId"] = memberId

		if occupationData["id"] != nil {
			id, ok := occupationData["id"].(int64)
			if ok {
				err = d.MemberOccupation.UpdateMemberOccupation(occupationData, id)
				if err != nil {
					return nil, fmt.Errorf("database.AddMember.5a: %s", err.Error())
				}
			}
		} else {
			_, err = d.MemberOccupation.AddMemberOccupation(occupationData)
			if err != nil {
				return nil, fmt.Errorf("database.AddMember.5b: %s", err.Error())
			}
		}
	}

	for _, beneficiaryData := range beneficiariesData {
		beneficiaryData["memberId"] = memberId

		if beneficiaryData["id"] != nil {
			id, ok := beneficiaryData["id"].(int64)
			if ok {
				err = d.MemberBeneficiary.UpdateMemberBeneficiary(beneficiaryData, id)
				if err != nil {
					return nil, fmt.Errorf("database.AddMember.6a: %s", err.Error())
				}
			}
		} else {
			_, err = d.MemberBeneficiary.AddMemberBeneficiary(beneficiaryData)
			if err != nil {
				return nil, fmt.Errorf("database.AddMember.6b: %s", err.Error())
			}
		}
	}

	return &memberId, nil
}

func (d *Database) MemberByPhoneNumber(phoneNumber string) (map[string]any, error) {
	member, err := d.Member.FetchMemberByPhoneNumber(phoneNumber, 0)
	if err != nil {
		return nil, fmt.Errorf("database.MemberByPhoneNumber.1: %s", err.Error())
	}

	fullRecord, err := d.Member.MemberDetails(member.ID)
	if err != nil {
		return nil, fmt.Errorf("database.MemberByPhoneNumber.2: %s", err.Error())
	}

	return fullRecord, nil
}

func (d *Database) GenericsSaveData(data map[string]any,
	model string,
	retries int,
) (*int64, error) {
	time.Sleep(time.Duration(retries) * time.Second)

	if !d.Mu.TryLock() {
		if retries < 5 {
			retries++

			return d.GenericsSaveData(data, model, retries)
		}

		return nil, fmt.Errorf("server.database.GenericsSaveData: failed to save due to lock error")
	}
	defer d.Mu.Unlock()

	if d.GenericModels[model] == nil {
		return nil, fmt.Errorf("server.database.GenericsSaveData: model %s does not exist", model)
	}

	id, err := d.GenericModels[model].AddRecord(data)
	if err != nil {
		return nil, err
	}

	return id, nil
}
