package times

import (
	"testing"
	"time"
)

func TestFormatTimestamp(t *testing.T) {
	// 测试案例
	tests := []struct {
		name      string
		timestamp int64
		format    string
		expected  string
	}{
		{
			name:      "Normal Case",
			timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			format:    "2006-01-02 15:04:05",
			expected:  time.Now().UTC().Format("2006-01-02 15:04:05"), // 期望值是当前时间的UTC格式
		},
		{
			name:      "Epoch Start",
			timestamp: 0,
			format:    "2006-01-02 15:04:05",
			expected:  "1970-01-01 00:00:00",
		},
		{
			name:      "Future Time",
			timestamp: (time.Now().AddDate(10, 0, 0).UnixNano() / int64(time.Millisecond)) + 5000, // 未来的时间戳
			format:    "2006-01-02 15:04:05",
			expected:  time.Now().AddDate(10, 0, 0).Add(time.Duration(5000) * time.Millisecond).UTC().Format("2006-01-02 15:04:05"),
		},
		{
			name:      "Different Format",
			timestamp: 0,
			format:    "2006/01/02",
			expected:  "1970/01/01",
		},
		// 可以添加更多的测试案例来覆盖不同的场景
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formattedTime := FormatTimestamp(tt.timestamp, tt.format)
			if formattedTime == "" {
				t.Errorf("FormatTimestamp() is empty")
				return
			}
			if formattedTime != tt.expected {
				t.Errorf("FormatTimestamp() = %v, want %v", formattedTime, tt.expected)
			}
		})
	}
}

func TestSQLServerDatetime(t *testing.T) {
	ts := FormatTimestamp(1731542400000, "2006-01-02 15:04:05")
	println(ts)
}
