package menus

import (
	"fmt"
	"regexp"
	"sacco/server/parser"
	"slices"
)

func RegistrationMenu(session *parser.Session, phoneNumber, text, preferencesFolder, cacheFolder string, preferredLanguage *string) string {
	var response string

	switch text {
	case "00":
		session.WorkflowsMapping["member"].NavNext(text)
		session.CurrentMenu = "main"
		text = "0"
		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)

	case "1":
		session.CurrentMenu = "registration.1"
		if session.ActiveMemberData != nil {
			data := map[string]any{}

			if regexp.MustCompile(`^\d+$`).MatchString(phoneNumber) {
				data["phoneNumber"] = phoneNumber
			}

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

			session.WorkflowsMapping["member"].Data = data
		}

		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)

	case "2":
		session.CurrentMenu = "registration.2"

		if session.ActiveMemberData != nil && session.ActiveMemberData["memberOccupation"] != nil {
			val, ok := session.ActiveMemberData["memberOccupation"].(map[string]any)
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

				session.WorkflowsMapping["memberOccupation"].Data = data
			}
		}

		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)

	case "3":
		session.CurrentMenu = "registration.3"

		if session.ActiveMemberData != nil && session.ActiveMemberData["memberContact"] != nil {
			val, ok := session.ActiveMemberData["memberContact"].(map[string]any)
			if ok {
				data := map[string]any{}

				targetKeys := []string{
					"id",
					"memberId",
					"postalAddress",
					"residentialAddress",
					"phoneNumber",
					"homeVillage",
					"homeTraditionalAuthority",
					"homeDistrict",
				}
				for key, value := range val {
					if slices.Contains(targetKeys, key) {
						data[key] = fmt.Sprintf("%v", value)
					}
				}

				session.WorkflowsMapping["memberContact"].Data = data
			}
		}

		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)

	case "4":
		session.CurrentMenu = "registration.4"

		if session.ActiveMemberData != nil && session.ActiveMemberData["memberNominee"] != nil {
			val, ok := session.ActiveMemberData["memberNominee"].(map[string]any)
			if ok {
				data := map[string]any{}

				targetKeys := []string{
					"id",
					"memberId",
					"name",
					"phoneNumber",
					"address",
				}
				for key, value := range val {
					if slices.Contains(targetKeys, key) {
						data[key] = fmt.Sprintf("%v", value)
					}
				}

				session.WorkflowsMapping["memberNominee"].Data = data
			}
		}

		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)

	case "5":
		session.CurrentMenu = "registration.5"

		if session.ActiveMemberData != nil && session.ActiveMemberData["memberBeneficiary"] != nil {
			memberBeneficiary := []map[string]any{}

			val, ok := session.ActiveMemberData["memberBeneficiary"].([]map[string]any)
			if ok {
				memberBeneficiary = val
			} else {
				val, ok := session.ActiveMemberData["memberBeneficiary"].([]any)
				if ok {
					for _, row := range val {
						v, ok := row.(map[string]any)
						if ok {
							memberBeneficiary = append(memberBeneficiary, v)
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
				for i, val := range memberBeneficiary {
					for key, value := range val {
						if slices.Contains(targetKeys, key) {
							label := fmt.Sprintf("%s%d", key, i+1)

							data[label] = fmt.Sprintf("%v", value)
						}
					}
				}

				session.WorkflowsMapping["memberBeneficiary"].Data = data
			}
		}

		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)

	case "6":
		session.CurrentMenu = "registration.6"
		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)

	case "7":
		session.CurrentMenu = "registration.7"
		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)

	default:
		memberAdded := ""
		occupationAdded := ""
		contactAdded := ""
		nomineeAdded := ""
		beneficiariesAdded := ""
		businessInfoAdded := ""

		if session.MemberId != nil {
			if phoneNumber == "default" {
				memberAdded = "&#10003;"
			} else {
				memberAdded = "(*)"
			}
		}
		if session.AddedModels["memberOccupation"] {
			if phoneNumber == "default" {
				occupationAdded = "&#10003;"
			} else {
				occupationAdded = "(*)"
			}
		}
		if session.AddedModels["memberContact"] {
			if phoneNumber == "default" {
				contactAdded = "&#10003;"
			} else {
				contactAdded = "(*)"
			}
		}
		if session.AddedModels["memberNominee"] {
			if phoneNumber == "default" {
				nomineeAdded = "&#10003;"
			} else {
				nomineeAdded = "(*)"
			}
		}
		if session.AddedModels["memberBeneficiary"] {
			if phoneNumber == "default" {
				beneficiariesAdded = "&#10003;"
			} else {
				beneficiariesAdded = "(*)"
			}
		}
		if session.AddedModels["memberBusiness"] {
			if phoneNumber == "default" {
				businessInfoAdded = "&#10003;"
			} else {
				businessInfoAdded = "(*)"
			}
		}

		if preferredLanguage != nil && *preferredLanguage == "ny" {
			response = "CON Sankhani Zochita\n" +
				fmt.Sprintf("1. Zokhudza Membala %s\n", memberAdded) +
				fmt.Sprintf("2. Zokhudza Ntchito %s\n", occupationAdded) +
				fmt.Sprintf("3. Adiresi Yamembela\n %s", contactAdded) +
				fmt.Sprintf("4. Wachibale wa Membala %s\n", nomineeAdded) +
				fmt.Sprintf("5. Odzalandila %s\n", beneficiariesAdded) +
				fmt.Sprintf("6. Zabizinesi %s\n", businessInfoAdded) +
				"\n" +
				"7. Onani Zonse Zamembala\n" +
				"\n" +
				"00. Tiyambirenso"
		} else {
			response = "CON Choose Activity\n" +
				fmt.Sprintf("1. Member Details %s\n", memberAdded) +
				fmt.Sprintf("2. Occupation Details %s\n", occupationAdded) +
				fmt.Sprintf("3. Contact Details %s\n", contactAdded) +
				fmt.Sprintf("4. Next of Kin Details %s\n", nomineeAdded) +
				fmt.Sprintf("5. Beneficiaries %s\n", beneficiariesAdded) +
				fmt.Sprintf("6. Business Details %s\n", businessInfoAdded) +
				"\n" +
				"7. View Member Details\n" +
				"\n" +
				"00. Main Menu"
		}
	}

	return response
}
