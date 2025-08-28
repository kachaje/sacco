package menufuncs

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sacco/server/database"
	"sacco/server/parser"
)

func DevConsole(
	loadMenu func(
		menuName string, session *parser.Session,
		phoneNumber, text, preferencesFolder, cacheFolder string,
	) string,
	db *database.Database,
	data map[string]any,
) string {
	var session *parser.Session
	var response, content, text, title string

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

	if session != nil {
		if text == "99" {
			session.CurrentMenu = "console"
		} else if session.CurrentMenu == "console" && regexp.MustCompile(`^\d+$`).MatchString(text) {
			session.CurrentMenu = fmt.Sprintf("%s.%s", session.CurrentMenu, text)
			text = ""
		}

		switch session.CurrentMenu {
		case "console.1":
			title = "WorkflowsMapping"

			if session.WorkflowsMapping != nil {
				data := map[string]any{}

				for key, wflow := range session.WorkflowsMapping {
					row := map[string]any{
						"data":           wflow.Data,
						"optionalFields": wflow.OptionalFields,
						"screenOrder":    wflow.ScreenOrder,
						"history":        wflow.History,
					}

					data[key] = row
				}

				payload, err := json.MarshalIndent(data, "", "  ")
				if err != nil {
					content = err.Error()
				} else {
					content = string(payload)
				}
			}
		case "console.2":
			title = "AddedModels"

			if session.AddedModels != nil {
				payload, err := json.MarshalIndent(session.AddedModels, "", "  ")
				if err != nil {
					content = err.Error()
				} else {
					content = string(payload)
				}
			}
		case "console.3":
			title = "ActiveData"

			if session.ActiveData != nil {
				payload, err := json.MarshalIndent(session.ActiveData, "", "  ")
				if err != nil {
					content = err.Error()
				} else {
					content = string(payload)
				}
			}
		case "console.4":
			title = "Data"

			if session.Data != nil {
				payload, err := json.MarshalIndent(session.Data, "", "  ")
				if err != nil {
					content = err.Error()
				} else {
					content = string(payload)
				}
			}
		case "console.5":
			title = "MemberId"

			if session.MemberId != nil {
				content = fmt.Sprint(*session.MemberId)
			}
		case "console.6":
			title = "SessionId"

			content = session.SessionId
		case "console.7":
			title = "PhoneNumber"

			content = session.PhoneNumber

		case "console.8":
			title = "SQL Query"

			result, err := db.SQLQuery(text)
			if err != nil {
				content = fmt.Sprintf(`query: %s

response: %s`, text, err.Error())
			} else {
				payload, err := json.MarshalIndent(result, "", " ")
				if err != nil {
					content = fmt.Sprintf(`query: %s

response: %s`, text, err.Error())
				} else {
					content = fmt.Sprintf(`query: %s

response: %s`, text, payload)
				}
			}

		case "console.9":
			title = "Member By PhoneNumber"

			if text != "" {
				result, err := db.MemberByPhoneNumber(text, nil, nil)
				if err != nil {
					content = fmt.Sprintf("%s\n", err.Error())
				} else {
					payload, err := json.MarshalIndent(result, "", " ")
					if err != nil {
						content = fmt.Sprintf("%s\n", err.Error())
					} else {
						content = fmt.Sprintf("%s\n", payload)
					}
				}
			} else {
				content = ""
			}

		default:
			session.CurrentMenu = "console"

			content = "Available Dumps:\n" +
				"1. WorkflowsMapping\n" +
				"2. AddedModels\n" +
				"3. ActiveData\n" +
				"4. Data\n" +
				"5. MemberId\n" +
				"6. SessionId\n" +
				"7. PhoneNumber\n" +
				"8. SQL Query\n" +
				"9. Member By PhoneNumber"
		}
	} else {
		content = "No active session provided"
	}

	response = "Dev Console\n\n" +
		title +
		fmt.Sprintf("\n\n%s\n", content) +
		"\n99. Cancel\n" +
		"00. Main Menu\n"

	return response
}
