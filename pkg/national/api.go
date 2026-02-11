package national

import (
	"errors"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log/slog"
)

type Sentence string

var Ok Sentence = "Ok"

func (t Sentence) Error() string { return string(t) }

func (t Sentence) Raw() *string {
	s := string(t)
	return &s
}

func (t Sentence) Trans(tr *i18n.Localizer, args ...any) *string {
	var data = map[string]interface{}{}
	if len(args) > 0 && len(args)%2 != 0 {
		args = append(args, nil)
	}
	var length = len(args)
	for i := 0; i < length; i += 2 {
		if k, ok := args[i].(string); ok {
			data[k] = args[i+1]
		} else {
			data[fmt.Sprintf("%v", args[i])] = args[i+1]
		}
	}
	s, err := tr.Localize(&i18n.LocalizeConfig{
		MessageID:    string(t),
		TemplateData: data,
	})
	if err != nil {
		logger.Warn("Unable to translate message", "error", err, "message", t, "language", language.Und)
		s = string(t)
	}
	return &s
}

func (t Sentence) LocalE(tr *i18n.Localizer, log *slog.Logger, args ...any) error {
	err := fmt.Errorf(*t.Trans(tr, args...))
	if log != nil {
		log.Error(string(t)+": %s", append(args, err)...)
	}
	return err
}

func (t Sentence) E(log *slog.Logger, args ...any) error {
	err := fmt.Errorf(*t.Trans(En, args...))
	if log != nil {
		log.Error(string(t)+": %s", append(args, err)...)
	}
	return err
}

func (t Sentence) Fine() bool {
	return errors.Is(t, Ok)
}

func (t Sentence) Msg(argv ...any) *Message {
	return &Message{
		sentence: &t, args: &argv,
	}
}

func (t Sentence) Build(a ...any) *Message {
	return &Message{&t, &a}
}

func (t Sentence) Register() {
	err := bundle.AddMessages(language.Make(en), message(string(t), string(t), string(t)))
	if err != nil {
		logger.Error("Unable to register sentence as message", "error", err, "lang", en, "error", t)
		return
	}
}

var OkMessage = Message{&Ok, &[]any{}}

type Message struct {
	sentence *Sentence
	args     *[]any
}

func New(s Sentence, args ...any) *Message {
	return &Message{
		sentence: &s, args: &args,
	}
}

func (t *Message) String(tr *i18n.Localizer) *string {
	return t.sentence.Trans(tr, *t.args...)
}

func (t *Message) LocalE(tr *i18n.Localizer, log *slog.Logger) error {
	return t.sentence.LocalE(tr, log, *t.args...)
}

func (t *Message) E(log *slog.Logger) error {
	return t.sentence.LocalE(En, log, *t.args...)
}

func (t *Message) Fine() bool {
	return t.sentence.Fine()
}

func (t *Message) Body() *Sentence {
	return t.sentence
}
