package filehandling

import (
	"fmt"
	"os"
	"path/filepath"
	"sacco/server/parser"
	"strconv"
)

func HandleCommonModels(data any, model, phoneNumber, sessionId, cacheFolder *string,
	saveFunc func(
		data map[string]any,
		model string,
		parentId *int64,
	) (*int64, error), sessions map[string]*parser.Session, sessionFolder string) error {
	val, ok := data.(map[string]any)
	if ok {
		filename := filepath.Join(sessionFolder, fmt.Sprintf("%s.json", *model))

		transactionDone := false

		// By default cache the data first in case we lose database connection
		CacheFile(filename, val, 0)
		defer func() {
			if transactionDone {
				os.Remove(filename)
			}
		}()

		for _, key := range []string{
			"netPay", "grossPay", "periodEmployed", "yearsInBusiness",
			"totalIncome", "totalCostOfGoods", "employeesWages", "ownSalary",
			"transport", "loanInterest", "utilities", "rentals", "otherCosts",
			"totalCosts", "netProfitLoss", "numberOfShares", "pricePerShare",
			"loanAmount", "repaymentPeriod", "amountRecommended",
			"amountApproved", "value",
		} {
			if val[key] != nil {
				nv, ok := val[key].(string)
				if ok {
					real, err := strconv.ParseFloat(nv, 64)
					if err == nil {
						val[key] = real
					}
				}
			}
		}

		if sessions[*sessionId].MemberId != nil {
			val["memberId"] = *sessions[*sessionId].MemberId

			if saveFunc == nil {
				return fmt.Errorf("server.SaveData.%s.1:missing saveFunc", *model)
			}

			_, err := saveFunc(val, *model, sessions[*sessionId].MemberId)
			if err != nil {
				return fmt.Errorf("server.SaveData.%s.2:%s", *model, err.Error())
			}

			transactionDone = true
		}

		sessions[*sessionId].ActiveMemberData[*model] = val

		sessions[*sessionId].AddedModels[*model] = true

		sessions[*sessionId].RefreshSession()

		sessions[*sessionId].LoadMemberCache(*phoneNumber, *cacheFolder)
	}

	return nil
}
