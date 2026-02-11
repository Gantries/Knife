package lang

import (
	"testing"

	"github.com/gantries/knife/pkg/lists"
	"github.com/stretchr/testify/assert"
)

func TestIf(t *testing.T) {
	tf := func(a []string) bool { return true }
	ff := func(a []string) bool { return false }
	assert.True(t, If(true, tf, ff, *lists.Of("hello", "world")))
	assert.False(t, If(false, tf, ff, *lists.Of("hello")))
}
