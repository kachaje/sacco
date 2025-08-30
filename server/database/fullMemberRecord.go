package database

import (
	"fmt"
	"log"
	"sacco/utils"
)

func (d *Database) LoadModelChildren(model string, id int64) (map[string]any, error) {
	data, err := d.GenericModels[model].FetchById(id)
	if err != nil {
		return nil, err
	}

	if len(data) <= 0 {
		return nil, fmt.Errorf("no match found")
	}

	capModel := utils.CapitalizeFirstLetter(model)

	if arrayChidren, ok := ArrayChildren[fmt.Sprintf("%sArrayChildren", capModel)]; ok {
		for _, childModel := range arrayChidren {
			parentKey := fmt.Sprintf("%sId", model)

			results, err := d.GenericModels[childModel].FilterBy(fmt.Sprintf(`WHERE %s = %v AND active = 1`, parentKey, id))
			if err != nil {
				log.Println(err)
				continue
			}

			if len(results) > 0 {
				data[childModel] = results
			}
		}
	}
	if singleChidren, ok := SingleChildren[fmt.Sprintf("%sSingleChildren", capModel)]; ok {
		for _, childModel := range singleChidren {
			parentKey := fmt.Sprintf("%sId", model)

			results, err := d.GenericModels[childModel].FilterBy(fmt.Sprintf(`WHERE %s = %v AND active = 1 ORDER by updated_at DESC LIMIT 1`, parentKey, id))
			if err != nil {
				log.Println(err)
				continue
			}

			if len(results) > 0 {
				data[childModel] = results[0]
			}
		}
	}

	return data, nil
}

func (d *Database) FullMemberRecord(phoneNumber string) (map[string]any, error) {
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
