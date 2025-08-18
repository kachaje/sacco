package menus

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sacco/parser"
	"strings"
	"sync"
)

type Session struct {
	CurrentMenu           string
	Data                  map[string]string
	PIWorkflow            *parser.WorkFlow
	LanguageWorkflow      *parser.WorkFlow
	OccupationWorkflow    *parser.WorkFlow
	ContactsWorkflow      *parser.WorkFlow
	NomineeWorkflow       *parser.WorkFlow
	BeneficiariesWorkflow *parser.WorkFlow
	PreferredLanguage     string
	MemberId              *int64
	SessionId             string
	PhoneNumber           string

	ContactsAdded      bool
	NomineeAdded       bool
	OccupationAdded    bool
	BeneficiariesAdded bool
	ActiveMemberData   map[string]any
}

var Sessions = make(map[string]*Session)
var mu sync.Mutex

func CheckPreferredLanguage(phoneNumber, preferencesFolder string) *string {
	settingsFile := filepath.Join(preferencesFolder, phoneNumber)

	_, err := os.Stat(settingsFile)
	if !os.IsNotExist(err) {
		content, err := os.ReadFile(settingsFile)
		if err != nil {
			log.Println(err)
			return nil
		}

		data := map[string]any{}

		err = json.Unmarshal(content, &data)
		if err != nil {
			log.Println(err)
			return nil
		}

		var preferredLanguage string

		if data["language"] != nil {
			val, ok := data["language"].(string)
			if ok {
				preferredLanguage = val
			}
		}

		return &preferredLanguage
	}

	return nil
}

func MainMenu(session *Session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder string) string {
	preferredLanguage := CheckPreferredLanguage(phoneNumber, preferencesFolder)

	if preferredLanguage != nil {
		session.PreferredLanguage = *preferredLanguage
	}

	var response string

	switch session.CurrentMenu {
	case "main":
		switch text {
		case "", "0":
			if preferredLanguage != nil && *preferredLanguage == "ny" {
				response = "CON Takulandilani ku Kaso SACCO\n" +
					"1. Membala Watsopano\n" +
					"2. Tengani Ngongole\n" +
					"3. Balansi\n" +
					"4. Matumizidwe\n" +
					"5. Chiyankhulo\n" +
					"6. Malizani"
			} else {
				response = "CON Welcome to Kaso SACCO\n" +
					"1. Membership Application\n" +
					"2. Loan Application\n" +
					"3. Check Balance\n" +
					"4. Banking Details\n" +
					"5. Preferred Language\n" +
					"6. Exit"
			}
		case "1":
			text = "000"
			session.CurrentMenu = "registration"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		case "2":
			text = "000"
			session.CurrentMenu = "loan"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		case "3":
			text = "000"
			session.CurrentMenu = "balance"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		case "4":
			text = "000"
			session.CurrentMenu = "banking"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		case "5":
			session.CurrentMenu = "language"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		case "6":
			if preferredLanguage != nil && *preferredLanguage == "ny" {
				response = "END Zikomo potidalila"
			} else {
				response = "END Thank you for using our service"
			}
			mu.Lock()
			delete(Sessions, sessionID)

			folderName := filepath.Join(cacheFolder, phoneNumber)

			_, err := os.Stat(folderName)
			if !os.IsNotExist(err) {
				files, err := os.ReadDir(folderName)
				if err == nil && len(files) == 0 {
					err = os.RemoveAll(folderName)
					if err != nil {
						log.Printf("server.MainMenu.removeFolder: %s\n", err.Error())
					}
				}
			}

			mu.Unlock()
		}
	case "language":
		if text == "" {
			session.CurrentMenu = "main"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		} else {
			response = session.LanguageWorkflow.NavNext(text)

			if strings.TrimSpace(response) == "" {
				session.CurrentMenu = "main"
				text = ""
				return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
			}
		}
	case "banking":
		if text == "0" {
			session.CurrentMenu = "main"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		} else {
			firstLine := "CON Banking Details\n"
			lastLine := "0. Back to Main Menu"
			name := "Name"
			number := "Number"
			branch := "Branch"

			if preferredLanguage != nil && *preferredLanguage == "ny" {
				firstLine = "CON Matumizidwe\n"
				lastLine = "0. Bwererani Pofikira"
				name = "Dzina"
				number = "Nambala"
				branch = "Buranchi"
			}

			switch text {
			case "1":
				response = "CON National Bank of Malawi\n" +
					fmt.Sprintf("%8s: Kaso SACCO\n", name) +
					fmt.Sprintf("%8s: 1006857589\n", number) +
					fmt.Sprintf("%8s: Lilongwe\n", branch) +
					lastLine
			case "2":
				response = "CON Airtel Money\n" +
					fmt.Sprintf("%8s: Kaso SACCO\n", name) +
					fmt.Sprintf("%8s: 0985 242 629\n", number) +
					lastLine
			default:
				response = firstLine +
					"1. National Bank\n" +
					"2. Airtel Money\n" +
					lastLine
			}
		}
	case "registration":
		return RegistrationMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder, preferredLanguage)

	case "registration.1":
		response = session.PIWorkflow.NavNext(text)

		if text == "00" {
			session.CurrentMenu = "main"
			text = "0"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		} else if strings.TrimSpace(response) == "" {
			session.CurrentMenu = "registration"
			text = ""
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		}

	case "registration.2":
		response = session.OccupationWorkflow.NavNext(text)

		if text == "00" {
			session.CurrentMenu = "main"
			text = "0"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		} else if strings.TrimSpace(response) == "" {
			session.CurrentMenu = "registration"
			text = ""
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		}

	case "registration.3":
		response = session.ContactsWorkflow.NavNext(text)

		if text == "00" {
			session.CurrentMenu = "main"
			text = "0"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		} else if strings.TrimSpace(response) == "" {
			session.CurrentMenu = "registration"
			text = ""
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		}

	case "registration.4":
		response = session.NomineeWorkflow.NavNext(text)

		if text == "00" {
			session.CurrentMenu = "main"
			text = "0"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		} else if strings.TrimSpace(response) == "" {
			session.CurrentMenu = "registration"
			text = ""
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		}

	case "registration.5":
		response = session.BeneficiariesWorkflow.NavNext(text)

		if text == "00" {
			session.CurrentMenu = "main"
			text = "0"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		} else if strings.TrimSpace(response) == "" {
			session.CurrentMenu = "registration"
			text = ""
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		}

	case "registration.6":
		if text == "00" {
			session.CurrentMenu = "main"
			text = "0"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		} else {
			if preferredLanguage != nil && *preferredLanguage == "ny" {
				response = "CON Zambiri za Membala\n" +
					"\n" +
					"00. Tiyambirenso"
			} else {
				response = "CON Member Details\n" +
					"\n" +
					"00. Main Menu"
			}
		}

	case "loan":
		if text == "0" {
			session.CurrentMenu = "main"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		} else {
			response = "CON Loan Application\n" +
				"0. Back to Main Menu"
		}

	case "balance":
		if text == "0" {
			session.CurrentMenu = "main"
			return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)
		} else {
			response = "CON Check Balance\n" +
				"0. Back to Main Menu"
		}
	}

	return response
}
