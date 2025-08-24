package menus

import (
	"sacco/server/parser"
	"strings"
)

func BusinessMenu(session *parser.Session, phoneNumber, text, preferencesFolder, cacheFolder string, preferredLanguage *string) string {
	var response string

	switch text {
	case "00":
		session.WorkflowsMapping["member"].NavNext(text)
		session.CurrentMenu = "main"
		text = "0"
		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)

	case "1":
		session.CurrentMenu = "business.1"

		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)

	case "2":
		session.CurrentMenu = "business.2"

		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)

	case "3":
		session.CurrentMenu = "business.3"

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

func BusinessMenu1(session *parser.Session, phoneNumber, text, preferencesFolder, cacheFolder string, preferredLanguage *string) string {
	response := session.WorkflowsMapping["memberBusiness"].NavNext(text)

	if text == "00" {
		session.CurrentMenu = "main"
		text = "0"
		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)
	} else if strings.TrimSpace(response) == "" {
		session.CurrentMenu = "business"
		text = ""
		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)
	}

	return response
}

func BusinessMenu2(session *parser.Session, phoneNumber, text, preferencesFolder, cacheFolder string, preferredLanguage *string) string {
	response := session.WorkflowsMapping["memberLastYearBusinessHistory"].NavNext(text)

	if text == "00" {
		session.CurrentMenu = "main"
		text = "0"
		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)
	} else if strings.TrimSpace(response) == "" {
		session.CurrentMenu = "business"
		text = ""
		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)
	}

	return response
}

func BusinessMenu3(session *parser.Session, phoneNumber, text, preferencesFolder, cacheFolder string, preferredLanguage *string) string {
	response := session.WorkflowsMapping["memberNextYearBusinessProjection"].NavNext(text)

	if text == "00" {
		session.CurrentMenu = "main"
		text = "0"
		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)
	} else if strings.TrimSpace(response) == "" {
		session.CurrentMenu = "business"
		text = ""
		return MainMenu(session, phoneNumber, text, preferencesFolder, cacheFolder)
	}

	return response
}
