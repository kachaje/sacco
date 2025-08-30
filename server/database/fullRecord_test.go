package database_test

import (
	"sacco/server/database"
	"testing"
)

func TestFullRecord(t *testing.T) {
	dbname := ":memory:"
	db := database.NewDatabase(dbname)
	defer db.Close()

}
