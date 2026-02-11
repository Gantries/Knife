package eval

import (
	"bytes"
	"sync"
	"text/template"

	"github.com/gantries/knife/pkg/errors"
	"github.com/gantries/knife/pkg/maps"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Renderer struct {
	tpl *template.Template
}

var renderers = maps.Map[string, *Renderer]{}

var renderLock = sync.Mutex{}

func Template(tpl *string, functions template.FuncMap) (*Renderer, error) {
	renderLock.Lock()
	defer renderLock.Unlock()
	if _, ok := renderers[*tpl]; ok {
		return renderers[*tpl], nil
	}
	t, err := template.New("").Funcs(functions).Parse(*tpl)
	if err != nil {
		return nil, err
	}
	renderers[*tpl] = &Renderer{t}
	return renderers[*tpl], nil
}

func (r *Renderer) Render(tr *i18n.Localizer, vars any) (*string, error) {
	if r.tpl == nil {
		return nil, errors.MissingTemplateError.LocalE(tr, logger)
	}
	wr := bytes.Buffer{}
	if err := r.tpl.Execute(&wr, vars); err != nil {
		return nil, err
	}
	ps := wr.String()
	return &ps, nil
}
