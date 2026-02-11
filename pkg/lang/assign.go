package lang

func Assignable[T any](v any) bool {
	if _, ok := v.(T); ok {
		return true
	}
	return false
}
