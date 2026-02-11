package eval

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/gantries/knife/pkg/errors"
	"github.com/gantries/knife/pkg/lang"
	"github.com/gantries/knife/pkg/lists"
	"github.com/gantries/knife/pkg/log"
	"github.com/gantries/knife/pkg/maps"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var logger = log.New("knife/eval")

var templates = maps.Map[string, *vm.Program]{}

var exprLock = sync.Mutex{}

func compileExpressionOrGetFromCache(tpl *string) *vm.Program {
	exprLock.Lock()
	defer exprLock.Unlock()
	if p, ok := templates[*tpl]; ok {
		return p
	}
	p, err := expr.Compile(*tpl)
	if err != nil {
		logger.Error("Unable to compile expression", "error", err, "expression", tpl)
		return nil
	}
	templates[*tpl] = p
	return templates[*tpl]
}

var functions = maps.Map[string, any]{
	"fmt":   fmt.Sprintf,
	"log":   logger,
	"s2i":   maps.Of[string, any],
	"i2i":   maps.Of[any, any],
	"arr":   lists.Of[any],
	"p2i32": lang.OrDefault[int],
	"p2i":   lang.OrDefault[int64],
	"p2s":   lang.OrDefault[string],
	"p2a":   lang.OrDefault[any],
	"p2f32": lang.OrDefault[float32],
	"p2f":   lang.OrDefault[float64],
	"p2b":   lang.OrDefault[bool],
}

func Evaluate(tr *i18n.Localizer, tpl *string, vars maps.Map[string, interface{}],
	envs ...maps.Map[string, interface{}]) (*string, error) {
	program := compileExpressionOrGetFromCache(tpl)
	if program == nil {
		return nil, errors.CompileExpressionError.LocalE(tr, logger, "expression", tpl)
	}
	r, err := expr.Run(program, vars.PutAll(envs...).PutAll(functions))
	if err != nil {
		return nil, errors.EvaluateExpressionError.LocalE(tr, logger, "error", err)
	}
	output, ok := r.(string)
	if !ok {
		return nil, errors.ExpectedTypeButError.LocalE(tr, logger, "expected",
			"string", "actual", reflect.ValueOf(r))
	}
	return &output, nil
}
