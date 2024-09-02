package database

import (
	"reflect"
	"testing"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/types"
	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/utils"
)

// Types
type TestParseStruct struct {
	Order_detail_id types.Default[int64] `json:"order_detail_id" db:"order_detail_id SERIAL PRIMARY KEY"`
	Order_id        int64                `json:"order_id" db:"order_id INTEGER NOT NULL" ref:"FOREIGN KEY (order_id) REFERENCES Assets(order_id)"`
	Product_id      int64                `json:"product_id" db:"product_id INTEGER NOT NULL"`
	Quantity        int64                `json:"quantity" db:"quantity INTEGER NOT NULL CHECK (quantity > 0)" ref:"FOREIGN KEY (quantity) REFERENCES Assets(asset_id)" idx:"CREATE INDEX idx_test_quantity ON TestParseStruct(quantity)"`
	Unit_price      float32              `json:"unit_price" db:"unit_price DECIMAL(10, 2) NOT NULL CHECK (unit_price > 0)" idx:"CREATE INDEX idx_test_unit_price ON TestParseStruct(unit_price)"`
}

type TestCheckIfTableExistsInput struct {
	Input   any
	Create  bool
	Correct bool
}

type TestAddEntryInput struct {
	Input any
}

type TestInput struct {
	Input   any
	Correct any
}

// Tests

// Parse Struct
var PARSING_TO_TABLE_SAMPLES = []TestInput{
	{
		Input: types.Asset{
			Id: types.Default[uint64]{
				Default: true,
				Value:   0,
			},
			Ticker:   "BTC",
			Source:   "Binance",
			Decimals: 18,
		},
		Correct: `CREATE TABLE IF NOT EXISTS Asset (
	id SERIAL PRIMARY KEY,
	ticker VARCHAR(16) NOT NULL,
	source VARCHAR(16) NOT NULL,
	decimals SMALLINT NOT NULL CHECK (decimals >= 0)
);`,
	},
	{
		Input: types.Price{
			Asset_id: 1234,
			Price:    12400,
			Timestamp: types.Timestamp{
				Now:  true,
				Unix: 12,
			},
		},
		Correct: `CREATE TABLE IF NOT EXISTS Price (
	id SERIAL PRIMARY KEY,
	asset_id INTEGER NOT NULL,
	price BIGINT NOT NULL,
	timestamp TIMESTAMP DEFAULT NOW() NOT NULL,
	FOREIGN KEY (asset_id) REFERENCES asset(id)
);
CREATE INDEX idx_price_asset_id ON Price(asset_id);`,
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
		Correct: `CREATE TABLE IF NOT EXISTS TestParseStruct (
	order_detail_id SERIAL PRIMARY KEY,
	order_id INTEGER NOT NULL,
	product_id INTEGER NOT NULL,
	quantity INTEGER NOT NULL CHECK (quantity > 0),
	unit_price DECIMAL(10, 2) NOT NULL CHECK (unit_price > 0),
	FOREIGN KEY (order_id) REFERENCES Assets(order_id),
	FOREIGN KEY (quantity) REFERENCES Assets(asset_id)
);
CREATE INDEX idx_test_quantity ON TestParseStruct(quantity);
CREATE INDEX idx_test_unit_price ON TestParseStruct(unit_price);`,
	},
}

func TestParseStructToTableFunc(t *testing.T) {
	for _, ps := range PARSING_TO_TABLE_SAMPLES {
		str, err := ParseStructToTable(reflect.TypeOf(ps.Input))

		if err != nil {
			t.Errorf("error parsing the struct:\n%v", err)
			continue
		}

		if str != ps.Correct {
			t.Errorf("incorrect parsing:\n%v\n%v", str, ps.Correct)
		}
	}
}

var CHECK_IF_TABLE_EXISTS = []TestCheckIfTableExistsInput{
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
		Create:  false,
		Correct: false,
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
		Create:  true,
		Correct: true,
	},
	{
		Input: types.Price{
			Asset_id: 12,
			Price:    444,
			Timestamp: types.Timestamp{
				Now:  true,
				Unix: 10,
			},
		},
		Create:  true,
		Correct: true,
	},
}

func TestCheckIfTableExistsFunc(t *testing.T) {
	// Start Mock db
	_, db, cleanup, err := InitMockSqlDB()
	if err != nil {
		utils.HandleFatalError(err)
	}
	defer cleanup()

	err = db.Ping()
	if err != nil {
		t.Errorf("Failed to connect to the database: %v", err)
		return
	} else {
		t.Log("\n\nSuccessfully connected to the database!\n\n")
	}

	for _, input := range CHECK_IF_TABLE_EXISTS {
		if input.Create {
			// Create Table
			query, err := ParseStructToTable(reflect.TypeOf(input.Input))
			if err != nil {
				t.Errorf("error on struct parsing: %v", err)
			}

			_, err = MakeQuery(db, query)
			if err != nil {
				t.Errorf("error when making the query: %v", err)
			}
		}

		exists, err := CheckIfTableExists(db, input.Input)

		if err != nil {
			t.Errorf("error on table check: %v", err)
			continue
		}

		if exists != input.Correct {
			t.Errorf("wrong table result: %v %v", exists, input.Correct)
		}
	}
}

func TestCreateTableFunc(t *testing.T) {
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

	for _, input := range CHECK_IF_TABLE_EXISTS {
		if input.Create {
			// Create Table
			_, err := CreateTable(db, input.Input)

			if err != nil {
				t.Errorf("failed to create the table: %v", err)
			}
		}

		exists, err := CheckIfTableExists(db, input.Input)

		if err != nil {
			t.Errorf("error on table check: %v", err)
			continue
		}

		if exists != input.Correct {
			t.Errorf("wrong table result: %v %v", exists, input.Correct)
		}
	}

}
