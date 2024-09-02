package types

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	DEFAULT   = reflect.TypeOf(Default[any]{})
	NULL      = reflect.TypeOf(Null[any]{})
	TIMESTAMP = reflect.TypeOf(Timestamp{})
	ASSET     = reflect.TypeOf(Asset{})
	PRICE     = reflect.TypeOf(Price{})
)

// Default type to add for each column that can be added as default
//
// Note that if `Default` is false, `Value` will be used instead of `DEFAULT`
type Default[T any] struct {
	Default bool
	Value   T
}

// Null type to add for each column that can be added as null:
//
// Note that if `Null` is false, `Value` will be used instead of `NULL`
type Null[T any] struct {
	Null  bool
	Value T
}

// Timestamp type to add for each column that can be added as null:
//
// Note that if `Null` is false, `Value` will be used instead of `NULL`
type Timestamp struct {
	Now      bool
	Datetime string
	Unix     int
}

// Where condition type
type WhereCondition struct {
	Column    int
	Condition string
}

// Table interface
type Table interface {
	GetPrimaryKeyNameDB() (string, error)
}

// Price struct
//
// Many to One relation with Assset
type Price struct {
	Id        Default[int64] `json:"id"  db:"id SERIAL PRIMARY KEY"`
	Asset_id  int            `json:"asset_id"  db:"asset_id INTEGER NOT NULL" ref:"FOREIGN KEY (asset_id) REFERENCES asset(id)" idx:"CREATE INDEX idx_price_asset_id ON Price(asset_id)"`
	Price     int            `json:"price"     db:"price BIGINT NOT NULL"`
	Timestamp Timestamp      `json:"timestamp" db:"timestamp TIMESTAMP DEFAULT NOW() NOT NULL"`
}

func (p Price) GetPrimaryKeyNameDB() (string, error) {
	str, err := getFieldNameDB(reflect.TypeOf(p).Field(0))
	if err != nil {
		return "", err
	}

	return str, nil
}

// Asset struct
type Asset struct {
	Id       Default[uint64] `josn:"id"       db:"id SERIAL PRIMARY KEY"`
	Ticker   string          `json:"ticker"   db:"ticker VARCHAR(16) NOT NULL"`
	Source   string          `json:"source"   db:"source VARCHAR(16) NOT NULL"`
	Decimals int8            `json:"decimals" db:"decimals SMALLINT NOT NULL CHECK (decimals >= 0)"`
}

func (a Asset) GetPrimaryKeyNameDB() (string, error) {
	str, err := getFieldNameDB(reflect.TypeOf(a).Field(0))
	if err != nil {
		return "", err
	}

	return str, nil
}

// Returns the db tag name
//
// Parameters:
//   - t:		the reflect type
//
// Returns:
//   - string:	the string tag name
func getFieldNameDB(sf reflect.StructField) (string, error) {
	str, ok := sf.Tag.Lookup("db")
	if !ok {
		return "", fmt.Errorf("struct field doesn't have a db tag")
	}
	return strings.Split(str, " ")[0], nil
}
