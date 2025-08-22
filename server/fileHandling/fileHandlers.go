package filehandling

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sacco/server/parser"
	"sacco/utils"
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
		map[string]any,
		string,
		int,
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

	case "memberBusiness", "memberOccupation", "memberNominee", "memberContact":
		return HandleCommonModels(
			data, model, phoneNumber, sessionId, cacheFolder,
			saveFunc, sessions, sessionFolder,
		)

	case "memberBeneficiary":
		return HandleBeneficiaries(data, phoneNumber, sessionId, cacheFolder, saveFunc, sessions, refData, sessionFolder)

	case "member":
		return HandleMemberDetails(data, phoneNumber, sessionId, cacheFolder, saveFunc, sessions, sessionFolder)

	// case "memberContact":
	// 	val, ok := data.(map[string]any)
	// 	if ok {
	// 		filename := filepath.Join(sessionFolder, "memberContact.json")

	// 		transactionDone := false

	// 		// By default cache the data first in case we lose database connection
	// 		CacheFile(filename, val, 0)
	// 		defer func() {
	// 			if transactionDone {
	// 				os.Remove(filename)
	// 			}
	// 		}()

	// 		if sessions[*sessionId].MemberId != nil {
	// 			val["memberId"] = *sessions[*sessionId].MemberId

	// 			if saveFunc == nil {
	// 				log.Println("Missing saveFunc")
	// 				return fmt.Errorf("server.SaveData.memberContact.1:missing saveFunc")
	// 			}

	// 			_, err := saveFunc(val, "memberContact", 0)
	// 			if err != nil {
	// 				return fmt.Errorf("server.SaveData.memberContact.2:%s", err.Error())
	// 			}

	// 			transactionDone = true
	// 		}

	// 		sessions[*sessionId].ActiveMemberData["memberContact"] = val

	// 		sessions[*sessionId].ContactsAdded = true

	// 		sessions[*sessionId].RefreshSession()

	// 		sessions[*sessionId].LoadMemberCache(*phoneNumber, *cacheFolder)
	// 	}

	// case "memberNominee":
	// 	val, ok := data.(map[string]any)
	// 	if ok {
	// 		filename := filepath.Join(sessionFolder, "memberNominee.json")

	// 		transactionDone := false

	// 		// By default cache the data first in case we lose database connection
	// 		CacheFile(filename, val, 0)
	// 		defer func() {
	// 			if transactionDone {
	// 				os.Remove(filename)
	// 			}
	// 		}()

	// 		if sessions[*sessionId].MemberId != nil {
	// 			val["memberId"] = *sessions[*sessionId].MemberId

	// 			if saveFunc == nil {
	// 				log.Println("Missing saveFunc")
	// 				return fmt.Errorf("server.SaveData.memberNominee.1:missing saveFunc")
	// 			}

	// 			_, err := saveFunc(val, "memberNominee", 0)
	// 			if err != nil {
	// 				return fmt.Errorf("server.SaveData.memberNominee.2:%s", err.Error())
	// 			}

	// 			transactionDone = true
	// 		}

	// 		sessions[*sessionId].ActiveMemberData["memberNominee"] = val

	// 		sessions[*sessionId].NomineeAdded = true

	// 		sessions[*sessionId].LoadMemberCache(*phoneNumber, *cacheFolder)
	// 	}

	// case "memberOccupation":
	// 	val, ok := data.(map[string]any)
	// 	if ok {
	// 		filename := filepath.Join(sessionFolder, "memberOccupation.json")

	// 		transactionDone := false

	// 		// By default cache the data first in case we lose database connection
	// 		CacheFile(filename, val, 0)
	// 		defer func() {
	// 			if transactionDone {
	// 				os.Remove(filename)
	// 			}
	// 		}()

	// 		for _, key := range []string{
	// 			"netPay", "grossPay", "periodEmployed", "yearsInBusiness",
	// 			"totalIncome", "totalCostOfGoods", "employeesWages", "ownSalary",
	// 			"transport", "loanInterest", "utilities", "rentals", "otherCosts",
	// 			"totalCosts", "netProfitLoss", "numberOfShares", "pricePerShare",
	// 			"loanAmount", "repaymentPeriod", "amountRecommended",
	// 			"amountApproved", "value",
	// 		} {
	// 			if val[key] != nil {
	// 				nv, ok := val[key].(string)
	// 				if ok {
	// 					real, err := strconv.ParseFloat(nv, 64)
	// 					if err == nil {
	// 						val[key] = real
	// 					}
	// 				}
	// 			}
	// 		}

	// 		if sessions[*sessionId].MemberId != nil {
	// 			val["memberId"] = *sessions[*sessionId].MemberId

	// 			if saveFunc == nil {
	// 				return fmt.Errorf("server.SaveData.memberOccupation.1:missing saveFunc")
	// 			}

	// 			_, err := saveFunc(val, "memberOccupation", 0)
	// 			if err != nil {
	// 				return fmt.Errorf("server.SaveData.memberOccupation.2:%s", err.Error())
	// 			}

	// 			transactionDone = true
	// 		}

	// 		sessions[*sessionId].ActiveMemberData["memberOccupation"] = val

	// 		sessions[*sessionId].OccupationAdded = true

	// 		sessions[*sessionId].RefreshSession()

	// 		sessions[*sessionId].LoadMemberCache(*phoneNumber, *cacheFolder)
	// 	}

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
