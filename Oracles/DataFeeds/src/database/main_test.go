package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/config"
	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/utils"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/joho/godotenv"
)

// Utils
// InitMockSqlDB inits a mock database for testing, it provides it
// with a sql driver and a cleanup function to defer
func InitMockSqlDB() (*embeddedpostgres.EmbeddedPostgres, *sql.DB, func(), error) {
	port, err := strconv.ParseUint(os.Getenv(config.DB_EMBED_PORT), 10, 32)
	if err != nil {
		utils.HandleFatalError(err)
	}

	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Username(os.Getenv(config.DB_EMBED_USERNAME)).
		Password(os.Getenv(config.DB_EMBED_PASSWORD)).
		Database(os.Getenv(config.DB_EMBED_DATABASE)).
		Port(uint32(port)).
		Version(embeddedpostgres.V16))

	// RuntimePath(RUNTIME_PATH).

	err = postgres.Start()
	if err != nil {
		return nil, nil, nil, err
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s port=%d",
		os.Getenv(config.DB_EMBED_USERNAME),
		os.Getenv(config.DB_EMBED_PASSWORD),
		os.Getenv(config.DB_EMBED_DATABASE),
		os.Getenv(config.DB_EMBED_SSL_MODE), port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, nil, err
	}

	cleanup := func() {
		fmt.Println("\n\nclenaing up database connection")
		// Close db connection
		if err := db.Close(); err != nil {
			utils.HandleGracefulError("error when closing the database connection", err)
		}

		fmt.Println("\n\nclenaing up mock database instance")
		// Stop postgres
		if err := postgres.Stop(); err != nil {
			utils.HandleGracefulError("error when shutting down mock embedded postgres database", err)
		}
	}

	return postgres, db, cleanup, nil
}

// Utils
// PrintRowsValues prints all row values
func PrintRowsValues(rows *sql.Rows) error {
	fmt.Println("**********************************")
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("unable to get columns: %v", err)
	}
	fmt.Println("Columns:", columns)

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for i := range values {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		fmt.Printf("\n")
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return fmt.Errorf("unable to scan columns: %v", err)
		}
		for i, col := range columns {
			fmt.Printf("%s: %v\n", col, values[i])
		}
	}
	fmt.Println("**********************************")

	return nil
}

// TestMain runs before any test in this package
func TestMain(m *testing.M) {
	// Setup code before tests
	fmt.Println("Running tests")
	// Path to src
	utils.HandleFatalError(godotenv.Load("../../.env"))

	// Run the tests
	exitCode := m.Run()

	// Teardown code after tests
	fmt.Println("Cleaning up after tests")

	// Exit with the appropriate code
	os.Exit(exitCode)
}
