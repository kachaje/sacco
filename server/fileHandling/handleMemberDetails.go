package filehandling

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

			for _, key := range parser.MemberChildren {
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

				for _, model := range parser.MemberArrayChildren {
					file := filepath.Join(sessionFolder, fmt.Sprintf("%s.json", model))

					_, err = os.Stat(file)
					if !os.IsNotExist(err) {
						fileData := []map[string]any{}

						content, err := os.ReadFile(file)
						if err != nil {
							log.Printf("server.HandleMemberDetails:%s\n", err.Error())
						} else {
							err = json.Unmarshal(content, &fileData)
							if err != nil {
								log.Printf("server.HandleMemberDetails:%s\n", err.Error())
							}
						}

						if len(fileData) > 0 {
							for i := range fileData {
								fileData[i]["memberId"] = id

								_, err = saveFunc(fileData[i], model, 0)
								if err != nil {
									log.Printf("server.HandleMemberDetails:%s\n", err.Error())
									continue
								}
							}

							memberData[model] = fileData

							if os.Getenv("DEBUG") == "true" {
								CacheFile(file, fileData, 0)
							} else {
								os.Remove(file)
							}
						}
					}
				}

				for _, model := range parser.MemberChildren {
					file := filepath.Join(sessionFolder, fmt.Sprintf("%s.json", model))

					_, err = os.Stat(file)
					if !os.IsNotExist(err) {
						fileData := map[string]any{}

						content, err := os.ReadFile(file)
						if err != nil {
							log.Printf("server.HandleMemberDetails:%s\n", err.Error())
						} else {
							err = json.Unmarshal(content, &fileData)
							if err != nil {
								log.Printf("server.HandleMemberDetails:%s\n", err.Error())
							}
						}

						if len(fileData) > 0 {
							fileData["memberId"] = id

							_, err = saveFunc(fileData, model, 0)
							if err != nil {
								log.Printf("server.HandleMemberDetails:%s\n", err.Error())
								continue
							}
						}

						memberData[model] = fileData

						if os.Getenv("DEBUG") == "true" {
							CacheFile(file, fileData, 0)
						} else {
							os.Remove(file)
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

			sessions[*phoneNumber].UpdateActiveMemberData(memberData, 0)

			sessions[*phoneNumber].RefreshSession()

			sessions[*phoneNumber].LoadMemberCache(*phoneNumber, *cacheFolder)
		}
	}

	return nil
}
