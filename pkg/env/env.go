package env

import (
	"os"
	"strings"

	"github.com/gantries/knife/pkg/maps"
)

func Environ() maps.Map[string, string] {
	e := maps.Map[string, string]{}
	for _, v := range os.Environ() {
		pair := strings.SplitN(v, "=", 2)
		if len(pair) == 2 {
			e.Put(pair[0], pair[1])
		}
	}
	return e
}
