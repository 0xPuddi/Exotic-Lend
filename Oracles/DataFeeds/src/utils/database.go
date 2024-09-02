package utils

import (
	"time"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/types"
)

// DatetimeToUnix translates a postgres timestamp to a unix timestamp
//
// Parameters:
//   - datetime:	the postgress timestamp
//
// Returns:
//   - int64:	the unix timestamp
//   - error: 	if any error occured
func DatetimeToUnix(datetime string) (int, error) {
	// Define the layout for the datetime format
	layout := "2006-01-02 15:04:05"

	// Parse the datetime string into a time.Time object
	t, err := time.Parse(layout, datetime)
	if err != nil {
		return 0, err
	}

	// Convert the time.Time object to Unix timestamp
	return int(t.Unix()), nil
}

// UpdateUnixTimestampFromDatetime takes the current Datetime
// timestamp of types.Price and updates its Unix timestamp
//
// Parameters:
//   - p:	the types.Price reference
//
// Returns:
//   - error: 	if any error occured
func UpdateUnixTimestampFromDatetime(p *types.Price) error {
	unix, err := DatetimeToUnix((*p).Timestamp.Datetime)
	if err != nil {
		return err
	}

	(*p).Timestamp.Unix = unix
	return nil
}
