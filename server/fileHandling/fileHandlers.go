package filehandling

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sacco/server/database"
	"sacco/server/parser"
	"sacco/utils"
	"slices"
	"strconv"
	"time"
)

var (
	FloatFields = []string{
		"netPay", "grossPay", "periodEmployedInMonths", "yearsInBusiness",
		"totalIncome", "totalCostOfGoods", "employeesWages", "ownSalary",
		"transport", "loanInterest", "utilities", "rentals", "otherCosts",
		"totalCosts", "netProfitLoss", "numberOfShares", "pricePerShare",
		"loanAmount", "repaymentPeriodInMonths", "amountRecommended",
		"amountApproved", "value",
	}
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
	data any, model, phoneNumber, cacheFolder, preferenceFolder *string,
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
	case "language":
		val, ok := data.(map[string]any)
		if ok {
			if val["language"] != nil && phoneNumber != nil {
				language, ok := val["language"].(string)
				if ok {
					SavePreference(*phoneNumber, "language", language, *preferenceFolder)
				}
			}
		}

	case "memberBeneficiary":
		return HandleBeneficiaries(data, phoneNumber, cacheFolder, saveFunc, sessions, refData, sessionFolder)

	case "member":
		return HandleMemberDetails(data, phoneNumber, cacheFolder, saveFunc, sessions, sessionFolder)

	default:
		models := []string{}

		for _, group := range database.ArrayChildren {
			models = append(models, group...)
		}

		for _, group := range database.SingleChildren {
			models = append(models, group...)
		}

		if slices.Contains(models, *model) {
			return HandleCommonModels(
				data, model, phoneNumber, cacheFolder,
				saveFunc, sessions, sessionFolder,
			)
		} else {
			if val, ok := data.(map[string]any); ok {
				filename := filepath.Join(sessionFolder, fmt.Sprintf("%s.json", *model))

				transactionDone := false

				CacheFile(filename, val, 0)
				defer func() {
					if transactionDone {
						os.Remove(filename)
					}
				}()

				for _, key := range FloatFields {
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

				_, err := saveFunc(val, *model, 0)
				if err != nil {
					return err
				}

				transactionDone = true
			} else {
				fmt.Println("##########", *model, *phoneNumber, data)
			}
		}
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
