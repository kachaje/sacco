package menus

import "fmt"

func RegistrationMenu(session *Session, phoneNumber, text, sessionID, preferencesFolder string, preferredLanguage *string) string {
	var response string

	switch text {
	case "00":
		session.PIWorkflow.NavNext(text)
		session.CurrentMenu = "main"
		text = "0"
		return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder)
	case "1":
		session.CurrentMenu = "registration.1"
		return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder)
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
				"\n" +
				"00. Tiyambirenso"
		} else {
			response = "CON Choose Activity\n" +
				fmt.Sprintf("1. Add Member Details %s\n", memberAdded) +
				fmt.Sprintf("2. Add Occupation Details %s\n", occupationAdded) +
				fmt.Sprintf("3. Add Contact Details %s\n", contactAdded) +
				fmt.Sprintf("4. Add Next of Kin Details %s\n", nomineeAdded) +
				fmt.Sprintf("5. Add Beneficiaries %s\n", beneficiariesAdded) +
				"\n" +
				"00. Main Menu"
		}
	}

	return response
}
