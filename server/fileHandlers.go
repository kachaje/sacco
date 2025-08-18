package server

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sacco/server/menus"
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

func SaveData(data any, model, phoneNumber, sessionId, cacheFolder, preferenceFolder *string) {
	sessionFolder := filepath.Join(*cacheFolder, *phoneNumber)

	_, err := os.Stat(sessionFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(sessionFolder, 0755)
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

			if *phoneNumber != "default" {
				memberData["defaultPhoneNumber"] = *phoneNumber
			}

			if menus.Sessions[*sessionId].ContactsAdded ||
				menus.Sessions[*sessionId].BeneficiariesAdded ||
				menus.Sessions[*sessionId].NomineeAdded ||
				menus.Sessions[*sessionId].OccupationAdded {
				var contactsData, nomineeData, occupationData map[string]any
				var beneficiariesData []map[string]any

				contactsFile := filepath.Join(sessionFolder, "contactDetails.json")

				_, err = os.Stat(contactsFile)
				if !os.IsNotExist(err) {
					content, err := os.ReadFile(contactsFile)
					if err != nil {
						log.Println(err)
					} else {
						err = json.Unmarshal(content, &contactsData)
						if err != nil {
							log.Println(err)
						}
					}

					os.Remove(contactsFile)
				}

				nomineeFile := filepath.Join(sessionFolder, "nomineeDetails.json")
				if !os.IsNotExist(err) {
					content, err := os.ReadFile(nomineeFile)
					if err != nil {
						log.Println(err)
					} else {
						err = json.Unmarshal(content, &nomineeData)
						if err != nil {
							log.Println(err)
						}
					}

					os.Remove(nomineeFile)
				}

				occupationFile := filepath.Join(sessionFolder, "occupationDetails.json")
				if !os.IsNotExist(err) {
					content, err := os.ReadFile(occupationFile)
					if err != nil {
						log.Println(err)
					} else {
						err = json.Unmarshal(content, &occupationData)
						if err != nil {
							log.Println(err)
						}
					}

					os.Remove(occupationFile)
				}

				beneficiariesFile := filepath.Join(sessionFolder, "beneficiaries.json")
				if !os.IsNotExist(err) {
					content, err := os.ReadFile(beneficiariesFile)
					if err != nil {
						log.Println(err)
					} else {
						err = json.Unmarshal(content, &beneficiariesData)
						if err != nil {
							log.Println(err)
						}
					}

					os.Remove(beneficiariesFile)
				}

				mid, err := db.AddMember(memberData, contactsData, nomineeData, occupationData, beneficiariesData, nil)
				if err != nil {
					log.Println(err)
					return
				}

				id = *mid

				if os.Getenv("DEBUG") == "true" {
					if len(contactsData) > 0 {
						contactsData["memberId"] = id

						CacheFile(contactsFile, contactsData)
					}

					if len(nomineeData) > 0 {
						nomineeData["memberId"] = id

						CacheFile(nomineeFile, nomineeData)
					}

					if len(occupationData) > 0 {
						occupationData["memberId"] = id

						CacheFile(occupationFile, occupationData)
					}

					if len(beneficiariesData) > 0 {
						for i := range beneficiariesData {
							beneficiariesData[i]["memberId"] = id
						}

						CacheFile(beneficiariesFile, beneficiariesData)
					}
				}
			} else {
				mid, err := db.AddMember(memberData, nil, nil, nil, nil, nil)
				if err != nil {
					log.Println(err)
					return
				}

				id = *mid
			}

			menus.Sessions[*sessionId].MemberId = &id

			memberData["id"] = id

			menus.Sessions[*sessionId].ActiveMemberData = memberData

			if os.Getenv("DEBUG") == "true" {
				filename := filepath.Join(sessionFolder, "memberDetails.json")

				CacheFile(filename, memberData)
			}
		}

	case "contactDetails":
		val, ok := data.(map[string]any)
		if ok {
			if menus.Sessions[*sessionId].MemberId != nil {
				val["memberId"] = *menus.Sessions[*sessionId].MemberId

				_, err := db.AddMember(nil, val, nil, nil, nil, menus.Sessions[*sessionId].MemberId)
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				filename := filepath.Join(sessionFolder, "contactDetails.json")

				CacheFile(filename, val)
			}

			menus.Sessions[*sessionId].ContactsAdded = true
		}

	case "nomineeDetails":
		val, ok := data.(map[string]any)
		if ok {
			if menus.Sessions[*sessionId].MemberId != nil {
				val["memberId"] = *menus.Sessions[*sessionId].MemberId

				_, err := db.AddMember(nil, nil, val, nil, nil, menus.Sessions[*sessionId].MemberId)
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				filename := filepath.Join(sessionFolder, "nomineeDetails.json")

				CacheFile(filename, val)
			}

			menus.Sessions[*sessionId].NomineeAdded = true
		}

	case "occupationDetails":
		val, ok := data.(map[string]any)
		if ok {
			for _, key := range []string{"netPay", "grossPay"} {
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

			if menus.Sessions[*sessionId].MemberId != nil {
				val["memberId"] = *menus.Sessions[*sessionId].MemberId

				_, err := db.AddMember(nil, nil, nil, val, nil, menus.Sessions[*sessionId].MemberId)
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				filename := filepath.Join(sessionFolder, "occupationDetails.json")

				CacheFile(filename, val)
			}

			menus.Sessions[*sessionId].OccupationAdded = true
		}

	case "beneficiaries":
		rawData, ok := data.(map[string]any)
		if ok {
			records := []map[string]any{}

			for i := range 4 {
				var name, contact string
				var percentage float64

				index := i + 1

				nameLabel := fmt.Sprintf("beneficiary%vName", index)
				percentLabel := fmt.Sprintf("beneficiary%vPercent", index)
				contactLabel := fmt.Sprintf("beneficiary%vContact", index)

				if rawData[nameLabel] == nil {
					break
				}

				name = fmt.Sprintf("%v", rawData[nameLabel])
				contact = fmt.Sprintf("%v", rawData[contactLabel])

				v, err := strconv.ParseFloat(fmt.Sprintf("%v", rawData[percentLabel]), 64)
				if err == nil {
					percentage = v
				} else {
					log.Println(err)
				}

				row := map[string]any{
					"name":       name,
					"percentage": percentage,
					"contact":    contact,
				}

				if menus.Sessions[*sessionId].MemberId != nil {
					row["memberId"] = *menus.Sessions[*sessionId].MemberId
				}

				records = append(records, row)
			}

			if menus.Sessions[*sessionId].MemberId != nil {
				_, err := db.AddMember(nil, nil, nil, nil, records, menus.Sessions[*sessionId].MemberId)
				if err != nil {
					log.Println(err)
					return
				}
			} else {
				filename := filepath.Join(sessionFolder, "beneficiaries.json")

				CacheFile(filename, records)
			}

			menus.Sessions[*sessionId].BeneficiariesAdded = true
		}

	default:
		fmt.Println("##########", *phoneNumber, *sessionId, data)
	}
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
