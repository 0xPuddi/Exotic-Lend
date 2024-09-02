package database

import (
	"reflect"
	"testing"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/types"
	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/utils"
)

// Parse struct to entry
var PARSING_TO_ENTRY = []TestInput{
	{
		Input: types.Asset{
			Id: types.Default[uint64]{
				Default: true,
				Value:   1,
			},
			Ticker:   "BTC",
			Source:   "Binance",
			Decimals: 18,
		},
		Correct: `INSERT INTO Asset (id, ticker, source, decimals)
VALUES (DEFAULT, 'BTC', 'Binance', 18)`,
	},
	{
		Input: TestParseStruct{
			Order_detail_id: types.Default[int64]{
				Default: true,
				Value:   1234,
			},
			Order_id:   12400,
			Product_id: 123456789,
			Quantity:   10,
			Unit_price: 123,
		},
		Correct: `INSERT INTO TestParseStruct (order_detail_id, order_id, product_id, quantity, unit_price)
VALUES (DEFAULT, 12400, 123456789, 10, 123)`,
	},
	{
		Input: types.Asset{
			Id: types.Default[uint64]{
				Default: true,
				Value:   1,
			},
			Ticker:   "BTC",
			Source:   "Binance",
			Decimals: 18,
		},
		Correct: `INSERT INTO Asset (id, ticker, source, decimals)
VALUES (DEFAULT, 'BTC', 'Binance', 18)`,
	},
	{
		Input: types.Asset{
			Id: types.Default[uint64]{
				Default: true,
				Value:   1,
			},
			Ticker:   "BTC",
			Source:   "Bin",
			Decimals: 18,
		},
		Correct: `INSERT INTO Asset (id, ticker, source, decimals)
VALUES (DEFAULT, 'BTC', 'Bin', 18)`,
	},
	{
		Input: types.Price{
			Id: types.Default[int64]{
				Default: true,
				Value:   0,
			},
			Asset_id: 12,
			Price:    444,
			Timestamp: types.Timestamp{
				Now:  false,
				Unix: 1724440501,
			},
		},
		Correct: `INSERT INTO Price (id, asset_id, price, timestamp)
VALUES (DEFAULT, 12, 444, TO_TIMESTAMP(1724440501))`,
	},
	{
		Input: types.Price{
			Id: types.Default[int64]{
				Default: false,
				Value:   0,
			},
			Asset_id: 12,
			Price:    444,
			Timestamp: types.Timestamp{
				Now:  true,
				Unix: 1,
			},
		},
		Correct: `INSERT INTO Price (id, asset_id, price, timestamp)
VALUES (0, 12, 444, NOW())`,
	},
}

func TestParseStructToEntry(t *testing.T) {
	for _, ps := range PARSING_TO_ENTRY {
		tps := reflect.TypeOf(ps.Input)
		vps := reflect.ValueOf(ps.Input)
		str, err := ParseStructToEntry(tps, vps)

		if err != nil {
			t.Errorf("error during parsing: %v", err)
		}

		if str != ps.Correct {
			t.Errorf("error when getting base type name,\nname: %v\ncorrect: %v", str, ps.Correct)
		}

	}
}

// Add entry
var INSERT_ENTRY = []TestAddEntryInput{
	{
		Input: types.Asset{
			Id: types.Default[uint64]{
				Default: true,
				Value:   1,
			},
			Ticker:   "BTC",
			Source:   "Binance",
			Decimals: 18,
		},
	},
	{
		Input: types.Asset{
			Id: types.Default[uint64]{
				Default: true,
				Value:   1,
			},
			Ticker:   "BTC",
			Source:   "Bin",
			Decimals: 18,
		},
	},
	{
		Input: types.Price{
			Id: types.Default[int64]{
				Default: true,
				Value:   0,
			},
			Asset_id: 1,
			Price:    444,
			Timestamp: types.Timestamp{
				Now:  true,
				Unix: 1724440501,
			},
		},
	},
	{
		Input: types.Price{
			Id: types.Default[int64]{
				Default: true,
				Value:   0,
			},
			Asset_id: 2,
			Price:    444,
			Timestamp: types.Timestamp{
				Now:  false,
				Unix: 1724440501,
			},
		},
	},
}

func TestInsertEntryFunc(t *testing.T) {
	// Start Mock db
	_, db, cleanup, err := InitMockSqlDB()
	if err != nil {
		utils.HandleFatalError(err)
	}
	defer cleanup()

	// Check connection
	err = db.Ping()
	if err != nil {
		t.Errorf("Failed to connect to the database: %v", err)
		return
	}
	t.Log("\n\nSuccessfully connected to the database!\n\n")

	for _, i := range INSERT_ENTRY {
		_, err := CreateTable(db, i.Input)
		if err != nil && err != ErrTableExists {
			t.Errorf("error when creating a table: %v", err)
		}

		iv := reflect.ValueOf(i.Input)
		it := reflect.TypeOf(i.Input)
		t.Log(ParseStructToEntry(it, iv))
		resEntry, err := InsertEntry(db, i.Input)
		if err != nil {
			t.Errorf("error when adding entry: %v", err)
			continue
		}

		affectedRows, err := resEntry.RowsAffected()
		if err != nil {
			t.Errorf("error when retriving affected rows: %v", err)
		}
		t.Logf("Rows affected: %d", affectedRows)

		resRows, err := SelectAll(db, i.Input)
		if err != nil {
			t.Errorf("unable to get select query response: %v", err)
			continue
		}
		defer resRows.Close()

		PrintRowsValues(resRows)
	}

}

// Parse Value
var PARSING_VALUE_TO_ENTRY = []TestInput{
	{
		Input:   "String",
		Correct: "'String'",
	},
	{
		Input:   true,
		Correct: "TRUE",
	},
	{
		Input:   -123,
		Correct: `-123`,
	},
}

func TestParseValueToEntryFunc(t *testing.T) {
	for _, v := range PARSING_VALUE_TO_ENTRY {
		vi := reflect.ValueOf(v.Input)
		result := ParseValueToEntry(vi)

		if result != v.Correct {
			t.Errorf("error parsing value result %v wanted %v", result, v.Correct)
		}
	}
}

// Parse Struct
var PARSING_CUSTOM_STRUCT = []TestInput{
	{
		Input: types.Default[any]{
			Default: true,
			Value:   123,
		},
		Correct: `DEFAULT`,
	},
	{
		Input: types.Default[any]{
			Default: false,
			Value:   123,
		},
		Correct: `123`,
	},
	{
		Input: types.Null[any]{
			Null:  true,
			Value: 123,
		},
		Correct: `NULL`,
	},
	{
		Input: types.Null[any]{
			Null:  false,
			Value: 123,
		},
		Correct: `123`,
	},
	{
		Input: TestParseStruct{
			Order_detail_id: types.Default[int64]{
				Default: true,
				Value:   1234,
			},
			Order_id:   12400,
			Product_id: 123456789,
			Quantity:   10,
			Unit_price: 123,
		},
		Correct: "",
	},
}

func TestParseCustomStructFunc(t *testing.T) {
	for _, ps := range PARSING_CUSTOM_STRUCT {
		str, _ := ParseCustomStruct(reflect.TypeOf(ps.Input), reflect.ValueOf(ps.Input))

		if str != ps.Correct {
			t.Errorf("incorrect parsing:\n%v\n%v", str, ps.Correct)
		}
	}
}
