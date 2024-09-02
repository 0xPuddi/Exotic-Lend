package utils

import (
	"reflect"
	"testing"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/types"
)

type BaseNameStruct[T, V any] struct {
	amount      int64
	description string
	custom_one  T
	custom_two  V
}

type TestParseStruct struct {
	Order_detail_id types.Default[int64] `json:"order_detail_id" db:"order_detail_id SERIAL PRIMARY KEY"`
	Order_id        int64                `json:"order_id" db:"order_id INTEGER NOT NULL" ref:"FOREIGN KEY (order_id) REFERENCES Assets(order_id)"`
	Product_id      int64                `json:"product_id" db:"product_id INTEGER NOT NULL"`
	Quantity        int64                `json:"quantity" db:"quantity INTEGER NOT NULL CHECK (quantity > 0)" ref:"FOREIGN KEY (quantity) REFERENCES Assets(asset_id)" idx:"CREATE INDEX idx_test_quantity ON TestParseStruct(quantity)"`
	Unit_price      float32              `json:"unit_price" db:"unit_price DECIMAL(10, 2) NOT NULL CHECK (unit_price > 0)" idx:"CREATE INDEX idx_test_unit_price ON TestParseStruct(unit_price)"`
}

// Base Type Name
var BASE_NAME_SAMPLES = []TestInput[any, string]{
	{
		Input: BaseNameStruct[string, int64]{
			amount:      123,
			description: "lul",
			custom_one:  "lil",
			custom_two:  120,
		},
		Correct: "BaseNameStruct",
	},
	{
		Input: types.Price{
			Asset_id: 1234,
			Price:    12400,
			Timestamp: types.Timestamp{
				Now:  true,
				Unix: 123456789,
			},
		},
		Correct: "Price",
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
		Correct: "TestParseStruct",
	},
}

func TestBaseTypeNameFunc(t *testing.T) {
	for _, ps := range BASE_NAME_SAMPLES {
		tps := reflect.TypeOf(ps.Input)
		str := BaseTypeName(tps)

		if str != ps.Correct {
			t.Errorf("error when getting base type name,\nname: %v\ncorrect: %v", str, ps.Correct)
		}
	}
}

// Get db samples
var PARSING_GET_NAME_DB_SAMPLES = []TestInput[any, []string]{
	{
		Input: types.Asset{
			Id: types.Default[uint64]{
				Default: true,
				Value:   11,
			},
			Ticker:   "BTC",
			Source:   "Binance",
			Decimals: 18,
		},
		Correct: []string{"id", "ticker", "source", "decimals"},
	},
	{
		Input: types.Price{
			Asset_id: 1234,
			Price:    12400,
			Timestamp: types.Timestamp{
				Now:  false,
				Unix: 123456789,
			},
		},
		Correct: []string{"id", "asset_id", "price", "timestamp"},
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
		Correct: []string{"order_detail_id", "order_id", "product_id", "quantity", "unit_price"},
	},
}

func TestGetTableNameDBFunc(t *testing.T) {
	for _, ps := range PARSING_GET_NAME_DB_SAMPLES {
		tps := reflect.TypeOf(ps.Input)

		for i := 0; i < tps.NumField(); i++ {
			f := tps.Field(i)

			str, err := GetFieldNameDB(f)

			if err != nil {
				t.Errorf("error when getting")
			}

			if str != ps.Correct[i] {
				t.Errorf("error getting the db name")
			}
		}
	}
}
