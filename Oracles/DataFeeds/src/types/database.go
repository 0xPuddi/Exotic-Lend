package types

import (
	"database/sql"
)

// Query result wrapper
type QueryResult struct {
	Result sql.Result
	Error  error
}

// Query result wrapper
type QueryRows struct {
	Result *sql.Rows
	Error  error
}
