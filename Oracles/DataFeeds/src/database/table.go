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
	ErrTableExists = errors.New("table already exists")
)

// Takes a db driver and a data struct, it then creates the struct table
// in the db
//
// Parameters:
//   - db:		the database struct
//   - data:	the data struct
//
// Returns:
//   - QueriesResult:	an array of results of the same length as entries
//   - error:			if any error occured
func CreateTable(db *sql.DB, data any) (sql.Result, error) {
	// Check struct
	if !utils.ValidateStruct(reflect.TypeOf(data)) {
		return nil, fmt.Errorf("data is not a struct")
	}

	// Check existence
	exists, err := CheckIfTableExists(db, data)
	if exists {
		return nil, ErrTableExists
	} else if err != nil {
		return nil, err
	}

	// Create query
	query, err := ParseStructToTable(reflect.TypeOf(data))
	if err != nil {
		return nil, err
	}

	// Query
	return MakeQuery(db, query)
}

// Takes a db driver and a data struct, it makes a query to
// see if the struct table exists or not
//
// Parameters:
//   - db:		the database struct
//   - data:	the data struct
//
// Returns:
//   - bool:	if the table exists or not
//   - error:	if any error occured
func CheckIfTableExists[T any](db *sql.DB, data T) (bool, error) {
	d := reflect.TypeOf(data)
	name := strings.ToLower(utils.BaseTypeName(d))

	query := fmt.Sprintf(`SELECT EXISTS (
		SELECT 1 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = '%s'
	);`, name)

	queryRows, err := MakeQueryWithResult(db, query)

	if err != nil {
		return false, err
	}
	defer queryRows.Close()

	// Extract value
	var exists bool

	if queryRows.Next() {
		if err := queryRows.Scan(&exists); err != nil {
			return false, err
		}
	} else {
		// Handle the case where no rows are returned, which should not happen for EXISTS
		return false, fmt.Errorf("no rows returned")
	}

	return exists, nil
}

// Parses the struct into a table and returns its creation query
//
// Parameters:
//   - data:	the struct to be created as table
//
// Returns:
//   - string:	the table creation query
//   - error:	if any erro occured during parsing
func ParseStructToTable(data reflect.Type) (string, error) {
	var idx []string // idx
	var ref []string // references

	// Check it is a struct
	if !utils.ValidateStruct(data) {
		return "", fmt.Errorf("cannot parse non-struct into table: %v", data)
	}

	// Build String
	var builder strings.Builder
	builder.WriteString("CREATE TABLE IF NOT EXISTS ")
	builder.WriteString(data.Name())
	builder.WriteString(" (\n")

	// Parse each field into query
	numField := data.NumField()
	for i := 0; i < numField; i++ {
		f := data.Field(i)

		// Check if there is a nested struct other than types.Default
		if utils.ValidateStruct(f.Type) && !utils.ValidateCustomStruct(f.Type) {
			return "", fmt.Errorf("cannot have nested struct as tables: %v", f)
		}

		// Add db type
		str_db, ok_db := f.Tag.Lookup("db")
		if ok_db {
			if i == numField-1 && len(ref) == 0 {
				// If it is the last skip the triling comma
				builder.WriteString("\t")
				builder.WriteString(str_db)
			} else {
				builder.WriteString("\t")
				builder.WriteString(str_db)
				builder.WriteString(",")
			}
		} else {
			return "", fmt.Errorf("column is not defined")
		}

		// Add any references
		str_ref, ok_ref := f.Tag.Lookup("ref")
		if ok_ref {
			ref = append(ref, str_ref)

			if i == numField-1 {
				builder.WriteString(",")
			}
		}

		// Add any index
		str_idx, ok_idx := f.Tag.Lookup("idx")
		if ok_idx {
			idx = append(idx, str_idx)
		}

		// Close the line
		builder.WriteString("\n")
	}

	// Add any reference at the end of the table
	numRef := len(ref)
	for i, r := range ref {
		if i == numRef-1 {
			builder.WriteString("\t")
			builder.WriteString(r)
			builder.WriteString("\n")
			continue
		}

		builder.WriteString("\t")
		builder.WriteString(r)
		builder.WriteString(",\n")
	}
	builder.WriteString(");")

	// Add any indexes at the end of the query
	numIdx := len(idx)
	for i, id := range idx {
		if i == 0 {
			builder.WriteString("\n")
		}

		if i == numIdx-1 {
			builder.WriteString(id)
			builder.WriteString(";")
			continue
		}

		builder.WriteString(id)
		builder.WriteString(";\n")
	}

	// Return the query
	return builder.String(), nil
}
