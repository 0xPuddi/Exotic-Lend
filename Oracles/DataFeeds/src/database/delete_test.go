package database

import (
	"reflect"
	"testing"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/types"
)

type DeleteRowsPrimaryKey struct {
	InitElem              []types.Table
	BeforeDeletionOrderId []uint
	DeleteElemId          uint8
	AfterDeletionOrderId  []uint
}

var DELETE_ROWS_PRIMARY_KEY_SELECTION = DeleteRowsPrimaryKey{
	InitElem: []types.Table{
		types.Asset{
			Id: types.Default[uint64]{
				Default: true,
				Value:   1,
			},
			Ticker:   "BTC",
			Source:   "Binance",
			Decimals: 18,
		},
		types.Price{
			Id: types.Default[int64]{
				Default: true,
				Value:   1,
			},
			Asset_id: 1,
			Price:    696969,
			Timestamp: types.Timestamp{
				Now:      true,
				Datetime: "",
				Unix:     0,
			},
		},
		types.Price{
			Id: types.Default[int64]{
				Default: true,
				Value:   2,
			},
			Asset_id: 1,
			Price:    420,
			Timestamp: types.Timestamp{
				Now:      true,
				Datetime: "",
				Unix:     0,
			},
		},
		types.Price{
			Id: types.Default[int64]{
				Default: true,
				Value:   3,
			},
			Asset_id: 1,
			Price:    1,
			Timestamp: types.Timestamp{
				Now:      true,
				Datetime: "",
				Unix:     0,
			},
		},
	},
	BeforeDeletionOrderId: []uint{
		3, 2, 1,
	},
	DeleteElemId: 3,
	AfterDeletionOrderId: []uint{
		2, 1,
	},
}

func TestDeleteRowsByPrimaryKeyWithSelectionQueryFunc(t *testing.T) {
	_, db, cleanup, err := InitMockSqlDB()
	if err != nil {
		t.Fatalf("mock db didnt initialized properly: %v", err)
	}
	defer cleanup()

	for _, el := range DELETE_ROWS_PRIMARY_KEY_SELECTION.InitElem {
		_, err := InsertEntry(db, el)
		if err != nil {
			t.Errorf("Error when inserting element: %v\n", err)
		}
	}

	rows, err := SelectAllWhereAssetIdOrderedRow(db, types.Price{}, 1, 2, 3, false)
	if err != nil {
		t.Errorf("error when selecting rows, %v", err)
	}

	var resultPrices []types.Price
	for rows.Next() {
		price := types.Price{}

		err := ScanRowToStruct(rows, reflect.ValueOf(&price).Elem())
		if err != nil {
			t.Errorf("error scanning row: %v", err)
			continue
		}
		resultPrices = append(resultPrices, price)
	}

	for i, id := range DELETE_ROWS_PRIMARY_KEY_SELECTION.BeforeDeletionOrderId {
		if id != uint(resultPrices[i].Id.Value) {
			t.Errorf("error retriving correct rows, wanted: %d, given: %d\n", id, resultPrices[i].Id.Value)
		}
	}

	strWAI, err := buildSelectWhereAssetIdOrderedRowConditions(types.PRICE, 1, 2, 1, false)
	if err != nil {
		t.Errorf("error when building select where asset id ordered row conditions, %v", err)
	}
	str, err := buildSelectConditionsQuery(types.PRICE, []int{0}, strWAI...)
	if err != nil {
		t.Errorf("error when building select conditions query, %v", err)
	}
	rowsAffected, err := DeleteRowsByPrimaryKeyWithSelectionQuery(db, types.Price{}, str)
	if err != nil {
		t.Errorf("error deleting rows: %v", err)
		return
	}
	ra, err := rowsAffected.RowsAffected()
	if err != nil {
		t.Errorf("erro fetching rows affected: %d", err)
	}
	t.Logf("rows affected: %d", ra)

	rows, err = SelectAllWhereAssetIdOrderedRow(db, types.Price{}, 1, 2, 3, false)
	if err != nil {
		t.Errorf("error when selecting rows, %v", err)
	}

	var resultPricesTwo []types.Price
	for rows.Next() {
		price := types.Price{}

		err := ScanRowToStruct(rows, reflect.ValueOf(&price).Elem())
		if err != nil {
			t.Errorf("error scanning row: %v", err)
			continue
		}
		resultPricesTwo = append(resultPricesTwo, price)
	}

	for i, id := range DELETE_ROWS_PRIMARY_KEY_SELECTION.AfterDeletionOrderId {
		t.Logf("Row price id: %d", resultPricesTwo[i].Id.Value)
		if id != uint(resultPricesTwo[i].Id.Value) {
			t.Errorf("error retriving correct rows, wanted: %d, given: %d\n", id, resultPricesTwo[i].Id.Value)
		}
	}
}
