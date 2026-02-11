package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type typedInt int64

func Test_ToLong(t *testing.T) {
	tests := []interface{}{
		int(1), int8(1), int16(1), int32(1), int64(1),
		uint(1), uint8(1), uint16(1), uint32(1), uint64(1),
		"1", []byte{1}, byte(1)}
	def := typedInt(0)
	for i, a := range tests {
		v := *ToLong(a, func(i *int64) *typedInt { t := typedInt(*i); return &t }, &def)
		assert.Equal(t, v, typedInt(1), "case %d", i)
	}
	assert.Equal(t, *ToLong([]byte{1, 1}, func(i *int64) *typedInt { t := typedInt(*i); return &t }, &def),
		typedInt(257), "case []byte")
}
