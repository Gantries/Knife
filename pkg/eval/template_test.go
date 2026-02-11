package eval

import (
	"testing"

	"github.com/gantries/knife/pkg/lists"
	"github.com/gantries/knife/pkg/maps"
	"github.com/gantries/knife/pkg/national"
	"github.com/stretchr/testify/assert"
)

func TestRenderer_Render(t *testing.T) {
	tpl := "hello {{.}}!"
	r, _ := Template(&tpl, nil)
	s, err := r.Render(national.En, "world")
	assert.Nil(t, err)
	assert.True(t, *s == "hello world!")

	tpl = `{{range .}}{{.Name}} {{.Type}} {{.Default}}{{end}}`
	r, _ = Template(&tpl, nil)
	s, err = r.Render(national.En, lists.Of[maps.Map[string, interface{}]](maps.Of[string, interface{}](
		"Name", "foo", "Type", "bar", "Default", "baz")))
	assert.Nil(t, err)
	assert.True(t, *s == "foo bar baz")
}
