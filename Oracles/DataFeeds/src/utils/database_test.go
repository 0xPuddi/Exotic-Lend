package utils

import (
	"testing"

	"github.com/0xPuddi/Exotic-Lend/Oracles/DataFeeds/types"
)

type DateToUnix struct {
	date string
	unix int
}

var DATE_TIMES_TO_UNIX = []DateToUnix{
	{
		date: "2024-08-25 12:00:00",
		unix: 1724587200,
	},
	{
		date: "2023-12-31 23:59:59",
		unix: 1704067199,
	},
	{
		date: "2022-01-01 00:00:00",
		unix: 1640995200,
	},
	{
		date: "2000-01-01 00:00:00",
		unix: 946684800,
	},
	{
		date: "2010-07-15 15:30:00",
		unix: 1279207800,
	},
	{
		date: "1995-03-25 08:45:00",
		unix: 796121100,
	},
	{
		date: "2015-11-11 11:11:11",
		unix: 1447240271,
	},
}

func TestDatetimeToUnixFunc(t *testing.T) {
	for _, dtu := range DATE_TIMES_TO_UNIX {
		unix, err := DatetimeToUnix(dtu.date)

		if err != nil {
			t.Errorf("error when parsing datetime to unix: %v\n", err)
		}

		if unix != dtu.unix {
			t.Errorf("worng unix time: wanted %d, given %d\n", dtu.unix, unix)
		}
	}
}

func TestUpdateUnixTimestampFromDatetimeFunc(t *testing.T) {
	for _, dtu := range DATE_TIMES_TO_UNIX {
		p := &types.Price{
			Timestamp: types.Timestamp{
				Now:      false,
				Datetime: dtu.date,
				Unix:     0,
			},
		}

		err := UpdateUnixTimestampFromDatetime(p)
		if err != nil {
			t.Errorf("error when parsing datetime to unix: %v\n", err)
		}

		if (*p).Timestamp.Unix != dtu.unix {
			t.Errorf("worng unix time: wanted %d, given %d\n", dtu.unix, (*p).Timestamp.Unix)
		}
	}
}
