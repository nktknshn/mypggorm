package mypggorm

import (
	"fmt"

	"gorm.io/gorm"
)

type ColumnInfo struct {
	Name string
	Type string
}

func (c ColumnInfo) String() string {
	return fmt.Sprintf("%v:`%v`", c.Name, c.Type)
}

type DatabaseSchema map[string][]ColumnInfo

func (s DatabaseSchema) String() string {
	var result string

	for tableName, columns := range s {
		result += fmt.Sprintf("%v:\n", tableName)
		for _, column := range columns {
			result += fmt.Sprintf("\t%v\n", column.String())
		}
	}

	return result
}

func GetSchemaPublic(db *gorm.DB) (DatabaseSchema, error) {
	return GetSchema(db, "public")
}

func GetSchema(db *gorm.DB, tableSchema string) (DatabaseSchema, error) {
	var schema DatabaseSchema = make(map[string][]ColumnInfo)

	rows, err := db.Raw("SELECT table_name, column_name, data_type FROM information_schema.columns WHERE table_schema=?", tableSchema).Rows()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var tableName string
		var columnInfo ColumnInfo

		err = rows.Scan(&tableName, &columnInfo.Name, &columnInfo.Type)
		if err != nil {
			return nil, err
		}

		schema[tableName] = append(schema[tableName], columnInfo)
	}

	return schema, nil
}
