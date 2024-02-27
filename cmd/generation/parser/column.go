package parser

import "database/sql"

type Column struct {
	Style                  string
	TableName              string
	TableComment           string
	ColumnName             string
	IsNullable             string
	DataType               string
	CharacterMaximumLength sql.NullInt64
	NumericPrecision       sql.NullInt64
	NumericScale           sql.NullInt64
	ColumnType             string
	ColumnComment          string
}
