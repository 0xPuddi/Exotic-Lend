package database

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/types"
	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/utils"
)

var (
	ErrNotValidTable = errors.New("not a valid table passed in the function")
)

// Takes a table and a SELECT query, using that query it then eliminates all
// results based on their unique primary key
//
// Parameters:
//   - db:			the database driver
//   - table:		the struct table
//   - selectQuery:	the SELECT query
//
// Returns:
//   - sql.Result:	an array of results of the same length as entries
//   - error:		if an error occured during the process
func DeleteRowsByPrimaryKeyWithSelectionQuery(db *sql.DB, table types.Table, selectQuery string) (sql.Result, error) {
	tt := reflect.TypeOf(table)
	if utils.ValidateDefaultStruct(tt) {
		return nil, ErrNotValidTable
	}

	table_name := utils.BaseTypeName(tt)
	var builder strings.Builder
	builder.WriteString("DELETE FROM ")
	builder.WriteString(table_name)
	builder.WriteString("\n")

	table_id, err := table.GetPrimaryKeyNameDB()
	if err != nil {
		return nil, err
	}

	builder.WriteString("WHERE ")
	builder.WriteString(table_id)
	builder.WriteString(" IN (\n")
	builder.WriteString(selectQuery)
	builder.WriteString("\n);")

	return MakeQuery(db, builder.String())
}
