package utils

import (
	"reflect"
	"testing"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/types"
)

type TestInput[I, T any] struct {
	Input   I
	Correct T
}

// Check if Custom
var CHECK_CUSTOM_SAMPLES = []TestInput[any, bool]{
	{
		Input: types.Default[string]{
			Default: false,
			Value:   "test",
		},
		Correct: true,
	},
	{
		Input: types.Default[any]{
			Default: false,
			Value:   123,
		},
		Correct: true,
	},
	{
		Input: types.Null[any]{
			Null:  false,
			Value: 123,
		},
		Correct: true,
	},
	{
		Input: types.Asset{
			Id: types.Default[uint64]{
				Default: true,
				Value:   0,
			},
			Ticker:   "USD",
			Source:   "BIN",
			Decimals: 123,
		},
		Correct: false,
	},
	{
		Input:   17,
		Correct: false,
	},
	{
		Input:   "string",
		Correct: false,
	},
}

func TestValidateStructIsCustomFunc(t *testing.T) {
	for _, ps := range CHECK_CUSTOM_SAMPLES {
		tps := reflect.TypeOf(ps.Input)
		isDefault := ValidateCustomStruct(tps)

		if isDefault != ps.Correct {
			t.Errorf("error struct is not default: %v", tps)
		}
	}
}

// Check if default
var CHECK_DEFAULT_SAMPLES = []TestInput[any, bool]{
	{
		Input: types.Default[string]{
			Default: false,
			Value:   "test",
		},
		Correct: true,
	},
	{
		Input: types.Default[any]{
			Default: false,
			Value:   123,
		},
		Correct: true,
	},
	{
		Input: types.Asset{
			Id: types.Default[uint64]{
				Default: true,
				Value:   1,
			},
			Ticker:   "USD",
			Source:   "BIN",
			Decimals: 123,
		},
		Correct: false,
	},
	{
		Input:   17,
		Correct: false,
	},
	{
		Input:   "string",
		Correct: false,
	},
}

func TestValidateDefaultFunc(t *testing.T) {
	for _, ps := range CHECK_DEFAULT_SAMPLES {
		tps := reflect.TypeOf(ps.Input)
		isDefault := ValidateDefaultStruct(tps)

		if isDefault != ps.Correct {
			t.Errorf("error struct is not default: %v", tps)
		}
	}
}

func TestValidateStruct(t *testing.T) {
	for _, ps := range CHECK_DEFAULT_SAMPLES {
		tps := reflect.TypeOf(ps.Input)
		isDefault := ValidateDefaultStruct(tps)

		if isDefault != ps.Correct {
			t.Errorf("error struct is not default: %v", tps)
		}
	}
}

var VALIDATE_QUERY_SAMPLES = []TestInput[string, bool]{
	{Input: "", Correct: false},
	// Checks for whitespace
	{Input: "INSERT", Correct: true},
	{Input: "CREATE IF NOT EXISTS", Correct: true},
	{Input: "UPDATE", Correct: true},
	{Input: "test", Correct: false},
	{Input: "NOTVALID", Correct: false},
	{Input: "SELECT NOT NULL", Correct: false},
}

func TestValidateQueryFunc(t *testing.T) {
	for _, ti := range VALIDATE_QUERY_SAMPLES {
		result := ValidateQuery(ti.Input)

		if result != ti.Correct {
			t.Errorf("Validate Query error, should have been %t instead is %t, input: %s", ti.Correct, result, ti.Input)
		}

	}
}

var VALIDATE_QUERY_WITH_RESULT_SAMPLES = []TestInput[string, bool]{
	{Input: "", Correct: false},
	// Checks for whitespace
	{Input: "INSERT", Correct: false},
	{Input: "CREATE", Correct: false},
	{Input: "UPDATE", Correct: false},
	{Input: "test", Correct: false},
	{Input: "NOTVALID", Correct: false},
	{Input: "SELECT * FROM Asset", Correct: true},
}

func TestValidateQueryWithResultFunc(t *testing.T) {
	for _, ti := range VALIDATE_QUERY_WITH_RESULT_SAMPLES {
		result := ValidateQueryWithResult(ti.Input)

		if result != ti.Correct {
			t.Errorf("Validate Query error, should have been %t instead is %t, input: %s", ti.Correct, result, ti.Input)
		}
	}
}
