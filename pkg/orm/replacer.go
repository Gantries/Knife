package orm

import (
	"gorm.io/gorm/schema"
	"strings"
)

type replacer struct {
	properties DatabaseProperties
}

func newReplacer(properties DatabaseProperties) schema.Replacer {
	return replacer{properties}
}

func (r replacer) Replace(name string) string {
	if r.properties.GetNoLowerCase() {
		return strings.ToUpper(name)
	}
	return name
}
