// Package errors provides internationalized error definitions and utilities.
//
// It uses the national package for i18n support and defines common error
// sentences that can be translated and formatted with context.
package errors

import (
	i "github.com/gantries/knife/pkg/national"
)

const (
	CompileExpressionError        i.Sentence = "Compile expression {{.express}} error"
	EvaluateExpressionError       i.Sentence = "Evaluate expression error: {{.error}}"
	ExpectedTypeButError          i.Sentence = "Type {{.expected}} is expected but got {{.actual}}"
	MissingTemplateError          i.Sentence = "Template is missing"
	MissingValueError             i.Sentence = "Missing required value"
	MissingAuthenticationToken    i.Sentence = "Missing authentication token"
	OverwriteInternalBuiltinError i.Sentence = "Internal builtin {{.builtin}} can't be overwritten"
	OverwriteBuiltinError         i.Sentence = "Builtin {{.builtin}} can't be overwritten"
	OverwriteIsForbiddenError     i.Sentence = "Overwrite {{.target}} of {{.type}} is not allowed"
	Unauthorized                  i.Sentence = "Unauthorized"
	UnexpectedTypeError           i.Sentence = "Got unexpected {{.type}}"
	UnexpectedValueError          i.Sentence = "Got unexpected {{.type}} {{.value}}"
	UnrecognizedError             i.Sentence = "Unrecognized {{.type}} {{.value}}"
	UnsupportedValueError         i.Sentence = "Unsupported {{.type}} {{.value}}"
	NotFoundError                 i.Sentence = "{{.type}} {{.value}} not found"
)

func Yes(actors ...func()) *i.Message {
	for _, act := range actors {
		act()
	}
	return &i.OkMessage
}

func No(e error, argv ...any) *i.Message {
	return i.New(i.Sentence(e.Error()), argv)
}

func init() {
	CompileExpressionError.Register()
	EvaluateExpressionError.Register()
	ExpectedTypeButError.Register()
	MissingValueError.Register()
	OverwriteInternalBuiltinError.Register()
	OverwriteBuiltinError.Register()
	OverwriteIsForbiddenError.Register()
	UnexpectedValueError.Register()
	UnrecognizedError.Register()
	UnsupportedValueError.Register()
	NotFoundError.Register()
}
