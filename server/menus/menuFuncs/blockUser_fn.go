package menufuncs

import (
	"sacco/server/database"
	"sacco/server/parser"
)

func BlockUser(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder, cacheFolder string,
	) string,
	db *database.Database,
	data map[string]any,
) string {
	var response string = "Block User\n\n00. Main Menu"

	_ = data

	return response
}
