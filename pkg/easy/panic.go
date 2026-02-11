package easy

import "fmt"

func Panic(err error, cleanup func(error)) {
	if err != nil {
		if cleanup != nil {
			cleanup(err)
		}
		panic(err)
	}
}

func PanicE[T any](v *T, err error) *T {
	if err != nil {
		panic(err)
	}
	return v
}

func PanicN[T any](v *T, cleanup func(error)) *T {
	if v != nil {
		return v
	}
	err := fmt.Errorf("expected value got nil")
	if cleanup != nil {
		cleanup(err)
	}
	panic(err)
}
