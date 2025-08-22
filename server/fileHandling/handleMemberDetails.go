package filehandling

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sacco/server/parser"
)

func HandleMemberDetails(data any, phoneNumber, sessionId, cacheFolder *string,
	saveFunc func(map[string]any, string, int) (*int64, error), sessions map[string]*parser.Session, sessionFolder string) error {
	memberData, ok := data.(map[string]any)
	if ok {
		var id int64
		var err error

		if phoneNumber != nil && *phoneNumber != "default" {
			if memberData["phoneNumber"] == nil {
				memberData["phoneNumber"] = *phoneNumber
			}
		}

		filename := filepath.Join(sessionFolder, "member.json")

		transactionDone := false

		// By default cache the data first in case we lose database connection
		CacheFile(filename, memberData, 0)
		defer func() {
			if transactionDone {
				os.Remove(filename)
			}
		}()

		if sessions[*sessionId].AddedModels["memberContact"] ||
			sessions[*sessionId].AddedModels["memberNominee"] ||
			sessions[*sessionId].AddedModels["memberOccupation"] ||
			sessions[*sessionId].AddedModels["memberBeneficiary"] {

			var contactsData, nomineeData, occupationData map[string]any
			var beneficiariesData []map[string]any

			contactsFile := filepath.Join(sessionFolder, "memberContact.json")

			_, err = os.Stat(contactsFile)
			if !os.IsNotExist(err) {
				content, err := os.ReadFile(contactsFile)
				if err != nil {
					log.Printf("server.SaveData.member.1:%s\n", err.Error())
				} else {
					err = json.Unmarshal(content, &contactsData)
					if err != nil {
						log.Printf("server.SaveData.member.2:%s\n", err.Error())
					}
				}
			}

			nomineeFile := filepath.Join(sessionFolder, "memberNominee.json")

			_, err = os.Stat(nomineeFile)
			if !os.IsNotExist(err) {
				content, err := os.ReadFile(nomineeFile)
				if err != nil {
					log.Printf("server.SaveData.member.3:%s\n", err.Error())
				} else {
					err = json.Unmarshal(content, &nomineeData)
					if err != nil {
						log.Printf("server.SaveData.member.4:%s\n", err.Error())
					}
				}
			}

			occupationFile := filepath.Join(sessionFolder, "memberOccupation.json")

			_, err = os.Stat(occupationFile)
			if !os.IsNotExist(err) {
				content, err := os.ReadFile(occupationFile)
				if err != nil {
					log.Printf("server.SaveData.member.5:%s\n", err.Error())
				} else {
					err = json.Unmarshal(content, &occupationData)
					if err != nil {
						log.Printf("server.SaveData.member.6:%s\n", err.Error())
					}
				}
			}

			beneficiariesFile := filepath.Join(sessionFolder, "memberBeneficiary.json")

			_, err = os.Stat(beneficiariesFile)
			if !os.IsNotExist(err) {
				content, err := os.ReadFile(beneficiariesFile)
				if err != nil {
					log.Printf("server.SaveData.member.7:%s\n", err.Error())
				} else {
					err = json.Unmarshal(content, &beneficiariesData)
					if err != nil {
						log.Printf("server.SaveData.member.8:%s\n", err.Error())
					}
				}
			}

			if saveFunc == nil {
				return fmt.Errorf("server.SaveData.member.9:missing saveFunc")
			}

			mid, err := saveFunc(memberData, "member", 0)
			if err != nil {
				log.Println(err)
				return err
			}

			id = *mid

			if len(contactsData) > 0 {
				contactsData["memberId"] = id

				memberData["memberContact"] = contactsData

				_, err = saveFunc(contactsData, "memberContact", 0)
				if err != nil {
					log.Println(err)
					return err
				}

				if os.Getenv("DEBUG") == "true" {
					CacheFile(contactsFile, contactsData, 0)
				} else {
					os.Remove(contactsFile)
				}
			}

			if len(nomineeData) > 0 {
				nomineeData["memberId"] = id

				memberData["memberNominee"] = nomineeData

				_, err = saveFunc(nomineeData, "memberNominee", 0)
				if err != nil {
					log.Println(err)
					return err
				}

				if os.Getenv("DEBUG") == "true" {
					CacheFile(nomineeFile, nomineeData, 0)
				} else {
					os.Remove(nomineeFile)
				}
			}

			if len(occupationData) > 0 {
				occupationData["memberId"] = id

				memberData["memberOccupation"] = occupationData

				_, err = saveFunc(occupationData, "memberOccupation", 0)
				if err != nil {
					log.Println(err)
					return err
				}

				if os.Getenv("DEBUG") == "true" {
					CacheFile(occupationFile, occupationData, 0)
				} else {
					os.Remove(occupationFile)
				}
			}

			if len(beneficiariesData) > 0 {
				for i := range beneficiariesData {
					beneficiariesData[i]["memberId"] = id

					_, err = saveFunc(beneficiariesData[i], "memberBeneficiary", 0)
					if err != nil {
						log.Println(err)
						continue
					}
				}

				memberData["memberBeneficiary"] = beneficiariesData

				if os.Getenv("DEBUG") == "true" {
					CacheFile(beneficiariesFile, beneficiariesData, 0)
				} else {
					os.Remove(beneficiariesFile)
				}
			}

			transactionDone = true
		} else {
			if saveFunc == nil {
				return fmt.Errorf("server.SaveData.member.10:missing saveFunc")
			}

			mid, err := saveFunc(memberData, "member", 0)
			if err != nil {
				log.Println(err)
				return fmt.Errorf("server.SaveData.member.11:%s", err.Error())
			}

			id = *mid

			transactionDone = true
		}

		sessions[*sessionId].MemberId = &id

		memberData["id"] = id

		sessions[*sessionId].UpdateActiveMemberData(memberData, 0)

		sessions[*sessionId].RefreshSession()

		sessions[*sessionId].LoadMemberCache(*phoneNumber, *cacheFolder)
	}

	return nil
}
