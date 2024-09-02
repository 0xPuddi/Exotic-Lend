package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/utils"
)

// Takes multiple rows to add to the table, it adds them
// respecting the order
//
// Parameters:
//   - db:		the database struct
//   - data:	the structs rows
//
// Returns:
//   - []sql.Result:	an array of results of the same length as entries
//   - []error:			an array of errors of the queries
func InsertEntries[T any](db *sql.DB, data []T) ([]sql.Result, []error) {
	var results []sql.Result
	var errors []error

	for _, d := range data {
		result, err := InsertEntry(db, d)

		results = append(results, result)
		errors = append(errors, err)
	}

	return results, errors
}

// Takes multiple queries, it checks and performs them
//
// Parameters:
//   - db:		the database struct
//   - queries:	the queries strings
//
// Returns:
//   - sql.Result:	an array of results of the same length as entries
//   - error:		if an error occured during the process
func InsertEntry(db *sql.DB, data any) (sql.Result, error) {
	// Check if it is a struct
	ty := reflect.TypeOf(data)
	if !utils.ValidateStruct(ty) {
		return nil, fmt.Errorf("data format is wrong")
	}

	// Create table if it doesn't exist
	exists, err := CheckIfTableExists(db, data)
	if err != nil && err != ErrTableExists {
		return nil, err
	}
	if err != ErrTableExists && !exists {
		CreateTable(db, data)
	}

	// Build SQL query
	query, err := ParseStructToEntry(ty, reflect.ValueOf(data))
	if err != nil {
		return nil, err
	}

	// Make Query
	return MakeQuery(db, query)
}

// Parses the value to the correct SQL formatting based on type
// Correctly supports: string, bool, integers, unsigned integers
// Doesn't support: runes (retunred as digit)
// Other types will be converted to string through their interface
//
// Parameters:
//   - value:	the struct to be inserted in the database
//
// Returns:
//   - string:	the formatted SQL value
//   - error:	if any error occured during parsing
func ParseValueToEntry(value reflect.Value) string {
	// String
	if value.Kind() == reflect.String {
		return fmt.Sprintf(`'%s'`, value.String())
	}

	// Bool
	if value.Kind() == reflect.Bool {
		if value.Bool() {
			return "TRUE"
		} else {
			return "FALSE"
		}
	}

	return fmt.Sprint(value.Interface())
}

// Parses the struct into a insertion query with its current parameters
//
// Parameters:
//   - data:	the struct to be inserted in the database
//
// Returns:
//   - string:	the insertion query
//   - error:	if any error occured during parsing
func ParseStructToEntry(data reflect.Type, value reflect.Value) (string, error) {
	// Check it is a struct
	if !utils.ValidateStruct(data) {
		return "", fmt.Errorf("cannot parse non-struct into insertion query: %v", data)
	}

	// Build String
	// Create query string
	var builder strings.Builder
	builder.WriteString("INSERT INTO ")
	builder.WriteString(data.Name())
	builder.WriteString(" (")

	// Add columns
	numFields := data.NumField()
	for i := 0; i < numFields; i++ {
		f := data.Field(i)

		// Check if there is a nested struct other than types.Default
		if utils.ValidateStruct(f.Type) && !utils.ValidateCustomStruct(f.Type) {
			return "", fmt.Errorf("cannot have nested not custom struct as tables: %v", f)
		}

		name_db, err_db := utils.GetFieldNameDB(f)
		if err_db != nil {
			return "", err_db
		}

		if i == numFields-1 {
			builder.WriteString(name_db)
			continue
		}

		builder.WriteString(name_db)
		builder.WriteString(", ")
	}

	builder.WriteString(")\n")
	builder.WriteString("VALUES (")

	// Add Values
	for i := 0; i < numFields; i++ {
		f := data.Field(i)
		v := value.Field(i)

		// Check if there is a nested struct other than types.Default
		if utils.ValidateStruct(f.Type) {
			val, err := ParseCustomStruct(f.Type, v)
			if err != nil {
				return "", err
			}

			builder.WriteString(val)

			if i == numFields-1 {
				continue
			}

			builder.WriteString(", ")
			continue
		}

		// Write value
		builder.WriteString(ParseValueToEntry(v))

		if i == numFields-1 {
			continue
		}
		builder.WriteString(", ")
	}

	builder.WriteString(")")
	return builder.String(), nil
}

// Parses the custom struct into the correct insertion query parameter
// Supported structs: Default, Null
//
// Parameters:
//   - t:	the reflect.Type of the struct
//   - v:	the reflect.Value of the struct
//
// Returns:
//   - string:	the insertion query parameter
//   - error:	if any error occured during parsing
func ParseCustomStruct(t reflect.Type, v reflect.Value) (string, error) {
	if utils.ValidateDefaultStruct(t) {
		defaultField := v.FieldByName("Default")
		valueField := v.FieldByName("Value")

		if defaultField.Bool() {
			// Write default
			return "DEFAULT", nil
		} else {
			// Write value
			return fmt.Sprint(valueField.Interface()), nil
		}
	}

	if utils.ValidateNullStruct(t) {
		nullField := v.FieldByName("Null")
		valueField := v.FieldByName("Value")

		if nullField.Bool() {
			// Write default
			return "NULL", nil
		} else {
			// Write value
			return fmt.Sprint(valueField.Interface()), nil
		}
	}

	if utils.ValidateTimestampStruct(t) {
		nowField := v.FieldByName("Now")
		unixField := v.FieldByName("Unix")

		if nowField.Bool() {
			// Write default
			return "NOW()", nil
		} else {
			// Write value
			return fmt.Sprintf("TO_TIMESTAMP(%d)", unixField.Interface()), nil
		}
	}

	return "", fmt.Errorf("cannot have nested struct as tables: %v", v.Interface())
}
