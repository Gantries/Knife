package lang

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStringNow(t *testing.T) {
	n := StringTs(time.Now())
	assert.NotNil(t, n)
}
