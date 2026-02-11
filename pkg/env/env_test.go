package env

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEnviron(t *testing.T) {
	e := Environ()
	assert.True(t, len(e) > 0)
}
