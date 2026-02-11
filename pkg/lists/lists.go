package lists

import (
	"slices"
	"strings"
)

type List[T any] []T

func (l *List[T]) Length() int {
	return len(*l)
}

func (l *List[T]) Last() T {
	return (*l)[len(*l)-1]
}

func (l *List[T]) Add(a ...T) *List[T] {
	*l = append(*l, a...)
	return l
}

func (l *List[T]) Delete(a ...int) *List[T] {
	slices.Sort(a) // 1,6,3 -> 1,3,6
	var target List[T]
	last := 0
	length := len(*l)
	for _, pos := range a {
		if pos > length {
			break
		}
		target = append(target, (*l)[last:pos]...)
		// 0:1 2:3 4:6
		last = pos + 1
	}
	end := a[len(a)-1]
	if end < length-1 {
		// rest elements
		target = append(target, (*l)[end+1:]...)
	}
	*l = target
	return l
}

func (l *List[T]) Empty() bool {
	return len(*l) <= 0
}

func (l *List[T]) Sub(start, stop int) List[T] {
	if len(*l) <= start {
		return *Of[T]()
	}
	if len(*l) >= stop {
		return (*l)[start:stop]
	}
	return (*l)[start:]
}

func (l *List[T]) FirstOrDefault(d *T) *T {
	if l.Empty() {
		return d
	}
	return &(*l)[0]
}

func (l *List[T]) ForRest(n int, fn func(t T), rest func(t T)) {
	if len(*l) <= n {
		l.For(fn)
	}
	before := l.Sub(0, n)
	(&before).For(fn)
	after := l.Sub(n, l.Length())
	(&after).For(rest)
}

func FirstOrDefault[T any](a []T, d *T) *T {
	if len(a) <= 0 {
		return d
	}
	return &a[0]
}

func (l *List[T]) For(fn func(t T)) {
	for _, a := range *l {
		fn(a)
	}
}

func For[K any, V any](arr *[]K, m func(K) V) *[]V {
	var res []V
	for _, a := range *arr {
		res = append(res, m(a))
	}
	return &res
}

func ForArray[K any, V any](arr *List[K], m func(K) V) *[]V {
	var res []V
	for _, a := range *arr {
		res = append(res, m(a))
	}
	return &res
}

func ForEach[K interface{}](arr *[]*K, act func(*K)) {
	for _, a := range *arr {
		act(a)
	}
}

func ForFlatten[E any](arr *[]E, f func(E) []E) *[]E {
	var res []E
	for _, a := range *arr {
		res = append(res, f(a)...)
	}
	return &res
}

func Process[K any](arr *[]K, actors ...func(K) error) error {
	if len(*arr) <= 0 || len(actors) <= 0 {
		return nil
	}
	for _, e := range *arr {
		for _, a := range actors {
			err := a(e)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Collect[E any, O any](arr []E, f func(E) O) []O {
	var res []O
	for _, e := range arr {
		res = append(res, f(e))
	}
	return res
}

func Collects[E any, O any](arr []E, f func(E) []O) []O {
	var res []O
	for _, e := range arr {
		res = append(res, f(e)...)
	}
	return res
}

func Join[E any](arr *[]E, separator string, f func(E) string) *string {
	b := strings.Builder{}
	last := len(*arr) - 1
	for v := 0; v < last; v += 1 {
		b.WriteString(f((*arr)[v]))
		b.WriteString(separator)
	}
	if last > 0 {
		b.WriteString(f((*arr)[last]))
	}
	s := b.String()
	return &s
}

func VaJoin(separator string, arr ...string) *string {
	b := strings.Builder{}
	last := len(arr) - 1

	for v := 0; v < last; v += 1 {
		b.WriteString(arr[v])
		b.WriteString(separator)
	}
	if last > 0 {
		b.WriteString(arr[last])
	}

	s := b.String()
	return &s
}

func VaJoinFn[E any](separator string, f func(E) string, arr ...E) *string {
	b := strings.Builder{}
	last := len(arr) - 1

	for v := 0; v < last; v += 1 {
		b.WriteString(f((arr)[v]))
		b.WriteString(separator)
	}
	if last > 0 {
		b.WriteString(f((arr)[last]))
	}

	s := b.String()
	return &s
}

func Of[T any](a ...T) *[]T {
	return &a
}

func Filter[E any](arr []E, filter func(e E) bool) (List[E], List[E]) {
	matched := &List[E]{}
	dismiss := &List[E]{}
	for _, e := range arr {
		if filter(e) {
			matched.Add(e)
		} else {
			dismiss.Add(e)
		}
	}
	return *matched, *dismiss
}
