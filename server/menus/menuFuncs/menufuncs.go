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
	FunctionsMap["bankingDetails"] = BankingDetails
	FunctionsMap["blockUser"] = BlockUser
	FunctionsMap["businessSummary"] = BusinessSummary
	FunctionsMap["changePassword"] = ChangePassword
	FunctionsMap["checkBalance"] = CheckBalance
	FunctionsMap["devConsole"] = DevConsole
	FunctionsMap["doExit"] = DoExit
	FunctionsMap["editUser"] = EditUser
	FunctionsMap["employmentSummary"] = EmploymentSummary
	FunctionsMap["landing"] = Landing
	FunctionsMap["listUsers"] = ListUsers
	FunctionsMap["memberLoansSummary"] = MemberLoansSummary
	FunctionsMap["signIn"] = SignIn
	FunctionsMap["signUp"] = SignUp
	FunctionsMap["viewMemberDetails"] = ViewMemberDetails
}
