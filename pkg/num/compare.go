package num

import (
	"strconv"

	"github.com/gantries/knife/pkg/log"
	"golang.org/x/exp/constraints"
)

var logger = log.New("knife/num")

type Converter[I any, O constraints.Integer | constraints.Float] interface {
	Convert(i I) (*O, error)
}

type Comparator[T any, O constraints.Integer | constraints.Float] interface {
	Between(min T, max T, v T, converter Converter[T, O]) (bool, error)
}

type GenericComparator[T any, O constraints.Integer | constraints.Float] struct{}

func (c *GenericComparator[T, O]) Between(min T, max T, v T, converter Converter[T, O]) (bool, error) {
	minValue, err := converter.Convert(min)
	if err != nil {
		logger.Error("Unable to convert minimum value", "error", err, "min", min)
		return false, err
	}
	maxValue, err := converter.Convert(max)
	if err != nil {
		logger.Error("Unable to convert maximum value", "error", err, "max", max)
		return false, err
	}
	if value, err := converter.Convert(v); err != nil {
		logger.Error("Unable to convert value", "error", err, "value", v)
		return false, err
	} else {
		return *minValue <= *value && *value < *maxValue, nil
	}
}

type StringFloat64Converter struct{}

func (c *StringFloat64Converter) Convert(s string) (*float64, error) {
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return &v, nil
	} else {
		return nil, err
	}
}
