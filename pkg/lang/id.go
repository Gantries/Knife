package lang

import (
	"time"

	"github.com/google/uuid"
)

// You can note that numbers follow each other: 1 (January), 2 (day), 3 (hour), 4(minutes)…
const format = "20060102030405"

func StringNow() string {
	return time.Now().Format(format)
}

func StringTs(t time.Time) string {
	return t.Format(format)
}

func StringUUID() string {
	return uuid.New().String()
}
