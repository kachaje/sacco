package database

import "fmt"

func (d *Database) FullRecord(phoneNumber string) (map[string]any, error) {
	var data = map[string]any{}

	results, err := d.GenericModels["member"].FilterBy(fmt.Sprintf(`WHERE phoneNumber = "%s" AND active = 1`, phoneNumber))
	if err != nil {
		return nil, err
	}

	if len(results) > 0 {
		data = results[0]
	} else {
		return nil, fmt.Errorf("no match found")
	}

	return data, nil
}
