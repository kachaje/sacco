package menus

import "sacco/server/parser"

func EmploymentMenu(session *parser.Session, phoneNumber, text, preferencesFolder, cacheFolder string, preferredLanguage *string) string {
	var response string

	switch text {
	case "00":
		session.WorkflowsMapping["member"].NavNext(text)
		session.CurrentMenu = "main"
		text = "0"
		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)

	default:
		if text == "0" {
			session.CurrentMenu = "main"
			return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)
		} else {
			response = "CON Employment Details\n" +
				"1. Occupation Details\n" +
				"2. Occupation Verification\n" +
				"\n" +
				"00. Main Menu"
		}
	}

	return response
}
