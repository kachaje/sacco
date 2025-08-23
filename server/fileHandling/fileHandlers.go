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
