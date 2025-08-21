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
	"time"
)

func CacheFile(filename string, data any, retries int) {
	time.Sleep(time.Duration(retries) * time.Second)

	if utils.FileLocked(filename) {
		if retries < 5 {
			retries++

			CacheFile(filename, data, retries)
			return
		}
	}
	_, err := utils.LockFile(filename)
	if err != nil {
		log.Printf("server.Cachefile.1: %s", err.Error())
		retries = 0

		CacheFile(filename, data, retries)
		return
	}
	defer func() {
		err := utils.UnLockFile(filename)
		if err != nil {
			log.Printf("server.Cachefile.2: %s", err.Error())
		}
	}()

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

	case "memberBusiness":
		return HandleCommonModels(
			data, model, phoneNumber, sessionId, cacheFolder,
			func(
				data map[string]any,
				model string,
				parentId *int64,
			) (*int64, error) {
				return nil, nil
			},
			sessions, sessionFolder,
		)

	case "memberDetails":
		return handleMemberDetails(data, phoneNumber, sessionId, cacheFolder, saveFunc, sessions, sessionFolder)

	case "contactDetails":
		val, ok := data.(map[string]any)
		if ok {
			filename := filepath.Join(sessionFolder, "contactDetails.json")

			transactionDone := false

			// By default cache the data first in case we lose database connection
			CacheFile(filename, val, 0)
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
			CacheFile(filename, val, 0)
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
			CacheFile(filename, val, 0)
			defer func() {
				if transactionDone {
					os.Remove(filename)
				}
			}()

			for _, key := range []string{
				"netPay", "grossPay", "periodEmployed", "yearsInBusiness",
				"totalIncome", "totalCostOfGoods", "employeesWages", "ownSalary",
				"transport", "loanInterest", "utilities", "rentals", "otherCosts",
				"totalCosts", "netProfitLoss", "numberOfShares", "pricePerShare",
				"loanAmount", "repaymentPeriod", "amountRecommended",
				"amountApproved", "value",
			} {
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
		return handleBeneficiaries(data, phoneNumber, sessionId, cacheFolder, saveFunc, sessions, refData, sessionFolder)

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
