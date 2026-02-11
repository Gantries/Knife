package lang

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDuplicate(t *testing.T) {
	v := 1
	pv := Dup(v)
	assert.False(t, &v == pv, "Unexpected equal")
}
