package generation

import (
	"database/sql"
	"fmt"
	"log"
	"sort"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/ch3nnn/sql2pb/cmd/generation/parser"
	"github.com/ch3nnn/sql2pb/cmd/generation/tools/stringx"
)

func generateSchema(table string, ignoreTables, ignoreColumns []string, serviceName, fieldStyle, dbType string) (*parser.Schema, error) {
	db, err := db()
	if err != nil {
		return nil, err
	}

	dbs, err := dbSchema(db)
	if nil != err {
		return nil, err
	}

	cols, err := dbColumns(db, dbs, table, dbType)
	if nil != err {
		return nil, err
	}

	schema := parser.NewSchema("proto3", serviceName, goPackageName, packageName)
	if err := schema.TypesFromColumns(cols, ignoreTables, ignoreColumns, fieldStyle); nil != err {
		return nil, err
	}

	sort.Sort(schema.Imports)
	//sort.Sort(schema.Messages)
	//sort.Sort(schema.Enums)

	return schema, nil
}

func db() (db *sql.DB, err error) {
	var dataSourceName string

	switch dbType {
	case "mysql":
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname)
	case "postgres":
		dataSourceName = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	default:
		log.Fatal("dbType not supported")
	}

	return sql.Open(dbType, dataSourceName)

}
func dbSchema(db *sql.DB) (string, error) {
	var schema, query string
	switch dbType {
	case "mysql":
		query = `SELECT SCHEMA()`
	case "postgres":
		query = `SELECT CURRENT_DATABASE()`
	default:
		log.Fatal("dbType not supported")
	}

	if err := db.QueryRow(query).Scan(&schema); err != nil {
		return "", err
	}

	return schema, nil
}

func dbColumns(db *sql.DB, dbs, table, dbType string) ([]parser.Column, error) {
	rows, err := db.Query(querySQL(dbs, dbType, table))
	if nil != err {
		return nil, err
	}
	defer rows.Close()

	var cols []parser.Column
	for rows.Next() {
		var cs parser.Column
		err := rows.Scan(
			&cs.TableName,
			&cs.ColumnName,
			&cs.IsNullable,
			&cs.DataType,
			&cs.CharacterMaximumLength,
			&cs.NumericPrecision,
			&cs.NumericScale,
			&cs.ColumnType,
			&cs.ColumnComment,
			&cs.TableComment,
		)
		if err != nil {
			log.Fatal(errors.Wrapf(err, "scan error, table: %s, column: %s", cs.TableName, cs.ColumnName))
		}

		if cs.TableComment == "" {
			cs.TableComment = stringx.From(cs.TableName).ToCamelWithStartLower()
		}

		cols = append(cols, cs)
	}
	if err := rows.Err(); nil != err {
		return nil, err
	}

	return cols, nil
}

func querySQL(dbs, dbType, table string) (sql string) {
	switch dbType {
	case "mysql":
		sql = `SELECT
					c.TABLE_NAME,
					c.COLUMN_NAME,
					c.IS_NULLABLE,
					c.DATA_TYPE,
					c.CHARACTER_MAXIMUM_LENGTH,
					c.NUMERIC_PRECISION,
					c.NUMERIC_SCALE,
					c.COLUMN_TYPE ,
					c.COLUMN_COMMENT,
					t.TABLE_COMMENT
				FROM
					INFORMATION_SCHEMA.COLUMNS AS c
				LEFT JOIN INFORMATION_SCHEMA.TABLES AS t ON
					c.TABLE_NAME = t.TABLE_NAME
					AND c.TABLE_SCHEMA = t.TABLE_SCHEMA
				WHERE 
					c.TABLE_SCHEMA = '%s'
					AND c.TABLE_NAME = '%s' 
				ORDER BY 
					c.TABLE_NAME,
					c.ORDINAL_POSITION`
		sql = fmt.Sprintf(sql, dbs, table)
	case "postgres":
		sql = `SELECT
					col.table_name AS TABLE_NAME,  -- 表名
					col.column_name AS COLUMN_NAME, -- 字段名
					col.is_nullable AS IS_NULLABLE, -- 是否为 null
					-- col.data_type as DATA_TYPE,
					col.udt_name AS DATA_TYPE,
					col.character_maximum_length AS CHARACTER_MAXIMUM_LENGTH , -- 字符最大长度
					col.numeric_precision AS NUMERIC_PRECISION, -- 数值精度
					col.numeric_scale AS NUMERIC_SCALE , -- 小数点后的精度基本单位的数
					 col.udt_name AS COLUMN_TYPE,  -- 字段类型
					COALESCE(pd.description, '') AS COLUMN_COMMENT, -- 字段注释
					COALESCE(OBJ_DESCRIPTION('%s'::regclass, 'pg_class'), '') AS TABLE_COMMENT -- 表注释
				FROM
					information_schema.columns AS col
				LEFT JOIN
					pg_description AS pd
				ON
					col.table_name::regclass = pd.objoid
					AND col.ordinal_position = pd.objsubid
				WHERE
					col.table_name ='%s' 
					AND 
					col.table_catalog = '%s'  -- 数据库名称
				ORDER BY col.table_name, col.ORDINAL_POSITION`
		sql = fmt.Sprintf(sql, table, table, dbs)
	default:
		log.Fatal("dbType not supported")
	}

	return
}
