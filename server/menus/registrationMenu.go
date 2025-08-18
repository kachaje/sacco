package menus

import (
	"fmt"
	"slices"
)

func RegistrationMenu(session *Session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder string, preferredLanguage *string) string {
	var response string

	switch text {
	case "00":
		session.PIWorkflow.NavNext(text)
		session.CurrentMenu = "main"
		text = "0"
		return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)

	case "1":
		session.CurrentMenu = "registration.1"
		if session.ActiveMemberData != nil {
			data := map[string]any{}

			targetKeys := []string{
				"dateOfBirth", "firstName", "gender", "lastName",
				"maritalStatus", "nationalId", "otherName", "title",
				"utilityBillNumber", "utilityBillType", "id",
			}
			for key, value := range session.ActiveMemberData {
				if slices.Contains(targetKeys, key) {
					data[key] = fmt.Sprintf("%v", value)
				}
			}

			session.PIWorkflow.Data = data
		}

		return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)

	case "2":
		session.CurrentMenu = "registration.2"
		return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)

	case "3":
		session.CurrentMenu = "registration.3"
		return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)

	case "4":
		session.CurrentMenu = "registration.4"
		return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)

	case "5":
		session.CurrentMenu = "registration.5"
		return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)

	case "6":
		session.CurrentMenu = "registration.6"
		return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)

	default:
		memberAdded := ""
		occupationAdded := ""
		contactAdded := ""
		nomineeAdded := ""
		beneficiariesAdded := ""

		if session.MemberId != nil {
			if phoneNumber == "default" {
				memberAdded = "&#10003;"
			} else {
				memberAdded = "(*)"
			}
		}
		if session.OccupationAdded {
			if phoneNumber == "default" {
				occupationAdded = "&#10003;"
			} else {
				occupationAdded = "(*)"
			}
		}
		if session.ContactsAdded {
			if phoneNumber == "default" {
				contactAdded = "&#10003;"
			} else {
				contactAdded = "(*)"
			}
		}
		if session.NomineeAdded {
			if phoneNumber == "default" {
				nomineeAdded = "&#10003;"
			} else {
				nomineeAdded = "(*)"
			}
		}
		if session.BeneficiariesAdded {
			if phoneNumber == "default" {
				beneficiariesAdded = "&#10003;"
			} else {
				beneficiariesAdded = "(*)"
			}
		}

		if preferredLanguage != nil && *preferredLanguage == "ny" {
			response = "CON Sankhani Zochita\n" +
				fmt.Sprintf("1. Zokhudza Membala %s\n", memberAdded) +
				fmt.Sprintf("2. Zokhudza Ntchito %s\n", occupationAdded) +
				fmt.Sprintf("3. Adiresi Yamembela\n %s", contactAdded) +
				fmt.Sprintf("4. Wachibale wa Membala %s\n", nomineeAdded) +
				fmt.Sprintf("5. Odzalandila %s\n", beneficiariesAdded) +
				"6. Onani Zonse Zamembala\n" +
				"\n" +
				"00. Tiyambirenso"
		} else {
			response = "CON Choose Activity\n" +
				fmt.Sprintf("1. Member Details %s\n", memberAdded) +
				fmt.Sprintf("2. Occupation Details %s\n", occupationAdded) +
				fmt.Sprintf("3. Contact Details %s\n", contactAdded) +
				fmt.Sprintf("4. Next of Kin Details %s\n", nomineeAdded) +
				fmt.Sprintf("5. Beneficiaries %s\n", beneficiariesAdded) +
				"6. View Member Details\n" +
				"\n" +
				"00. Main Menu"
		}
	}

	return response
}
