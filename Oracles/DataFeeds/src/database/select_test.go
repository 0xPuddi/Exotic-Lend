package database

import (
	"reflect"
	"testing"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/types"
)

// Types
type SelectColumnsInput struct {
	Struct  any
	Columns []int
}
type SelectColumnsCorrect struct {
	Correct bool
	Query   string
}
type TestSelectColumnsInput struct {
	Input   SelectColumnsInput
	Correct SelectColumnsCorrect
}

// Tests

// Test Parse Struct Select Columns
var STRUCTS_TO_SELECT_COLUMNS = []TestSelectColumnsInput{
	{
		Input: SelectColumnsInput{
			Struct: types.Asset{
				Id: types.Default[uint64]{
					Default: true,
					Value:   1,
				},
				Ticker:   "BTC",
				Source:   "Binance",
				Decimals: 0,
			},
			Columns: []int{0, 2, 3},
		},
		Correct: SelectColumnsCorrect{
			Correct: true,
			Query: `SELECT id, source, decimals
FROM Asset`,
		},
	},
	{
		Input: SelectColumnsInput{
			Struct: types.Asset{
				Id: types.Default[uint64]{
					Default: true,
					Value:   1,
				},
				Ticker:   "BTC",
				Source:   "Binance",
				Decimals: 0,
			},
			Columns: []int{0, 2, 3, 1, 3, 2},
		},
		Correct: SelectColumnsCorrect{
			Correct: false,
			Query:   ``,
		},
	},
	{
		Input: SelectColumnsInput{
			Struct:  123,
			Columns: []int{},
		},
		Correct: SelectColumnsCorrect{
			Correct: false,
			Query:   ``,
		},
	},
	{
		Input: SelectColumnsInput{
			Struct: types.Asset{
				Id: types.Default[uint64]{
					Default: true,
					Value:   1,
				},
				Ticker:   "BTC",
				Source:   "Binance",
				Decimals: 0,
			},
			Columns: []int{},
		},
		Correct: SelectColumnsCorrect{
			Correct: false,
			Query:   ``,
		},
	},
	{
		Input: SelectColumnsInput{
			Struct: types.Asset{
				Id: types.Default[uint64]{
					Default: true,
					Value:   1,
				},
				Ticker:   "BTC",
				Source:   "Binance",
				Decimals: 0,
			},
			Columns: []int{6},
		},
		Correct: SelectColumnsCorrect{
			Correct: false,
			Query:   ``,
		},
	},
}

func TestParseStructToSelectColumns(t *testing.T) {
	for _, i := range STRUCTS_TO_SELECT_COLUMNS {
		query, err := parseStructToSelectColumns(reflect.TypeOf(i.Input.Struct), i.Input.Columns...)

		if i.Correct.Correct == true {
			if err != nil || query != i.Correct.Query {
				t.Errorf("error or incorrect query output: \n%s\n%v", query, err)
			}
		} else {
			if err == nil || query != "" {
				t.Errorf("incorrect input gives a correct output: \n%s\n%v", query, err)
			}
		}
	}
}

// Test Build Selection Query
type BuildSelectionQueriesInput struct {
	Conditions []string
	Columns    []int
	Table      any
}
type BuildSelectionQueries struct {
	Input   BuildSelectionQueriesInput
	Correct string
}

var BUILD_SELECTION_QUERIES = []BuildSelectionQueries{
	{
		Input: BuildSelectionQueriesInput{
			Conditions: []string{
				"ORDER BY timestamp DESC",
				"LIMIT 10",
			},
			Columns: []int{},
			Table: types.Asset{
				Id: types.Default[uint64]{
					Default: true,
					Value:   1,
				},
			},
		},
		Correct: `SELECT * FROM Asset
ORDER BY timestamp DESC
LIMIT 10`,
	},
	{
		Input: BuildSelectionQueriesInput{
			Conditions: []string{
				"ORDER BY timestamp DESC",
				"LIMIT 10",
			},
			Columns: []int{2, 3},
			Table: types.Asset{
				Id: types.Default[uint64]{
					Default: true,
					Value:   1,
				},
				Ticker:   "lul",
				Source:   "SOS",
				Decimals: 11,
			},
		},
		Correct: `SELECT source, decimals FROM Asset
ORDER BY timestamp DESC
LIMIT 10`,
	},
}

func TestBuildSelectionQueriesFunc(t *testing.T) {
	for _, i := range BUILD_SELECTION_QUERIES {
		str, err := buildSelectConditionsQuery(reflect.TypeOf(i.Input.Table), i.Input.Columns, i.Input.Conditions...)
		if err != nil {
			t.Errorf("error building conditions query: %v", err)
		}

		if str != i.Correct {
			t.Errorf("incorrect build result: \n%v\n%v", str, i.Correct)
		}
	}
}

// Most Recent Rows
type MostRecentRowsInput struct {
	EntryRows []any
	Column    int
	Limit     int
}
type MostRecentRows struct {
	Input   MostRecentRowsInput
	Correct []int
}

var MOST_RECENT_ROWS = MostRecentRows{
	Input: MostRecentRowsInput{
		EntryRows: []any{
			types.Asset{
				Id: types.Default[uint64]{
					Default: false,
					Value:   0,
				},
				Ticker:   "BTC",
				Source:   "Binance",
				Decimals: 10,
			},
			types.Price{
				Id: types.Default[int64]{
					Default: true,
					Value:   0,
				},
				Asset_id: 0,
				Price:    600000000000000,
				Timestamp: types.Timestamp{
					Now:  true,
					Unix: 0,
				},
			},
			types.Price{
				Id: types.Default[int64]{
					Default: true,
					Value:   0,
				},
				Asset_id: 0,
				Price:    610000000000000,
				Timestamp: types.Timestamp{
					Now:  false,
					Unix: 1724526459,
				},
			},
			types.Price{
				Id: types.Default[int64]{
					Default: true,
					Value:   0,
				},
				Asset_id: 0,
				Price:    590000000000000,
				Timestamp: types.Timestamp{
					Now:  false,
					Unix: 1624526459,
				},
			},
			types.Price{
				Id: types.Default[int64]{
					Default: true,
					Value:   0,
				},
				Asset_id: 0,
				Price:    605000000000000,
				Timestamp: types.Timestamp{
					Now:  false,
					Unix: 1724525459,
				},
			},
		},
		Column: 10,
		Limit:  2,
	},
	Correct: []int{
		600000000000000,
		610000000000000,
		605000000000000,
		590000000000000,
	},
}

func TestSelectMostRecentRowsFunc(t *testing.T) {
	_, db, cleanup, err := InitMockSqlDB()
	if err != nil {
		t.Fatalf("DB failed to start: %v", err)
	}
	defer cleanup()

	results, errors := InsertEntries(db, MOST_RECENT_ROWS.Input.EntryRows)
	for _, err := range errors {
		if err != nil {
			t.Errorf("error during insertion: %v", err)
			return
		}
	}
	for _, res := range results {
		row, err := res.RowsAffected()

		if err != nil {
			t.Errorf("error when catching response: %v", err)
			continue
		}

		t.Logf("Rows: %d", row)
	}

	rows, err := SelectAllWhereAssetIdOrderedRow(db, types.Price{}, 0, 3, 4, true)
	if err != nil {
		t.Errorf("error when selecting most recent rows: %v", err)
		return
	}
	defer rows.Close()

	var prices []types.Price
	for rows.Next() {
		price := types.Price{}

		err := ScanRowToStruct(rows, reflect.ValueOf(&price).Elem())
		if err != nil {
			t.Errorf("error when scanning row to struct: %v", err)
		}
		prices = append(prices, price)
	}

	for i, p := range prices {
		t.Logf("Scanned row: %d, %d, %s", p.Asset_id, p.Price, p.Timestamp.Datetime)

		if p.Price != MOST_RECENT_ROWS.Correct[i] {
			t.Errorf("price retrived incorrectly, wanted %d, retrived %d", MOST_RECENT_ROWS.Correct[i], p.Price)
		}
	}
}

type BuildSelectTableMatch struct {
	table    any
	matchCol []int
	matchVal []any
	Limit    int
}
type BuildSelectTableMatchInput struct {
	Input   BuildSelectTableMatch
	Correct []string
}

var BUILD_SELECT_TABLE_MATCH = []BuildSelectTableMatchInput{
	{
		Input: BuildSelectTableMatch{
			table: types.Price{
				Id: types.Default[int64]{
					Default: true,
					Value:   10,
				},
				Asset_id: 124,
				Timestamp: types.Timestamp{
					Now:  true,
					Unix: 100,
				},
				Price: 1000,
			},
			matchCol: []int{1, 3},
			matchVal: []any{
				124,
				"4444",
			},
			Limit: 10,
		},
		Correct: []string{
			`WHERE asset_id = 124`,
			`AND timestamp = '4444'`,
			`LIMIT 10`,
		},
	},
}

func TestBuildSelectTableByMatch(t *testing.T) {
	for _, bstm := range BUILD_SELECT_TABLE_MATCH {
		strings, err := buildSelectTableByMatch(reflect.TypeOf(bstm.Input.table), bstm.Input.matchCol, bstm.Input.matchVal, bstm.Input.Limit)

		if err != nil {
			t.Errorf("error when building select table by match query: %v", err)
		}

		for i := 0; i < len(strings); i++ {
			if strings[i] != bstm.Correct[i] {
				t.Errorf("error building the query: \ngiven %v\nwanted %v", strings[i], bstm.Correct[i])
			}
		}
	}
}

type SelectTableMatch struct {
	table    any
	matchCol []int
	matchVal []any
	Limit    int
}
type SelectTableMatchInput struct {
	Entries []any
	Input   SelectTableMatch
	Correct []any
}

var SELECT_TABLE_MATCH = SelectTableMatchInput{
	Entries: []any{
		types.Asset{
			Id: types.Default[uint64]{
				Default: true,
				Value:   1,
			},
			Ticker:   "VOLTA",
			Source:   "Binance",
			Decimals: 18,
		},
		types.Asset{
			Id: types.Default[uint64]{
				Default: true,
				Value:   2,
			},
			Ticker:   "VOLTA",
			Source:   "Binance",
			Decimals: 18,
		},
		types.Price{
			Id: types.Default[int64]{
				Default: true,
				Value:   10,
			},
			Asset_id: 1,
			Timestamp: types.Timestamp{
				Now:  true,
				Unix: 100,
			},
			Price: 99,
		},
		types.Price{
			Id: types.Default[int64]{
				Default: true,
				Value:   10,
			},
			Asset_id: 1,
			Timestamp: types.Timestamp{
				Now:  true,
				Unix: 100,
			},
			Price: 1000,
		},
		types.Price{
			Id: types.Default[int64]{
				Default: true,
				Value:   10,
			},
			Asset_id: 1,
			Timestamp: types.Timestamp{
				Now:  false,
				Unix: 10000,
			},
			Price: 500,
		},
		types.Price{
			Id: types.Default[int64]{
				Default: true,
				Value:   10,
			},
			Asset_id: 1,
			Timestamp: types.Timestamp{
				Now:  false,
				Unix: 1000,
			},
			Price: 1000,
		},
		types.Price{
			Id: types.Default[int64]{
				Default: true,
				Value:   10,
			},
			Asset_id: 1,
			Timestamp: types.Timestamp{
				Now:  true,
				Unix: 100,
			},
			Price: 1000,
		},
		types.Price{
			Id: types.Default[int64]{
				Default: true,
				Value:   10,
			},
			Asset_id: 2,
			Timestamp: types.Timestamp{
				Now:  true,
				Unix: 100,
			},
			Price: 99,
		},
		types.Price{
			Id: types.Default[int64]{
				Default: true,
				Value:   10,
			},
			Asset_id: 2,
			Timestamp: types.Timestamp{
				Now:  true,
				Unix: 100,
			},
			Price: 99,
		},
	},
	Input: SelectTableMatch{
		table: types.Price{},
		matchCol: []int{
			1, 2,
		},
		matchVal: []any{
			2, 99,
		},
		Limit: -1,
	},
	Correct: []any{
		99,
		99,
	},
}

func TestSelectTableByMatchColumnsFunc(t *testing.T) {
	_, db, cleanup, err := InitMockSqlDB()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	// Build Database
	results, errors := InsertEntries(db, SELECT_TABLE_MATCH.Entries)
	for _, err := range errors {
		if err != nil {
			t.Errorf("error during insertion: %v", err)
			return
		}
	}
	for _, res := range results {
		row, err := res.RowsAffected()

		if err != nil {
			t.Errorf("error when catching response: %v", err)
			continue
		}

		t.Logf("Rows: %d", row)
	}

	// Make Query
	rows, err := SelectTableByMatchColumns(db, SELECT_TABLE_MATCH.Input.table, []int{}, SELECT_TABLE_MATCH.Input.matchCol, SELECT_TABLE_MATCH.Input.matchVal, SELECT_TABLE_MATCH.Input.Limit)
	if err != nil {
		t.Errorf("error when selecting most recent rows: %v", err)
		return
	}
	defer rows.Close()

	// Scan rows
	var prices []types.Price
	for rows.Next() {
		price := types.Price{}

		err := ScanRowToStruct(rows, reflect.ValueOf(&price).Elem())
		if err != nil {
			t.Errorf("error when scanning row to struct: %v", err)
		}
		prices = append(prices, price)
	}

	// Check Validity
	for i, p := range prices {
		if p.Price != SELECT_TABLE_MATCH.Correct[i] {
			t.Errorf("error selecting rows: \ngiven %v\nwanted %v", p.Price, SELECT_TABLE_MATCH.Correct[i])
		}
	}
}
