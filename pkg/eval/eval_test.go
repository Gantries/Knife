package eval

import (
	"testing"

	"github.com/gantries/knife/pkg/maps"
	"github.com/gantries/knife/pkg/national"
	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	tr := national.En
	tpl := `"hello world"`
	r, err := Evaluate(tr, &tpl, maps.Map[string, interface{}]{})
	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "hello world", *r)
}

func TestEvaluateStringConcat(t *testing.T) {
	tr := national.En
	tpl := `"hello" + " " + "world"`
	r, err := Evaluate(tr, &tpl, maps.Map[string, interface{}]{})
	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "hello world", *r)
}

func TestEvaluateWithVars(t *testing.T) {
	tr := national.En
	tpl := `name`
	r, err := Evaluate(tr, &tpl, maps.Of[string, any]("name", "test user"))
	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "test user", *r)
}
