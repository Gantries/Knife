package maps

import (
	"strings"

	"github.com/gantries/knife/pkg/lists"
)

type Map[K comparable, V interface{}] map[K]V

func (m Map[K, V]) Has(k K) bool {
	if _, ok := m[k]; ok {
		return true
	}
	return false
}

// Get should be used cautiously
func (m Map[K, V]) Get(k K) *V {
	if r, ok := m[k]; ok {
		return &r
	}
	return nil
}

// GetOrDefault should be used cautiously
func (m Map[K, V]) GetOrDefault(k K, def V) *V {
	if r, ok := m[k]; ok {
		return &r
	}
	return &def
}

func (m Map[K, V]) Length() int {
	return len(m)
}

func (m Map[K, V]) Merge(merge func(o, n V) (V, error), a ...Map[K, V]) (Map[K, V], error) {
	for _, i := range a {
		for k, v := range i {
			if old, ok := m[k]; !ok {
				m[k] = v
			} else {
				r, err := merge(old, v)
				if err != nil {
					return nil, err
				}
				m[k] = r
			}
		}
	}
	return m, nil
}

func (m Map[K, V]) Put(k K, v V) Map[K, V] {
	m[k] = v
	return m
}

func (m Map[K, V]) PutAll(a ...Map[K, V]) Map[K, V] {
	return m.PutAllWith(func(o, n V) V { return n }, a...)
}

func (m Map[K, V]) PutIfAbsent(k K, f func(k K) V) V {
	if _, ok := m[k]; !ok {
		m[k] = f(k)
	}
	return m[k]
}

func (m Map[K, V]) PutAllWith(f func(o, n V) V, a ...Map[K, V]) Map[K, V] {
	for _, i := range a {
		for k, v := range i {
			if old, ok := m[k]; ok {
				m[k] = f(old, v)
			} else {
				m[k] = v
			}
		}
	}
	return m
}

func (m Map[K, V]) Keys() []K {
	r := lists.List[K]{}
	for k := range m {
		r.Add(k)
	}
	return r
}

func (m Map[K, V]) Values() []V {
	r := lists.List[V]{}
	for _, v := range m {
		r.Add(v)
	}
	return r
}

func (m Map[K, V]) Equals(t Map[K, V], eq func(left, right V) bool) bool {
	if len(m) != len(t) {
		return false
	}
	for k, v := range m {
		if r, ok := t[k]; ok {
			if !eq(v, r) {
				return false
			}
		}
	}
	return true
}

func (m Map[K, V]) String(ks func(k K) string, vs func(v V) string, kvsep, sep string) string {
	builder := strings.Builder{}
	for k, v := range m {
		if builder.Len() > 0 {
			builder.WriteString(sep)
		}
		builder.WriteString(ks(k))
		builder.WriteString(kvsep)
		builder.WriteString(vs(v))
	}
	return builder.String()
}

func Visit[K comparable, V interface{}, P interface{}](m map[K]V, parent *P, tracer func(parent *P, key K) *P, drill func(value V) (map[K]V, bool), action func(parent *P, key K, value V)) {
	for k, v := range m {
		if c, ok := drill(v); ok {
			Visit(c, tracer(parent, k), tracer, drill, action)
		} else {
			action(parent, k, v)
		}
	}
}

func Composite[K comparable, V, O any](m Map[K, V], fn func(k K, v V) O) lists.List[O] {
	res := lists.List[O]{}
	for k, v := range m {
		res.Add(fn(k, v))
	}
	return res
}

func FromFn[K comparable, T, O any](arr []T, key func(T) K, val func(T) O) Map[K, O] {
	res := Map[K, O]{}
	for _, v := range arr {
		res[key(v)] = val(v)
	}
	return res
}

func From[K comparable, T any](arr []T, key func(T) K) Map[K, T] {
	return FromFn[K, T, T](arr, key, func(v T) T { return v })
}

func FromFilter[K comparable, T any](arr []T, key func(T) K, filter func(T) bool) (Map[K, T], Map[K, T]) {
	matched, dismiss := Map[K, T]{}, Map[K, T]{}
	for _, e := range arr {
		if filter(e) {
			matched[key(e)] = e
		} else {
			dismiss[key(e)] = e
		}
	}
	return matched, dismiss
}

func Getter[K comparable, V any](field K) func(Map[K, V]) V {
	return func(m Map[K, V]) V {
		return m[field]
	}
}

func Of[K comparable, V any](a ...any) Map[K, V] {
	var res = Map[K, V]{}
	var length = len(a)

	last := length - 1 // if elements count is
	if length%2 != 0 {
		var v V
		res[a[last].(K)] = v
		last = length - 2 // skip last element
	}

	for i := 0; i < last; i += 2 {
		res[a[i].(K)] = a[i+1].(V)
	}

	return res
}

func Batch[K comparable, V any, T any](m Map[K, V], n int, fn func(Map[K, V]) ([]T, error)) ([]T, error) {
	rows := make([]T, 0)
	batch := Map[K, V]{}
	i := 0
	for k, v := range m {
		batch[k] = v
		i += 1
		if i%n == 0 {
			a, e := fn(batch)
			if e != nil {
				return rows, e
			}
			rows = append(rows, a...)
			batch = Map[K, V]{}
		}
	}
	if i%n != 0 {
		a, e := fn(batch)
		if e != nil {
			return rows, e
		}
		rows = append(rows, a...)
	}
	return rows, nil
}

func Default[K comparable, V any](v, def Map[K, V]) Map[K, V] {
	if v == nil {
		return def
	}
	return v
}
