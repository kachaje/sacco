package filehandling

import "sacco/server/parser"

func HandleNestedModel(data any, phoneNumber, cacheFolder *string,
	saveFunc func(map[string]any, string, int) (*int64, error), sessions map[string]*parser.Session, sessionFolder string) error {
	if modelData, ok := data.(map[string]any); ok {
		var id int64

		if phoneNumber != nil && *phoneNumber != "default" {
			if modelData["phoneNumber"] == nil {
				modelData["phoneNumber"] = *phoneNumber
			}
		}

		_ = id
	}

	return nil
}
