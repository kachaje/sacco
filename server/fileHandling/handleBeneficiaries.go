package filehandling

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sacco/server/parser"
	"strconv"
)

func handleBeneficiaries(data any, phoneNumber, sessionId, cacheFolder *string,
	saveFunc func(
		a map[string]any,
		b map[string]any,
		c map[string]any,
		d map[string]any,
		e []map[string]any,
		f *int64,
	) (*int64, error), sessions map[string]*parser.Session, refData map[string]any, sessionFolder string) error {
	rawData, ok := data.(map[string]any)
	if ok {
		var memberId int64

		records := []map[string]any{}

		if sessions != nil && sessions[*sessionId] != nil &&
			sessions[*sessionId].MemberId != nil {
			memberId = *sessions[*sessionId].MemberId
		}

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
		CacheFile(filename, records, 0)
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

	return nil
}
