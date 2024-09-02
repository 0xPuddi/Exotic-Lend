package database

import (
	"database/sql"
	"fmt"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/types"
	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/utils"
)

// Takes multiple queries, it checks and performs them
//
// Parameters:
//   - db:		the database struct
//   - queries:	the queries strings
//
// Returns:
//   - QueriesResult:	an array of results of the same length as entries
func MakeQueries(db *sql.DB, queries []string) []types.QueryResult {
	results := make([]types.QueryResult, len(queries))

	for i, query := range queries {
		result, err := MakeQuery(db, query)
		results[i] = types.QueryResult{
			Result: result,
			Error:  err,
		}
	}

	return results
}

// Checks the validity of the query and, if it pases, it makes it
// Valid queries are: INSERT, UPDATE, DELETE
//
// Parameters:
//   - db:		the database struct
//   - query:	the query string
//
// Returns:
//   - Result:	the query result
//   - error:	an error if occured
func MakeQuery(db *sql.DB, query string) (sql.Result, error) {
	if !utils.ValidateQuery(query) {
		return nil, fmt.Errorf("query is not valid")
	}

	return db.Exec(query)
}

// Takes multiple queries, it checks and performs them
//
// Parameters:
//   - db:		the database struct
//   - queries:	the queries strings
//
// Returns:
//   - QueriesResult:	an array of results of the same length as entries
func MakeQueriesWithResult(db *sql.DB, queries []string) []types.QueryRows {
	results := make([]types.QueryRows, len(queries))

	for i, query := range queries {
		result, err := MakeQueryWithResult(db, query)
		results[i] = types.QueryRows{
			Result: result,
			Error:  err,
		}
	}

	return results
}

// Checks the validity of the query with a result and, if it pases,
// it makes it. Valid queries are: SELECT
// Notice query must be complete
//
// Parameters:
//   - db:		the database struct
//   - query:	the query string
//
// Returns:
//   - Result:	the query result
//   - error:	an error if occured
func MakeQueryWithResult(db *sql.DB, query string) (*sql.Rows, error) {
	if !utils.ValidateQueryWithResult(query) {
		return nil, fmt.Errorf("query is not valid")
	}

	return db.Query(query)
}
