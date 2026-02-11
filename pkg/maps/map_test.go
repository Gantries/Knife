package maps

import (
	"strconv"
	"strings"
	"testing"

	"github.com/gantries/knife/pkg/lists"

	"github.com/stretchr/testify/assert"
)

func TestMap_Equals(t *testing.T) {
	v1 := Map[string, int]{"a": 1, "b": 2}
	v2 := Map[string, int]{"a": 1, "b": 2}
	assert.True(t, v1.Equals(v2, func(left, right int) bool {
		return left == right
	}))
}

func TestKeys(t *testing.T) {
	v := map[string]int{"a": 1, "b": 2}
	s := Map[string, int](v)
	assert.True(t, s["a"] == 1)
	assert.True(t, (Of[string, int]("a", 1, "b", 2))["a"] == 1)
	assert.True(t, (Of[string, int]("a", 1, "b"))["a"] == 1)
	assert.True(t, (Of[string, int]("a", 1, "b"))["b"] == 0)
	assert.True(t, (Of[string, interface{}]("a", 1, "b"))["b"] == nil)
	// Map key ordering is not deterministic, so check if the key exists in the result
	keys := s.Keys()
	assert.True(t, len(keys) == 2)
	assert.True(t, (keys[0] == "a" && keys[1] == "b") || (keys[0] == "b" && keys[1] == "a"))
	// Values order depends on keys order, so we check the total length instead
	values := s.Values()
	assert.True(t, len(values) == 2)
}

func TestSet(t *testing.T) {
	has := func(k string, s Set[string]) bool {
		return (s).Has(k)
	}
	v := map[string]int{"a": 1, "b": 2}
	vc := Map[string, int](v)
	assert.True(t, has("a", &vc))

	one, two := 1, 2
	p := map[string]*int{"a": &one, "b": &two}
	pc := Map[string, *int](p)
	assert.True(t, has("a", &pc))
}

func TestFrom(t *testing.T) {
	a := []int{1, 2, 3, 4}
	v := From(a, strconv.Itoa)
	assert.True(t, len(v) == len(a))
	assert.True(t, v["1"] == 1)
	f := FromFn(a, strconv.Itoa, func(i int) int64 { return int64(i) })
	assert.True(t, f["1"] == int64(1))
}

func TestFromFilter(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	m, d := FromFilter(a, strconv.Itoa, func(i int) bool { return i%2 == 1 })
	assert.Equal(t, len(m), 3)
	assert.Equal(t, len(d), 2)
	assert.Equal(t, m["1"], 1)
	assert.Equal(t, m["3"], 3)
	assert.Equal(t, m["5"], 5)
	assert.Equal(t, d["2"], 2)
	assert.Equal(t, d["4"], 4)
}

func TestGetter(t *testing.T) {
	a := Map[string, int]{
		"a": 1,
		"b": 2,
	}

	getter := Getter[string, int]("a")
	vl := getter(a)
	assert.True(t, vl == 1)
}

func TestMap_Put(t *testing.T) {
	a := Map[string, int]{
		"a": 1,
		"b": 2,
	}
	a.Put("c", 3)
	assert.True(t, a["c"] == 3)
}

func TestBatch(t *testing.T) {
	a := Map[string, int]{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}
	v, err := Batch(a, 2, func(m Map[string, int]) ([]float64, error) {
		a := m.Values()
		return *lists.For(&a, func(v int) float64 {
			return float64(v * 2.0)
		}), nil
	})
	assert.Nil(t, err)
	assert.True(t, len(v) == len(a))
}

func TestMap_Merge(t *testing.T) {
	old := Map[string, Map[string, int]]{
		"a": Map[string, int]{
			"a": 1,
			"b": 2,
		},
		"b": Map[string, int]{
			"c": 3,
			"d": 4,
		},
	}
	updates := Map[string, Map[string, int]]{
		"c": Map[string, int]{
			"a": 1,
			"b": 2,
		},
	}
	merged, err := old.Merge(func(o, n Map[string, int]) (Map[string, int], error) {
		return o.PutAll(n), nil
	}, updates)
	assert.Nil(t, err)
	assert.True(t, (merged)["c"]["a"] == 1)
}

func TestMap_PutAll(t *testing.T) {
	a := Map[string, int]{
		"a": 1,
		"b": 2,
	}
	a.PutAll(Map[string, int]{"c": 3, "d": 4})
	assert.True(t, a["c"] == 3)
}

func Test_ByRef(t *testing.T) {
	m := Map[string, int]{
		"a": 1,
		"b": 2,
		"c": 3,
	}
	f := func(m Map[string, int]) {
		m.Put("a", 100)
	}
	f(m)
	assert.True(t, m["a"] == 100)
}

func TestVisitor(t *testing.T) {
	m := Map[string, interface{}]{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": 1,
			},
		},
		"b": 2,
		"c": 3,
	}
	Visit(m, nil,
		func(parent *string, key string) *string {
			if parent == nil {
				return &key
			} else {
				p := strings.Join([]string{*parent, key}, "/")
				return &p
			}
		},
		func(value interface{}) (map[string]interface{}, bool) {
			if c, ok := value.(map[string]interface{}); ok {
				return c, ok
			}
			return nil, false
		},
		func(parent *string, key string, value interface{}) {
			var path string
			if parent == nil {
				path = key
			} else {
				path = strings.Join([]string{*parent, key}, "/")
			}
			t.Logf("path = %s, value = %x\n", path, value)
		},
	)
}
