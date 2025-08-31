package menufuncs

import (
	"fmt"
	"regexp"
	filehandling "sacco/server/fileHandling"
	"sacco/server/parser"
)

func CreateNewSession(phoneNumber, sessionId, preferencesFolder, preferredLanguage string, demoMode bool) *parser.Session {
	mu.Lock()
	session, exists := Sessions[phoneNumber]
	if !exists {
		session = parser.NewSession(DB.MemberByPhoneNumber, &phoneNumber, &sessionId)

		for model, data := range WorkflowsData {
			session.WorkflowsMapping[model] = parser.NewWorkflow(data, filehandling.SaveData, &preferredLanguage, &phoneNumber, &sessionId, &preferencesFolder, DB.GenericsSaveData, Sessions, nil)
		}

		if preferredLanguage != "" {
			session.PreferredLanguage = preferredLanguage
		}

		if demoMode {
			defaultUser := "admin"
			defaultUserId := int64(1)
			defaultRole := "admin"

			session.SessionUser = &defaultUser
			session.SessionUserId = &defaultUserId
			session.SessionUserRole = &defaultRole
		}

		Sessions[phoneNumber] = session
	}
	mu.Unlock()

	return session
}

func SetPhoneNumber(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder string,
	) string,
	data map[string]any,
) string {
	var response string
	var content, text, preferencesFolder string
	var session *parser.Session

	title := "CON Set PhoneNumber\n\n"
	footer := "\n00. Main Menu\n"

	if data["session"] != nil {
		if val, ok := data["session"].(*parser.Session); ok {
			session = val
		}
	}
	if data["text"] != nil {
		if val, ok := data["text"].(string); ok {
			text = val
		}
	}
	if data["preferencesFolder"] != nil {
		if val, ok := data["preferencesFolder"].(string); ok {
			preferencesFolder = val
		}
	}

	if text == "00" {
		session.CurrentMenu = "main"
		return loadMenu("main", session, session.PhoneNumber, "", preferencesFolder)
	}

	askPhoneNumber := func(msg string) string {
		return fmt.Sprintf("Enter phone number: %s\n", msg)
	}

	if text != "" && text != "000" {
		if !regexp.MustCompile(`^\d+$`).MatchString(text) {
			content = askPhoneNumber("(Invalid input)")
		} else {
			session.PhoneNumber = text

			_, err := session.RefreshSession()
			if err == nil {
				session.UpdateSessionFlags()
			}

			text = ""
			content = "Success. Phone Number set!\n"
		}
	} else {
		content = askPhoneNumber("")
	}

	response = fmt.Sprintf("%s%s%s", title, content, footer)

	return response
}
