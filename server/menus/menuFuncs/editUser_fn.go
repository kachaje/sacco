package menufuncs

import (
	"sacco/server/parser"
)

func EditUser(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder, cacheFolder string,
	) string,
	data map[string]any,
) string {
	var response string = "Edit User\n\n00. Main Menu"

	_ = data

	return response
}
