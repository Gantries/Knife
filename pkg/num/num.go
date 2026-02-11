package num

import (
	"strconv"

	"github.com/gantries/knife/pkg/errors"
)

func String(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Long(s *string, base int, bitSize int) (i int64, err error) {
	if s == nil {
		return 0, errors.MissingValueError.E(logger)
	}
	return strconv.ParseInt(*s, base, bitSize)
}

func Float64(s *string) (d float64, err error) {
	if s == nil {
		return 0, errors.MissingValueError.E(logger)
	}
	return strconv.ParseFloat(*s, 64)
}
