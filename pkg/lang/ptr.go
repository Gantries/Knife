package lang

func Dup[T any](v T) *T {
	return &v
}

func OrDefault[T any](v *T, d *T) *T {
	if v != nil {
		return v
	}
	return d
}
