package menufuncs

import (
	"sacco/server/parser"
)

func BusinessSummary(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder, cacheFolder string,
	) string,
	data map[string]any,
) string {
	var result string = "Business Summary\n\n" +
		"00. Main Menu\n"

	_ = data

	return result
}
