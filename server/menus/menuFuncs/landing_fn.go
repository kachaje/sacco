package menufuncs

import (
	"sacco/server/database"
	"sacco/server/parser"
)

func Landing(loadMenu func(
	menuName string, session *parser.Session,
	phoneNumber, text, preferencesFolder, cacheFolder string,
) string,
	db *database.Database,
	data map[string]any,
) string {
	var response string
	var phoneNumber, text, preferencesFolder, cacheFolder string
	var session *parser.Session

	if data["session"] != nil {
		if val, ok := data["session"].(*parser.Session); ok {
			session = val
		}
	}
	if data["phoneNumber"] != nil {
		if val, ok := data["phoneNumber"].(string); ok {
			phoneNumber = val
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
	if data["cacheFolder"] != nil {
		if val, ok := data["cacheFolder"].(string); ok {
			cacheFolder = val
		}
	}

	data = map[string]any{
		"phoneNumber":       phoneNumber,
		"cacheFolder":       cacheFolder,
		"session":           session,
		"preferencesFolder": preferencesFolder,
		"text":              text,
	}

	session.LastPrompt = ""
	session.Cache = map[string]string{}

	switch text {
	case "1":
		session.CurrentMenu = "signIn"
		data["text"] = ""
		return SignIn(loadMenu, db, data)
	case "2":
		session.CurrentMenu = "signUp"
		session.LastPrompt = "username"
		data["text"] = ""
		return SignUp(loadMenu, db, data)
	default:
		response = "Welcome! Select Action\n\n" +
			"1. Sign In\n" +
			"2. Sign Up\n"
	}

	return response
}
