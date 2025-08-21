package menus

import (
	"fmt"
	"sacco/parser"
	"slices"
)

func RegistrationMenu(session *parser.Session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder string, preferredLanguage *string) string {
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
				"utilityBillNumber", "utilityBillType", "id", "phoneNumber",
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

		if session.ActiveMemberData != nil && session.ActiveMemberData["occupationDetails"] != nil {
			val, ok := session.ActiveMemberData["occupationDetails"].(map[string]any)
			if ok {
				data := map[string]any{}

				targetKeys := []string{
					"employerAddress", "employerName", "employerPhone",
					"grossPay", "highestQualification", "jobTitle", "netPay",
					"periodEmployed", "id",
				}
				for key, value := range val {
					if slices.Contains(targetKeys, key) {
						data[key] = fmt.Sprintf("%v", value)
					}
				}

				session.OccupationWorkflow.Data = data
			}
		}

		return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)

	case "3":
		session.CurrentMenu = "registration.3"

		if session.ActiveMemberData != nil && session.ActiveMemberData["contactDetails"] != nil {
			val, ok := session.ActiveMemberData["contactDetails"].(map[string]any)
			if ok {
				data := map[string]any{}

				targetKeys := []string{
					"id",
					"memberId",
					"postalAddress",
					"residentialAddress",
					"phoneNumber",
					"homeVillage",
					"homeTA",
					"homeDistrict",
				}
				for key, value := range val {
					if slices.Contains(targetKeys, key) {
						data[key] = fmt.Sprintf("%v", value)
					}
				}

				session.ContactsWorkflow.Data = data
			}
		}

		return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)

	case "4":
		session.CurrentMenu = "registration.4"

		if session.ActiveMemberData != nil && session.ActiveMemberData["nomineeDetails"] != nil {
			val, ok := session.ActiveMemberData["nomineeDetails"].(map[string]any)
			if ok {
				data := map[string]any{}

				targetKeys := []string{
					"id",
					"memberId",
					"nomineeName",
					"nomineePhone",
					"nomineeAddress",
				}
				for key, value := range val {
					if slices.Contains(targetKeys, key) {
						data[key] = fmt.Sprintf("%v", value)
					}
				}

				session.NomineeWorkflow.Data = data
			}
		}

		return MainMenu(session, phoneNumber, text, sessionID, preferencesFolder, cacheFolder)

	case "5":
		session.CurrentMenu = "registration.5"

		if session.ActiveMemberData != nil && session.ActiveMemberData["beneficiaries"] != nil {
			beneficiaries := []map[string]any{}

			val, ok := session.ActiveMemberData["beneficiaries"].([]map[string]any)
			if ok {
				beneficiaries = val
			} else {
				val, ok := session.ActiveMemberData["beneficiaries"].([]any)
				if ok {
					for _, row := range val {
						v, ok := row.(map[string]any)
						if ok {
							beneficiaries = append(beneficiaries, v)
						}
					}
				}
			}

			{
				data := map[string]any{}

				targetKeys := []string{
					"id",
					"name",
					"percentage",
					"contact",
				}
				for i, val := range beneficiaries {
					for key, value := range val {
						if slices.Contains(targetKeys, key) {
							label := fmt.Sprintf("%s%d", key, i+1)

							data[label] = fmt.Sprintf("%v", value)
						}
					}
				}

				session.BeneficiariesWorkflow.Data = data
			}
		}

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
