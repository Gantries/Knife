package num

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenericComparator_Between(t *testing.T) {
	cmp := GenericComparator[string, float64]{}
	cvt := &StringFloat64Converter{}
	r, err := cmp.Between("1.0", "2.0", "1.5", cvt)
	assert.Nil(t, err)
	assert.True(t, r)
	r, err = cmp.Between("1.0", "2.0", "3.0", cvt)
	assert.Nil(t, err)
	assert.False(t, r)
	_, err = cmp.Between("ab", "2.0", "3.0", cvt)
	assert.NotNil(t, err)
}
