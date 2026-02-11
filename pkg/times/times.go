// Package times provides time formatting and timezone utilities.
//
// It includes functions for formatting timestamps and calculating
// UTC offsets for local time zones.
package times

import "time"

func init() {
	UTCOffset = getOffsetMicroseconds(time.Local)
}

func FormatTs(ts int64) string {
	return FormatTsByLayout(ts, "2006-01-02 15:04:05")
}

func FormatTsByLayout(ts int64, layout string) string {
	localTime := time.Unix(ts/1000, 0)
	localTimeStr := localTime.Format(layout)
	return localTimeStr
}
