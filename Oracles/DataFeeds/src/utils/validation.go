package utils

import (
	"reflect"
	"regexp"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/types"
)

// Takes a value and returns true if it is a custom defined struct
// Custom defined structs are: Default and Null
//
// Parameters:
//   - data:	the value
//
// Returns:
//   - bool:	if the value is a struct or not
func ValidateCustomStruct(t reflect.Type) bool {
	return ValidateStruct(t) && (ValidateDefaultStruct(t) || ValidateNullStruct(t) || ValidateTimestampStruct(t))
}

// Takes a value and returns true if it is a struct
//
// Parameters:
//   - data:	the value
//
// Returns:
//   - bool:	if the value is a struct or not
func ValidateStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

// Checks whether a type is a types.Default regardless of the
// generic type of the struct
//
// Parameters:
//   - t:		the reflect type
//
// Returns:
//   - bool:	if the type is a generic types.Default or not
func ValidateDefaultStruct(t reflect.Type) bool {
	return types.DEFAULT.Kind() == t.Kind() && types.DEFAULT.PkgPath() == t.PkgPath() && BaseTypeName(types.DEFAULT) == BaseTypeName(t)
}

// Checks whether a type is a types.Null regardless of the
// generic type of the struct
//
// Parameters:
//   - t:		the reflect type
//
// Returns:
//   - bool:	if the type is a generic types.Null or not
func ValidateNullStruct(t reflect.Type) bool {
	return types.NULL.Kind() == t.Kind() && types.NULL.PkgPath() == t.PkgPath() && BaseTypeName(types.NULL) == BaseTypeName(t)
}

// Checks whether a type is a types.Timestamp regardless of the
// generic type of the struct
//
// Parameters:
//   - t:		the reflect type
//
// Returns:
//   - bool:	if the type is a generic types.Timestamp or not
func ValidateTimestampStruct(t reflect.Type) bool {
	return types.TIMESTAMP.Kind() == t.Kind() && types.TIMESTAMP.PkgPath() == t.PkgPath() && BaseTypeName(types.TIMESTAMP) == BaseTypeName(t)
}

// Takes a string and checks its validity as a query
//
// Parameters:
//   - query:	the query string
//
// Returns:
//   - bool:	if the query is valid or not
func ValidateQuery(query string) bool {
	sqlRegex := regexp.MustCompile(`(INSERT|UPDATE|DELETE|CREATE|DROP|ALTER|TRUNCATE|REPLACE|GRANT|REVOKE).*`)
	return sqlRegex.MatchString(query)
}

// Takes a string and checks its validity as a query with a result
//
// Parameters:
//   - query:	the query string
//
// Returns:
//   - bool:	if the query is valid or not
func ValidateQueryWithResult(query string) bool {
	sqlRegex := regexp.MustCompile(`SELECT.*`)
	return sqlRegex.MatchString(query)
}
