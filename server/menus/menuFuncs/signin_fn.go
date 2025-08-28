package menufuncs

import (
	"fmt"
	"sacco/server/database"
	"sacco/server/parser"

	"github.com/google/uuid"
)

func SignIn(
	loadMenu func(
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

	if text == "00" {
		session.CurrentMenu = "main"
		return loadMenu("main", session, phoneNumber, "", preferencesFolder, cacheFolder)
	}

	if text == "" {
		response = "Login\n\nEnter username:\n"
	} else {
		if session.LastPrompt == "username" {
			session.Cache["username"] = text

			session.LastPrompt = "password"

			response = "Login\n\nEnter password:\n"
		} else {
			session.Cache["password"] = text

			text = ""

			if id, ok := db.ValidatePassword(session.Cache["username"], session.Cache["password"]); ok {
				token := uuid.NewString()
				session.SessionToken = &token
				session.SessionUserId = id

				session.CurrentMenu = "main"

				username := fmt.Sprintf("%v", session.Cache["username"])

				session.SessionUser = &username

				session.Cache = map[string]string{}
				session.LastPrompt = ""

				return loadMenu("main", session, phoneNumber, text, preferencesFolder, cacheFolder)
			} else {
				session.Cache = map[string]string{}
				session.LastPrompt = "username"

				response = "Login\n\nEnter username:\n"
			}
		}
	}

	response = fmt.Sprintf("%s\n00. Main Menu\n", response)

	return response
}
