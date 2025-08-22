package filehandling

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sacco/server/parser"
	"sacco/utils"
	"strconv"
)

func RerunFailedSaves(phoneNumber, sessionId, cacheFolder *string,
	saveFunc func(
		a map[string]any,
		b string,
		c int,
	) (*int64, error), sessions map[string]*parser.Session) error {
	sessionFolder := filepath.Join(*cacheFolder, *phoneNumber)

	_, err := os.Stat(sessionFolder)
	if !os.IsNotExist(err) {
		retryNumberLockfile := filepath.Join(sessionFolder, fmt.Sprintf("%s.number", *phoneNumber))

		if utils.FileLocked(retryNumberLockfile) {
			return nil
		}

		_, err = utils.LockFile(retryNumberLockfile)
		if err != nil {
			log.Printf("server.RerunFailedSaves.LockFile: %s", err.Error())
			return err
		}
		defer func() {
			err = utils.UnLockFile(retryNumberLockfile)
			if err != nil {
				log.Printf("server.RerunFailedSaves.UnLockFile: %s", err.Error())
				return
			}
		}()

		targetFiles := []string{
			"memberDetails",
			"memberContact",
			"memberNominee",
			"memberOccupation",
			"memberBeneficiary",
		}

		for _, target := range targetFiles {
			filename := filepath.Join(sessionFolder, fmt.Sprintf("%s.json", target))

			// Priority is on writes this this can wait
			if utils.FileLocked(filename) {
				continue
			}

			_, err := os.Stat(filename)
			if !os.IsNotExist(err) {
				content, err := os.ReadFile(filename)
				if err != nil {
					log.Printf("server.RerunFailedSaves.1: %s", err.Error())
					continue
				}

				log.Printf("server.RerunFailedSaves: Retrying to save %s\n", target)

				if target == "memberBeneficiary" {
					rawData := []map[string]any{}

					err = json.Unmarshal(content, &rawData)
					if err != nil {
						log.Printf("server.RerunFailedSaves.2: %s", err.Error())
						continue
					}

					data := map[string]any{}

					for i, row := range rawData {
						index := i + 1

						nameLabel := fmt.Sprintf("name%v", index)
						percentLabel := fmt.Sprintf("percentage%v", index)
						contactLabel := fmt.Sprintf("contact%v", index)
						idLabel := fmt.Sprintf("id%v", index)
						memberIdLabel := fmt.Sprintf("memberId%v", index)

						if row["name"] != nil {
							data[nameLabel] = fmt.Sprintf("%v", row["name"])
						}
						if row["contact"] != nil {
							data[contactLabel] = fmt.Sprintf("%v", row["contact"])
						}
						if row["percentage"] != nil {
							v, err := strconv.ParseFloat(fmt.Sprintf("%v", row["percentage"]), 64)
							if err == nil {
								data[percentLabel] = v
							}
						}
						if row["id"] != nil {
							v, err := strconv.ParseInt(fmt.Sprintf("%v", row["id"]), 10, 64)
							if err == nil {
								data[idLabel] = v
							}
						}
						if row["memberId"] != nil {
							v, err := strconv.ParseInt(fmt.Sprintf("%v", row["memberId"]), 10, 64)
							if err == nil {
								data[memberIdLabel] = v
							}
						}
					}

					err = SaveData(data, &target, phoneNumber, sessionId, cacheFolder, nil, saveFunc, sessions, nil)
					if err != nil {
						log.Printf("server.RerunFailedSaves.3: %s", err.Error())
						continue
					}
				} else {
					data := map[string]any{}

					err = json.Unmarshal(content, &data)
					if err != nil {
						log.Printf("server.RerunFailedSaves.4: %s", err.Error())
						continue
					}

					err = SaveData(data, &target, phoneNumber, sessionId, cacheFolder, nil, saveFunc, sessions, nil)
					if err != nil {
						log.Printf("server.RerunFailedSaves.5: %s", err.Error())
						continue
					}
				}
			}
		}
	}

	return nil
}
