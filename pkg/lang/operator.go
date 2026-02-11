package lang

func If[T any, V any](expression bool, tf func(T) V, ff func(T) V, arg T) V {
	if expression {
		return tf(arg)
	} else {
		return ff(arg)
	}
}

func Ternary[T any](expression bool, t T, f T) T {
	if expression {
		return t
	} else {
		return f
	}
}

func Default[T any](v *T, d T) T {
	if v == nil {
		return d
	}
	return *v
}

func ComputeIf[T interface{}](expression bool, tf func() T, ff func() T) T {
	if expression {
		return tf()
	}
	return ff()
}
