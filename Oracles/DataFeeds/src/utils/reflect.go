package utils

import (
	"fmt"
	"reflect"
	"strings"
)

// Returns the type name, regardless of its generic types
//
// Parameters:
//   - t:	the reflect type
//
// Returns:
//   - string:	the string name
func BaseTypeName(t reflect.Type) string {
	name := t.String()

	last_idx := strings.Index(name, "[")
	first_idx := strings.LastIndex(name, ".")

	if last_idx != -1 && first_idx != -1 {
		return name[first_idx+1 : last_idx]
	} else if first_idx != -1 {
		return name[first_idx+1:]
	} else if last_idx != -1 {
		return name[:last_idx]
	}

	return name
}

// Returns the db tag name
//
// Parameters:
//   - t:		the reflect type
//
// Returns:
//   - string:	the string tag name
func GetFieldNameDB(sf reflect.StructField) (string, error) {
	str, ok := sf.Tag.Lookup("db")
	if !ok {
		return "", fmt.Errorf("struct field doesn't have a db tag")
	}
	return strings.Split(str, " ")[0], nil
}
