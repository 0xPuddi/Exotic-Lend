package feeds

import (
	"fmt"
	"os"
	"testing"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/utils"
	"github.com/joho/godotenv"
)

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
