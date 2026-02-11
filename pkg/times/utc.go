package times

import (
	"time"

	"github.com/gantries/knife/pkg/types"
)

var UTCOffset int64 = 28800000000

// FormatTimestamp 将Unix时间戳格式化为指定的UTC时间字符串, 单位: ms
func FormatTimestamp(timestamp int64, format string) string {
	// 将时间戳转换为time.Time对象
	t := time.Unix(timestamp/int64(time.Microsecond), 0).UTC()

	// 使用指定的格式格式化时间
	formattedTime := t.Format(format)

	return formattedTime
}

// FixedTsByLocation adjusts the provided time.Time to a fixed timestamp
// by subtracting the microseconds offset of the given location from the Unix
// microsecond timestamp. This function is necessary because GORM returns all
// time.Time values in UTC, and this adjustment allows for converting them to
// the correct local time zone string, especially for database drivers that do
// not handle timezone conversions automatically.
func FixedTsByLocation(location *time.Location, time time.Time) types.IntCount {
	offset := getOffsetMicroseconds(location)
	ts := time.UnixMicro() - offset
	return types.IntCount(ts)
}

// FixedTs is a convenience function that adjusts the provided time.Time
// to a fixed timestamp by subtracting the predefined UTC offset (UTCOffset).
// Similar to FixedTsByLocation, this function addresses the issue of
// GORM returning time.Time values in UTC, which may not be correctly converted
// to local time zone strings by some database drivers. This version uses a
// predefined UTC offset, assuming the caller knows the required offset.
func FixedTs(time time.Time) types.IntCount {
	ts := time.UnixMicro() - UTCOffset
	return types.IntCount(ts)
}

// getOffsetMicroseconds calculates the difference in micro seconds between the current timezone and UTC.
func getOffsetMicroseconds(location *time.Location) int64 {
	// Get the time in both time zones
	timeInLoc1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	timeInLoc2 := time.Date(2024, 1, 1, 0, 0, 0, 0, location)

	// Calculate the difference in micro seconds
	return timeInLoc1.Sub(timeInLoc2).Microseconds()
}
