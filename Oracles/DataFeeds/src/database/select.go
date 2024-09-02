package database

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/utils"
)

var (
	ErrComlumnIndexOutOfBounds = errors.New("column index out of bounds")
)

// Selects all rows from the table
//
// Parameters:
//   - db:		the database struct
//   - table:	the table truct
//
// Returns:
//   - *sql.Rows:	query rows result
//   - error:		error if occured
func SelectAll(db *sql.DB, table any) (*sql.Rows, error) {
	tt := reflect.TypeOf(table)

	if !utils.ValidateStruct(tt) {
		return nil, fmt.Errorf("table is not a struct")
	}

	query := "SELECT * FROM " + utils.BaseTypeName(tt)

	return MakeQueryWithResult(db, query)
}

// Makes a selection query of specific columns of the table
//
// Parameters:
//   - db:		the database struct
//   - table:	the table struct
//
// Returns:
//   - *sql.Rows:	query rows result
//   - error:		error if occured
func SelectColumns(db *sql.DB, table any, columns ...int) (*sql.Rows, error) {
	query, err := parseStructToSelectColumns(reflect.TypeOf(table), columns...)
	if err != nil {
		return nil, err
	}

	return MakeQueryWithResult(db, query)
}

// Parses the struct into a column selection query
//
// Parameters:
//   - tt:		the reflect table type
//   - columns:	the columns to include in the query
//
// Returns:
//   - string:	the query
//   - error:	error if occured
func parseStructToSelectColumns(tt reflect.Type, columns ...int) (string, error) {
	if !utils.ValidateStruct(tt) {
		return "", fmt.Errorf("table is not a struct")
	}

	if len(columns) > tt.NumField() || len(columns) == 0 {
		return "", fmt.Errorf("selected more than actual columns")
	}

	var builder strings.Builder
	builder.WriteString("SELECT ")
	for i, c := range columns {
		if c >= tt.NumField() {
			return "", fmt.Errorf("field number out of range")
		}

		str_db := tt.Field(c).Tag.Get("db")
		name_col := strings.Split(str_db, " ")[0]

		builder.WriteString(name_col)

		if i+1 == len(columns) {
			continue
		}
		builder.WriteString(", ")
	}

	builder.WriteString("\nFROM ")
	builder.WriteString(utils.BaseTypeName(tt))
	return builder.String(), nil
}

// Makes a query selecting all columns with many custom conditions
//
// Parameters:
//   - db:			database sql driver
//   - table:		the table struct
//   - conditions:	the conditions to include in the query
//
// Returns:
//   - *sql.Rows:	queried rows
//   - error:		error if occured
func SelectAllConditions(db *sql.DB, table any, conditions ...string) (*sql.Rows, error) {
	query, err := buildSelectConditionsQuery(reflect.TypeOf(table), []int{}, conditions...)
	if err != nil {
		return nil, err
	}

	return MakeQueryWithResult(db, query)
}

// Makes a query selecting defined columns with many custom conditions
//
// Parameters:
//   - db:				database sql driver
//   - table:			the table struct
//   - selectColumns:	the columns to be selected
//   - conditions:		the conditions to include in the query
//
// Returns:
//   - *sql.Rows:	queried rows
//   - error:		error if occured
func SelectColumnsConditions(db *sql.DB, table any, selectColumns []int, conditions ...string) (*sql.Rows, error) {
	query, err := buildSelectConditionsQuery(reflect.TypeOf(table), selectColumns, conditions...)
	if err != nil {
		return nil, err
	}

	return MakeQueryWithResult(db, query)
}

// Builds a query with many conditions
//
// Parameters:
//   - tt:				the reflect table type
//   - selectColumns:	the columns to select in the query
//   - conditions:		the conditions to include in the query
//
// Returns:
//   - string:	the query
//   - error:	error if occured
func buildSelectConditionsQuery(tt reflect.Type, selectColumns []int, conditions ...string) (string, error) {
	if !utils.ValidateStruct(tt) {
		return "", fmt.Errorf("table is not a struct")
	}
	var builder strings.Builder

	builder.WriteString("SELECT ")
	if len(selectColumns) == 0 {
		builder.WriteString("*")
	} else {
		for i, sc := range selectColumns {
			if sc >= tt.NumField() {
				return "", ErrComlumnIndexOutOfBounds
			}

			f := tt.Field(sc)

			fieldName, err := utils.GetFieldNameDB(f)
			if err != nil {
				return "", err
			}

			builder.WriteString(fieldName)

			if i == len(selectColumns)-1 {
				continue
			}

			builder.WriteString(", ")
		}
	}
	builder.WriteString(" FROM ")
	builder.WriteString(utils.BaseTypeName(tt))

	for _, s := range conditions {
		builder.WriteString("\n")
		builder.WriteString(s)
	}

	return builder.String(), nil
}

// Makes a ordered query on the Price struct based on asset_id filtering
//
// Parameters:
//   - table:			the struct table
//   - selectColumns:	the columns to be selected in the query (empty to select all), note that order in indexes will be followed in the query result
//   - asset_id:		the asset_id to filter the query by
//   - orderByColumn:	the column to order by
//   - limit:			the maximum number of elements to retrive (>0)
//   - desc:			if you wish to sort descending or ascending
//
// Returns:
//   - *sql.Rows:	the rows result
//   - error:		error if occured
func SelectAllWhereAssetIdOrderedRow(db *sql.DB, table any, assetId int, orderByColumn int, limit int, desc bool) (*sql.Rows, error) {
	// Build Order
	tt := reflect.TypeOf(table)
	if orderByColumn >= tt.NumField() {
		return nil, ErrComlumnIndexOutOfBounds
	}

	conditions, err := buildSelectWhereAssetIdOrderedRowConditions(tt, assetId, orderByColumn, limit, desc)
	if err != nil {
		return nil, err
	}

	return SelectAllConditions(db, table, conditions...)
}

// Builds conditions for a most recent query on the Price struct based on asset_id
//
// Parameters:
//   - tt:				the reflect table type
//   - selectColumns:	the columns to be selected in the query (empty to select all), note that order in indexes will be followed in the query result
//   - asset_id:		the asset_id to filter the query by
//   - orderByColumn:	the column to order by
//   - limit:			the maximum number of elements to retrive (>0)
//   - desc:			if you wish to sort descending or ascending
//
// Returns:
//   - *sql.Rows:	the rows result
//   - error:		error if occured
func buildSelectWhereAssetIdOrderedRowConditions(tt reflect.Type, asset_id int, orderByColumn int, limit int, desc bool) ([]string, error) {
	var conditions []string

	dbColumnName, err := utils.GetFieldNameDB(tt.Field(orderByColumn))
	if err != nil {
		return nil, err
	}

	conditions = append(conditions, fmt.Sprintf("WHERE asset_id = %d", asset_id))

	if desc {
		conditions = append(conditions, fmt.Sprintf("ORDER BY %s DESC", dbColumnName))
	} else {
		conditions = append(conditions, fmt.Sprintf("ORDER BY %s ASC", dbColumnName))
	}

	// Insert limit
	if limit >= 0 {
		conditions = append(conditions, fmt.Sprintf("LIMIT %d", limit))
	}

	return conditions, nil
}

// Makes a query on assets based on the source
//
// Parameters:
//   - db:			the databse driver
//   - table:	the
//   - asset_id:		the asset_id to filter the query by
//   - orderByColumn:	the column to order by
//   - limit:			the maximum number of elements to retrive (>0)
//   - desc:			if you wish to sort descending or ascending
//
// Returns:
//   - *sql.Rows:	the rows result
//   - error:		error if occured
func SelectTableByMatchRow(db *sql.DB, table any, matchColumns []int, matchValues []any, limit int) (*sql.Rows, error) {
	// Build Order
	tt := reflect.TypeOf(table)
	if len(matchColumns) >= tt.NumField() || len(matchColumns) != len(matchValues) {
		return nil, ErrComlumnIndexOutOfBounds
	}

	conditions, err := buildSelectTableByMatch(tt, matchColumns, matchValues, limit)
	if err != nil {
		return nil, err
	}

	return SelectColumnsConditions(db, table, []int{}, conditions...)
}

// Makes the select query where you can select specific values for each column
//
// Parameters:
//   - db:				the databse driver
//   - table:			the struct table
//   - selectColumns:	the columns to be selected
//   - matchColumns:	the columns to be matched with a value
//   - matchValues:		values that are going to be matched with
//   - limit:			the limit of columns to return, negative to return them all
//
// Returns:
//   - *sql.Rows:	the rows result
//   - error:		error if occured
func SelectTableByMatchColumns(db *sql.DB, table any, selectColumns []int, matchColumns []int, matchValues []any, limit int) (*sql.Rows, error) {
	// Build Order
	tt := reflect.TypeOf(table)
	if len(matchColumns) >= tt.NumField() || len(selectColumns) >= tt.NumField() || len(matchColumns) != len(matchValues) {
		return nil, ErrComlumnIndexOutOfBounds
	}

	conditions, err := buildSelectTableByMatch(tt, matchColumns, matchValues, limit)
	if err != nil {
		return nil, err
	}

	return SelectColumnsConditions(db, table, selectColumns, conditions...)
}

// Builds the select query where you can select specific values for each column
//
// Parameters:
//   - tt:				the struct table reflect type
//   - matchColumns:	the columns to be matched with a value
//   - matchValues:		values that are going to be matched with
//   - limit:			the limit of columns to return, negative to return them all
//
// Returns:
//   - *sql.Rows:	the rows result
//   - error:		error if occured
func buildSelectTableByMatch(tt reflect.Type, matchColumns []int, matchValues []any, limit int) ([]string, error) {
	var conditions []string
	var builder strings.Builder

	builder.WriteString("WHERE ")
	mvv := reflect.ValueOf(matchValues)
	for i := 0; i < len(matchColumns); i++ {
		match_column, err := utils.GetFieldNameDB(tt.Field(matchColumns[i]))
		if err != nil {
			return []string{}, err
		}
		match_value := ParseValueToEntry(mvv.Index(i).Elem())

		if i == 0 {
			builder.WriteString(match_column)
			builder.WriteString(" = ")
			builder.WriteString(match_value)
			conditions = append(conditions, builder.String())
			builder.Reset()
			continue
		}

		builder.WriteString("AND ")
		builder.WriteString(match_column)
		builder.WriteString(" = ")
		builder.WriteString(match_value)
		conditions = append(conditions, builder.String())
		builder.Reset()
	}

	// Insert limit
	if limit >= 0 {
		conditions = append(conditions, fmt.Sprintf("LIMIT %d", limit))
	}

	return conditions, nil
}

// func parseCustomStructValue(cv reflect.Value) {

// }
