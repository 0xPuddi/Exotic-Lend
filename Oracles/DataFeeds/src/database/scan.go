package database

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/utils"
)

type ScannableRow interface {
	Scan(dest ...any) error
}

// ScanRowToStruct scans all rows of the selected table rows
//
// Parameters:
//
//   - *sql.Row:	the row to scan
//
//   - table:		the table to scan
//
//     note table has to be passed as: reflect.ValueOf(&yourStruct).Elem() I have to idea why
//
//     this is the only way to preserve the reflect.Struct type attribute while having an
//
//     addressable reflect.Type
//
// Returns:
//   - error: an error if occurs, nil otherwise
func ScanRowToStruct(row ScannableRow, table reflect.Value) error {
	if !utils.ValidateStruct(table.Type()) {
		return fmt.Errorf("cant scan a row without a struct table")
	}

	var addresses []any
	for i := 0; i < table.NumField(); i++ {
		tf := table.Field(i)

		if utils.ValidateCustomStruct(tf.Type()) {
			tf, err := handleCustomStructField(tf)
			if err != nil {
				return err
			}

			if !tf.CanAddr() {
				return fmt.Errorf("field is not addressable %v", tf)
			}
			addresses = append(addresses, tf.Addr().Interface())
			continue
		}

		if !tf.CanAddr() {
			return fmt.Errorf("field is not addressable %v", tf)
		}
		addresses = append(addresses, tf.Addr().Interface())
	}

	return row.Scan(addresses...)
}

// ScanSelectedRowsToParameters scans the selected table rows
//
// Parameters:
//   - *sql.Row:	the row to scan
//   - table:		the table to scan
//   - ...int:		the index of table rows to scan over the row. Indexes MUST be ordered with the queried rows
//
// Returns:
//   - error: an error if occurs, nil otherwise
func ScanSelectedRowsToParameters(row *sql.Row, table reflect.Value, tableRows ...int) error {
	if table.NumField() > len(tableRows) {
		return fmt.Errorf("more rows to scan than actual rows")
	}

	var addresses []any
	for _, indexRow := range tableRows {
		if table.NumField() >= 1+indexRow {
			return fmt.Errorf("more rows to scan than actual rows")
		}

		tf := table.Field(indexRow)

		if utils.ValidateCustomStruct(tf.Type()) {
			tf, err := handleCustomStructField(tf)
			if err != nil {
				return err
			}

			addresses = append(addresses, tf.Addr().Interface())
			continue
		}

		if !tf.CanAddr() {
			return fmt.Errorf("field is not addressable %v", tf)
		}
		addresses = append(addresses, tf.Addr().Interface())
	}

	return row.Scan(addresses...)
}

// handleCustomStruct handles the correct pointer selection for custom structs types
//
// Parameters:
//   - reflect.Value:	the custom struct
//
// Returns:
//   - reflect.Value: address to be scanned
func handleCustomStructField(tf reflect.Value) (reflect.Value, error) {
	if utils.ValidateDefaultStruct(tf.Type()) {
		tf = tf.FieldByName("Value")
		return tf, nil
	}

	if utils.ValidateNullStruct(tf.Type()) {
		tf = tf.FieldByName("Value")
		return tf, nil
	}

	if utils.ValidateTimestampStruct(tf.Type()) {
		tf = tf.FieldByName("Datetime")
		return tf, nil
	}

	return reflect.Value{}, fmt.Errorf("no custom struct has been found: %v", tf)
}
