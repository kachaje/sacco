package server

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sacco/parser"
	"strconv"
)

func CacheFile(filename string, data any) {
	payload, err := json.MarshalIndent(data, "", "  ")
	if err == nil {
		err = os.WriteFile(filename, payload, 0644)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		log.Println(err)
	}
}

func SaveData(
	data any, model, phoneNumber, sessionId, cacheFolder, preferenceFolder *string,
	saveFunc func(
		a map[string]any,
		b map[string]any,
		c map[string]any,
		d map[string]any,
		e []map[string]any,
		f *int64,
	) (*int64, error), sessions map[string]*parser.Session, refData map[string]any) error {
	var sessionFolder string

	if cacheFolder != nil {
		sessionFolder = filepath.Join(*cacheFolder, *phoneNumber)

		_, err := os.Stat(sessionFolder)
		if os.IsNotExist(err) {
			os.MkdirAll(sessionFolder, 0755)
		}
	}

	switch *model {
	case "preferredLanguage":
		val, ok := data.(map[string]any)
		if ok {
			if val["language"] != nil && phoneNumber != nil {
				language, ok := val["language"].(string)
				if ok {
					SavePreference(*phoneNumber, "language", language, *preferenceFolder)
				}
			}
		}

	case "memberDetails":
		memberData, ok := data.(map[string]any)
		if ok {
			var id int64
			var err error

			if phoneNumber != nil && *phoneNumber != "default" {
				memberData["defaultPhoneNumber"] = *phoneNumber
			}

			filename := filepath.Join(sessionFolder, "memberDetails.json")

			transactionDone := false

			// By default cache the data first in case we lose database connection
			CacheFile(filename, memberData)
			defer func() {
				if transactionDone {
					os.Remove(filename)
				}
			}()

			if sessions[*sessionId].ContactsAdded ||
				sessions[*sessionId].BeneficiariesAdded ||
				sessions[*sessionId].NomineeAdded ||
				sessions[*sessionId].OccupationAdded {

				var contactsData, nomineeData, occupationData map[string]any
				var beneficiariesData []map[string]any

				contactsFile := filepath.Join(sessionFolder, "contactDetails.json")

				_, err = os.Stat(contactsFile)
				if !os.IsNotExist(err) {
					content, err := os.ReadFile(contactsFile)
					if err != nil {
						log.Printf("server.SaveData.memberDetails.1:%s\n", err.Error())
					} else {
						err = json.Unmarshal(content, &contactsData)
						if err != nil {
							log.Printf("server.SaveData.memberDetails.2:%s\n", err.Error())
						}
					}
				}

				nomineeFile := filepath.Join(sessionFolder, "nomineeDetails.json")

				_, err = os.Stat(nomineeFile)
				if !os.IsNotExist(err) {
					content, err := os.ReadFile(nomineeFile)
					if err != nil {
						log.Printf("server.SaveData.memberDetails.3:%s\n", err.Error())
					} else {
						err = json.Unmarshal(content, &nomineeData)
						if err != nil {
							log.Printf("server.SaveData.memberDetails.4:%s\n", err.Error())
						}
					}
				}

				occupationFile := filepath.Join(sessionFolder, "occupationDetails.json")

				_, err = os.Stat(occupationFile)
				if !os.IsNotExist(err) {
					content, err := os.ReadFile(occupationFile)
					if err != nil {
						log.Printf("server.SaveData.memberDetails.5:%s\n", err.Error())
					} else {
						err = json.Unmarshal(content, &occupationData)
						if err != nil {
							log.Printf("server.SaveData.memberDetails.6:%s\n", err.Error())
						}
					}
				}

				beneficiariesFile := filepath.Join(sessionFolder, "beneficiaries.json")

				_, err = os.Stat(beneficiariesFile)
				if !os.IsNotExist(err) {
					content, err := os.ReadFile(beneficiariesFile)
					if err != nil {
						log.Printf("server.SaveData.memberDetails.7:%s\n", err.Error())
					} else {
						err = json.Unmarshal(content, &beneficiariesData)
						if err != nil {
							log.Printf("server.SaveData.memberDetails.8:%s\n", err.Error())
						}
					}
				}

				if saveFunc == nil {
					return fmt.Errorf("server.SaveData.memberDetails.9:missing saveFunc")
				}

				mid, err := saveFunc(memberData, contactsData, nomineeData, occupationData, beneficiariesData, nil)
				if err != nil {
					log.Println(err)
					return err
				}

				id = *mid

				if len(contactsData) > 0 {
					contactsData["memberId"] = id

					memberData["contactDetails"] = contactsData

					if os.Getenv("DEBUG") == "true" {
						CacheFile(contactsFile, contactsData)
					} else {
						os.Remove(contactsFile)
					}
				}

				if len(nomineeData) > 0 {
					nomineeData["memberId"] = id

					memberData["nomineeDetails"] = nomineeData

					if os.Getenv("DEBUG") == "true" {
						CacheFile(nomineeFile, nomineeData)
					} else {
						os.Remove(nomineeFile)
					}
				}

				if len(occupationData) > 0 {
					occupationData["memberId"] = id

					memberData["occupationDetails"] = occupationData

					if os.Getenv("DEBUG") == "true" {
						CacheFile(occupationFile, occupationData)
					} else {
						os.Remove(occupationFile)
					}
				}

				if len(beneficiariesData) > 0 {
					for i := range beneficiariesData {
						beneficiariesData[i]["memberId"] = id
					}

					memberData["beneficiaries"] = beneficiariesData

					if os.Getenv("DEBUG") == "true" {
						CacheFile(beneficiariesFile, beneficiariesData)
					} else {
						os.Remove(beneficiariesFile)
					}
				}

				transactionDone = true
			} else {
				if saveFunc == nil {
					return fmt.Errorf("server.SaveData.memberDetails.10:missing saveFunc")
				}

				mid, err := saveFunc(memberData, nil, nil, nil, nil, nil)
				if err != nil {
					log.Println(err)
					return fmt.Errorf("server.SaveData.memberDetails.11:%s", err.Error())
				}

				id = *mid

				transactionDone = true
			}

			sessions[*sessionId].MemberId = &id

			memberData["id"] = id

			sessions[*sessionId].ActiveMemberData = memberData

			sessions[*sessionId].RefreshSession()

			sessions[*sessionId].LoadMemberCache(*phoneNumber, *cacheFolder)
		}

	case "contactDetails":
		val, ok := data.(map[string]any)
		if ok {
			filename := filepath.Join(sessionFolder, "contactDetails.json")

			transactionDone := false

			// By default cache the data first in case we lose database connection
			CacheFile(filename, val)
			defer func() {
				if transactionDone {
					os.Remove(filename)
				}
			}()

			if sessions[*sessionId].MemberId != nil {
				val["memberId"] = *sessions[*sessionId].MemberId

				if saveFunc == nil {
					log.Println("Missing saveFunc")
					return fmt.Errorf("server.SaveData.contactDetails.1:missing saveFunc")
				}

				_, err := saveFunc(nil, val, nil, nil, nil, sessions[*sessionId].MemberId)
				if err != nil {
					return fmt.Errorf("server.SaveData.contactDetails.2:%s", err.Error())
				}

				transactionDone = true
			}

			sessions[*sessionId].ActiveMemberData["contactDetails"] = val

			sessions[*sessionId].ContactsAdded = true

			sessions[*sessionId].RefreshSession()

			sessions[*sessionId].LoadMemberCache(*phoneNumber, *cacheFolder)
		}

	case "nomineeDetails":
		val, ok := data.(map[string]any)
		if ok {
			filename := filepath.Join(sessionFolder, "nomineeDetails.json")

			transactionDone := false

			// By default cache the data first in case we lose database connection
			CacheFile(filename, val)
			defer func() {
				if transactionDone {
					os.Remove(filename)
				}
			}()

			if sessions[*sessionId].MemberId != nil {
				val["memberId"] = *sessions[*sessionId].MemberId

				if saveFunc == nil {
					log.Println("Missing saveFunc")
					return fmt.Errorf("server.SaveData.nomineeDetails.1:missing saveFunc")
				}

				_, err := saveFunc(nil, nil, val, nil, nil, sessions[*sessionId].MemberId)
				if err != nil {
					return fmt.Errorf("server.SaveData.nomineeDetails.2:%s", err.Error())
				}

				transactionDone = true
			}

			sessions[*sessionId].ActiveMemberData["nomineeDetails"] = val

			sessions[*sessionId].NomineeAdded = true

			sessions[*sessionId].LoadMemberCache(*phoneNumber, *cacheFolder)
		}

	case "occupationDetails":
		val, ok := data.(map[string]any)
		if ok {
			filename := filepath.Join(sessionFolder, "occupationDetails.json")

			transactionDone := false

			// By default cache the data first in case we lose database connection
			CacheFile(filename, val)
			defer func() {
				if transactionDone {
					os.Remove(filename)
				}
			}()

			for _, key := range []string{"netPay", "grossPay", "periodEmployed"} {
				if val[key] != nil {
					nv, ok := val[key].(string)
					if ok {
						real, err := strconv.ParseFloat(nv, 64)
						if err == nil {
							val[key] = real
						}
					}
				}
			}

			if sessions[*sessionId].MemberId != nil {
				val["memberId"] = *sessions[*sessionId].MemberId

				if saveFunc == nil {
					return fmt.Errorf("server.SaveData.occupationDetails.1:missing saveFunc")
				}

				_, err := saveFunc(nil, nil, nil, val, nil, sessions[*sessionId].MemberId)
				if err != nil {
					return fmt.Errorf("server.SaveData.occupationDetails.2:%s", err.Error())
				}

				transactionDone = true
			}

			sessions[*sessionId].ActiveMemberData["occupationDetails"] = val

			sessions[*sessionId].OccupationAdded = true

			sessions[*sessionId].RefreshSession()

			sessions[*sessionId].LoadMemberCache(*phoneNumber, *cacheFolder)
		}

	case "beneficiaries":
		rawData, ok := data.(map[string]any)
		if ok {
			var memberId int64

			records := []map[string]any{}

			for i := range 4 {
				var name, contact string
				var percentage float64
				var id int64

				index := i + 1

				nameLabel := fmt.Sprintf("name%v", index)
				percentLabel := fmt.Sprintf("percentage%v", index)
				contactLabel := fmt.Sprintf("contact%v", index)
				idLabel := fmt.Sprintf("id%v", index)
				memberIdLabel := fmt.Sprintf("memberId%v", index)

				var row map[string]any

				if refData != nil && refData[idLabel] != nil {
					v, err := strconv.ParseInt(fmt.Sprintf("%v", refData[idLabel]), 10, 64)
					if err == nil {
						id = v
					} else {
						log.Printf("server.SaveData.beneficiaries.2:%s", err.Error())
					}
				}

				if refData != nil && refData[memberIdLabel] != nil {
					v, err := strconv.ParseInt(fmt.Sprintf("%v", refData[memberIdLabel]), 10, 64)
					if err == nil {
						memberId = v
					} else {
						log.Printf("server.SaveData.beneficiaries.3:%s", err.Error())
					}
				}

				if rawData[nameLabel] == nil {
					if id != 0 {
						row = map[string]any{
							"active": 0,
						}
					} else {
						continue
					}
				} else if rawData[nameLabel] != nil {
					name = fmt.Sprintf("%v", rawData[nameLabel])
					contact = fmt.Sprintf("%v", rawData[contactLabel])

					v, err := strconv.ParseFloat(fmt.Sprintf("%v", rawData[percentLabel]), 64)
					if err == nil {
						percentage = v
					} else {
						log.Printf("server.SaveData.beneficiaries.1:%s", err.Error())
					}

					row = map[string]any{
						"name":       name,
						"percentage": percentage,
						"contact":    contact,
					}
				}

				if id != 0 {
					row["id"] = id
				}

				if sessions != nil && sessions[*sessionId].MemberId != nil {
					row["memberId"] = *sessions[*sessionId].MemberId

					memberId = *sessions[*sessionId].MemberId
				} else if memberId != 0 {
					row["memberId"] = memberId
				}

				records = append(records, row)
			}

			filename := filepath.Join(sessionFolder, "beneficiaries.json")

			transactionDone := false

			// By default cache the data first in case we lose database connection
			CacheFile(filename, records)
			defer func() {
				if transactionDone {
					os.Remove(filename)
				}
			}()

			if os.Getenv("DEBUG") == "true" {
				payload, _ := json.MarshalIndent(records, "", "  ")

				fmt.Println(string(payload))
			}

			if memberId != 0 {
				if saveFunc == nil {
					return fmt.Errorf("server.SaveData.beneficiaries.4:missing saveFunc")
				}

				_, err := saveFunc(nil, nil, nil, nil, records, &memberId)
				if err != nil {
					return fmt.Errorf("server.SaveData.beneficiaries.5:%s", err.Error())
				}

				transactionDone = true
			}

			if phoneNumber != nil && cacheFolder != nil && sessions != nil && sessionId != nil {
				sessions[*sessionId].ActiveMemberData["beneficiaries"] = records

				sessions[*sessionId].BeneficiariesAdded = true

				sessions[*sessionId].RefreshSession()

				sessions[*sessionId].LoadMemberCache(*phoneNumber, *cacheFolder)
			}
		}

	default:
		fmt.Println("##########", *phoneNumber, *sessionId, data)
	}

	return nil
}

func SavePreference(phoneNumber, key, value, preferencesFolder string) error {
	settingsFile := filepath.Join(preferencesFolder, phoneNumber)

	data := map[string]any{}

	_, err := os.Stat(settingsFile)
	if !os.IsNotExist(err) {
		content, err := os.ReadFile(settingsFile)
		if err != nil {
			log.Println(err)
			return err
		}

		err = json.Unmarshal(content, &data)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	data[key] = value

	payload, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Println(err)
		return err
	}

	return os.WriteFile(settingsFile, payload, 0644)
}
