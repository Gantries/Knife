package lists

import (
	"strconv"
	"testing"

	"github.com/gantries/knife/pkg/serde"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	v := []int64{1, 2, 3, 4}
	s := *For(&v, func(t int64) string {
		return strconv.FormatInt(t, 10)
	})
	assert.True(t, len(s) == len(v))
	assert.True(t, s[0] == "1")
}

type Counter struct{ counter int }

func TestEach(t *testing.T) {
	v := []*Counter{
		{1}, {2}, {3},
	}
	ForEach(&v, func(c *Counter) { c.counter += 1 })
	assert.True(t, v[0].counter == 2)
	assert.True(t, v[2].counter == 4)
}

func TestCollect(t *testing.T) {
	v := []*Counter{
		{1}, {2}, {3},
	}
	l := Collect(v, func(e *Counter) int {
		return e.counter
	})
	assert.True(t, len(l) == len(v))
	assert.True(t, (l)[0] == 1)
}

func TestJoin(t *testing.T) {
	l := []int{1, 2, 3, 4, 5}
	assert.True(t, *Join(&l, ",", strconv.Itoa) == "1,2,3,4,5")
	assert.True(t, *VaJoin(",", "1", "2", "3") == "1,2,3")
	assert.True(t, *VaJoinFn(",", strconv.Itoa, 1, 2, 3) == "1,2,3")
}

func TestList(t *testing.T) {
	l := List[int]{1, 2, 3, 4, 5, 6}
	l.Delete(4, 2)
	assert.True(t, l.Length() == 4)
	assert.True(t, l[2] == 4)
	assert.True(t, l[3] == 6)
	l.Add(7)
	assert.True(t, l.Last() == 7)
}

func TestFor(t *testing.T) {
	l := List[int]{1, 2, 3, 4, 5, 6}
	a := ForArray(&l, strconv.Itoa)
	assert.True(t, (*a)[0] == "1")
}

func TestForFlatten(t *testing.T) {
	l := []int{1, 2, 3, 4, 5, 6}
	a := ForFlatten(&l, func(i int) []int {
		return []int{i, i + 1}
	})
	buf, e := serde.Serialize(*a)
	assert.Nil(t, e)
	assert.True(t, string(buf) == "[1,2,2,3,3,4,4,5,5,6,6,7]")
}

func TestForFilter(t *testing.T) {
	l := List[int]{1, 2, 3, 4, 5, 6, 7, 8, 9}
	m, d := Filter(l, func(i int) bool { return i%2 == 1 })
	assert.Equal(t, m[0], 1)
	assert.Equal(t, m[4], 9)
	assert.Equal(t, d[0], 2)
	assert.Equal(t, d[3], 8)
	assert.Equal(t, len(m), 5)
	assert.Equal(t, len(d), 4)
}
