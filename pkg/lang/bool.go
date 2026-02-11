package lang

import (
	"strconv"

	"github.com/gantries/knife/pkg/errors"
)

func Bool(s *string) (i bool, err error) {
	if s == nil {
		return false, errors.MissingValueError.E(logger)
	}
	return strconv.ParseBool(*s)
}
