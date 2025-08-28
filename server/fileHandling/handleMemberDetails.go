package filehandling

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sacco/server/database"
	"sacco/server/parser"
)

func HandleMemberDetails(data any, phoneNumber, cacheFolder *string,
	saveFunc func(map[string]any, string, int) (*int64, error), sessions map[string]*parser.Session, sessionFolder string) error {
	memberData, ok := data.(map[string]any)
	if ok {
		var id int64

		if phoneNumber != nil && *phoneNumber != "default" {
			if memberData["phoneNumber"] == nil {
				memberData["phoneNumber"] = *phoneNumber
			}
		}

		if sessions[*phoneNumber] != nil {
			filename := filepath.Join(sessionFolder, "member.json")

			transactionDone := false

			// By default cache the data first in case we lose database connection
			CacheFile(filename, memberData, 0)
			defer func() {
				if transactionDone {
					os.Remove(filename)
				}
			}()

			someChildAdded := false

			for _, key := range database.MemberSingleChildren {
				if sessions[*phoneNumber].AddedModels[key] {
					someChildAdded = true
					break
				}
			}

			if someChildAdded {
				if saveFunc == nil {
					return fmt.Errorf("server.SaveData.member.9:missing saveFunc")
				}

				mid, err := saveFunc(memberData, "member", 0)
				if err != nil {
					log.Println(err)
					return err
				}

				id = *mid

				models := []string{}

				models = append(models, database.MemberArrayChildren...)

				models = append(models, database.MemberSingleChildren...)

				for _, model := range models {
					file := filepath.Join(sessionFolder, fmt.Sprintf("%s.json", model))

					_, err = os.Stat(file)
					if !os.IsNotExist(err) {
						fileArrayData := []map[string]any{}
						fileFlatData := map[string]any{}

						content, err := os.ReadFile(file)
						if err != nil {
							log.Printf("server.HandleMemberDetails:%s\n", err.Error())
						} else {
							err = json.Unmarshal(content, &fileArrayData)
							if err != nil {
								err = json.Unmarshal(content, &fileFlatData)
								if err != nil {
									log.Printf("server.HandleMemberDetails:%s\n", err.Error())
								}
							}
						}

						if len(fileArrayData) > 0 {
							for i := range fileArrayData {
								fileArrayData[i]["memberId"] = id

								mid, err = saveFunc(fileArrayData[i], model, 0)
								if err != nil {
									log.Printf("server.HandleMemberDetails:%s\n", err.Error())
									continue
								}

								switch model {
								case "memberBusiness":
									sessions[*phoneNumber].MemberBusinessId = mid
								case "memberLoan":
									sessions[*phoneNumber].MemberLoanId = mid
								}
							}

							memberData[model] = fileArrayData

							if os.Getenv("DEBUG") == "true" {
								CacheFile(file, fileArrayData, 0)
							} else {
								os.Remove(file)
							}
						} else if len(fileFlatData) > 0 {
							fileFlatData["memberId"] = id

							mid, err = saveFunc(fileFlatData, model, 0)
							if err != nil {
								log.Printf("server.HandleMemberDetails:%s\n", err.Error())
								continue
							}

							switch model {
							case "memberBusiness":
								sessions[*phoneNumber].MemberBusinessId = mid
							case "memberLoan":
								sessions[*phoneNumber].MemberLoanId = mid
							}

							memberData[model] = fileFlatData

							if os.Getenv("DEBUG") == "true" {
								CacheFile(file, fileFlatData, 0)
							} else {
								os.Remove(file)
							}
						}
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

			sessions[*phoneNumber].MemberId = &id

			memberData["id"] = id

			sessions[*phoneNumber].UpdateActiveData(memberData, 0)

			sessions[*phoneNumber].RefreshSession()

			sessions[*phoneNumber].LoadCacheData(*phoneNumber, *cacheFolder)
		}
	}

	return nil
}
