package synch

import (
	"math"
	"reflect"
	"testing"

	"github.com/gantries/knife/pkg/lists"
	"github.com/gantries/knife/pkg/serde"
	"github.com/stretchr/testify/assert"
)

func TestRunner(t *testing.T) {
	repeat := 3
	o := Runner(repeat, func() int {
		repeat -= 1
		return repeat
	})
	assert.True(t, o.count == 3)

	for i := 0; i < 5; i += 1 {
		assert.True(t, o.count == 3-i)
		v := o.Run()
		assert.True(t, o.count == 3-i-1)
		if i < 3 {
			assert.True(t, v[0].Int() == int64(math.Max(float64(3-i-1), 0)))
		} else {
			assert.Nil(t, v)
		}
	}

	values := Runner(1, func(argv ...any) string {
		s, e := serde.Serialize(argv)
		assert.Nil(t, e)
		assert.True(t, string(s) == "[\"hello\",\"world\"]")
		return "nothing serious"
	}).Run(reflect.ValueOf("hello"), reflect.ValueOf("world"))
	texts := *lists.For(&values, func(v reflect.Value) string { return v.String() })
	s, e := serde.Serialize(texts)
	assert.Nil(t, e)
	assert.True(t, string(s) == "[\"nothing serious\"]")
}
