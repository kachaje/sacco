package models

import "database/sql"

type Model struct {
	Fields     []string
	FieldTypes map[string]string

	db *sql.DB
}

func NewModel(db *sql.DB, fields []string, fieldTypes map[string]string) *Model {
	m := &Model{
		db:         db,
		Fields:     fields,
		FieldTypes: fieldTypes,
	}

	return m
}
