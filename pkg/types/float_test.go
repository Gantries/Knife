package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloat_Value(t *testing.T) {
	f := Float(0)
	err := f.UnmarshalGQL(123)
	assert.Nil(t, err)
	assert.Nil(t, f.UnmarshalGQL("123.1"))
	assert.True(t, f == 123.1)
}
