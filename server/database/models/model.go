package models

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
)

type Model struct {
	ModelName  string
	Fields     []string
	FieldTypes map[string]string

	db *sql.DB
}

func NewModel(
	db *sql.DB,
	modelName string,
	fields []string,
	fieldTypes map[string]string,
) (*Model, error) {
	if modelName == "" {
		return nil, fmt.Errorf("missing required modelName")
	}

	m := &Model{
		db:         db,
		ModelName:  modelName,
		Fields:     fields,
		FieldTypes: fieldTypes,
	}

	return m, nil
}

func (m *Model) AddRecord(data map[string]any) (*int64, error) {
	fields := []string{}
	values := []any{}
	markers := []string{}
	var id int64

	for key, value := range data {
		if !slices.Contains(m.Fields, key) {
			continue
		}

		fields = append(fields, key)
		markers = append(markers, "?")
		values = append(values, value)
	}

	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`,
		m.ModelName,
		strings.Join(fields, ", "),
		strings.Join(markers, ", "),
	)

	result, err := QueryWithRetry(m.db, context.Background(), query, values...)
	if err != nil {
		return nil, err
	}

	if id, err = result.LastInsertId(); err != nil {
		return nil, err
	}

	return &id, nil
}
