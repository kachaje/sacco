package menus

import "sacco/server/parser"

func BusinessMenu(session *parser.Session, phoneNumber, text, preferencesFolder, cacheFolder string, preferredLanguage *string) string {
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
			response = "CON Business\n" +
				"1. Business Details\n" +
				"2. Previous Year History\n" +
				"3. Next Year Projection\n" +
				"\n" +
				"00. Main Menu"
		}
	}

	return response
}
