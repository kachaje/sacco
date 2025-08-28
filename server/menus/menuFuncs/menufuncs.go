package menufuncs

import (
	"sacco/server/database"
	"sacco/server/parser"
)

var (
	DB       *database.Database
	Sessions = map[string]*parser.Session{}
)
