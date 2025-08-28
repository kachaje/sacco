package menufuncs

import (
	"sacco/server/database"
	"sacco/server/parser"
)

var (
	DB       *database.Database
	Sessions = map[string]*parser.Session{}

	FunctionsMap = map[string]func(
		func(
			string, *parser.Session,
			string, string, string, string,
		) string,
		map[string]any,
	) string{}
)

func init() {
	FunctionsMap["doExit"] = DoExit

	FunctionsMap["businessSummary"] = BusinessSummary

	FunctionsMap["employmentSummary"] = EmploymentSummary

	FunctionsMap["checkBalance"] = CheckBalance

	FunctionsMap["bankingDetails"] = BankingDetails

	FunctionsMap["viewMemberDetails"] = ViewMemberDetails

	FunctionsMap["devConsole"] = DevConsole

	FunctionsMap["memberLoansSummary"] = MemberLoansSummary

	FunctionsMap["signIn"] = SignIn

	FunctionsMap["listUsers"] = ListUsers

	FunctionsMap["blockUser"] = BlockUser

	FunctionsMap["editUser"] = EditUser

	FunctionsMap["changePassword"] = ChangePassword

	FunctionsMap["signUp"] = SignUp

	FunctionsMap["landing"] = Landing
}
