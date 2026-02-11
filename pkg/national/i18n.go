package national

import (
	"context"
	"net/http"
	"sync"

	"github.com/gantries/knife/pkg/log"
	"github.com/gantries/knife/pkg/maps"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

var logger = log.New("knife/pkg/national")

type i18nLang string

var keyI18nLang = i18nLang("")

var bundle = i18n.NewBundle(language.English)

const en = "en"

var configurator I18nConfigurator = nil

func message(id, one, other string) *i18n.Message {
	return &i18n.Message{
		ID:    id,
		One:   one,
		Other: other,
	}
}

func LoadMessages(translations maps.Map[string, maps.Map[string, string]]) {
	messages := map[string][]*i18n.Message{}
	newMessage := func(lang string, m *i18n.Message) {
		messages[lang] = append(messages[lang], m)
	}
	for id, m := range translations {
		newMessage(en, message(id, id, id))
		for lang, trans := range m {
			newMessage(lang, message(id, trans, trans))
		}
	}
	for l, m := range messages {
		err := bundle.AddMessages(language.Make(l), m...)
		if err != nil {
			logger.Error("Unable to add messages", "error", err, "lang", l)
			return
		}
	}
}

func LoadMessagesFromString(cfg string) {
	var translations = maps.Map[string, maps.Map[string, string]]{}
	err := yaml.Unmarshal([]byte(cfg), &translations)
	if err != nil {
		logger.Error("Unable to parse translations", "error", err)
		return
	}
	LoadMessages(translations)
}

func WithLanguage(ctxt context.Context, req *http.Request, fallback string) context.Context {
	lang := req.Header.Get("Accept-Language")
	if len(lang) == 0 {
		lang = fallback
	}
	return context.WithValue(ctxt, keyI18nLang, lang)
}

func Language(ctxt context.Context) string {
	bundled := ctxt.Value(keyI18nLang)
	if bundled != nil {
		if l, ok := bundled.(string); ok {
			return l
		}
	}
	if configurator != nil {
		return configurator.DefaultLang()
	}
	return en
}

var translators = struct{ sync.Map }{}

func FindOrCreateLocalizer(lang string) *i18n.Localizer {
	f, ok := translators.Load(lang)
	if ok {
		return f.(*i18n.Localizer)
	}
	f, loaded := translators.LoadOrStore(lang, i18n.NewLocalizer(bundle, lang, en))
	if loaded {
		return f.(*i18n.Localizer)
	}
	return f.(*i18n.Localizer)
}

func Tr(ctxt context.Context) *i18n.Localizer {
	return FindOrCreateLocalizer(Language(ctxt))
}

var En = FindOrCreateLocalizer(en)

type I18nConfigurator interface {
	DefaultLang() string
}

func Prepare(c I18nConfigurator) {
	configurator = c
}
