package nacos

import (
	"os"
	"strconv"
	"testing"

	"github.com/gantries/knife/pkg/assert"
	. "github.com/gantries/knife/pkg/dot"
	"github.com/gantries/knife/pkg/national"
)

func Test_SetupInternaltional(t *testing.T) {
	// Skip if NACOS_SERVER_ADDRESS environment variable is not set
	if _, exists := os.LookupEnv("NACOS_SERVER_ADDRESS"); !exists {
		t.Skip("NACOS_SERVER_ADDRESS environment variable not set, skipping Nacos tests")
		return
	}

	lang := "zh_CN"
	cfg := Configuration{
		ServerAddress: Must(Env("NACOS_SERVER_ADDRESS", MaxLineLength)),
		Timeout:       Must(strconv.Atoi(Must(Env("NACOS_TIME_OUT", MaxLineLength)))),
		Port:          Must(strconv.Atoi(Must(Env("NACOS_PORT", MaxLineLength)))),
		ContextPath:   Must(Env("NACOS_CONTEXT_PATH", MaxLineLength)),
		DataID:        Must(Env("NACOS_DATA_ID", MaxLineLength)),
		Group:         Must(Env("NACOS_GROUP", MaxLineLength)),
		Namespace:     Must(Env("NACOS_NAMESPACE", MaxLineLength)),
		Scheme:        Must(Env("NACOS_SCHEME", MaxLineLength)),
		Username:      Must(Env("NACOS_USERNAME", MaxLineLength)),
		Password:      Must(Env("NACOS_PASSWORD", MaxLineLength)),
		Languages: []Data{
			{
				Group: Must(Env("NACOS_DATA_0_GROUP", MaxLineLength)),
				ID:    Must(Env("NACOS_DATA_0_ID", MaxLineLength)),
			},
		},
		DefaultLanguage: lang,
	}
	w, err := NewWatcher(cfg)
	assert.True(t, err == nil, "no error expected")
	assert.True(t, w != nil, "watcher should be created")

	l := national.FindOrCreateLocalizer(lang)
	w.SetupInternational(cfg)
	s := national.Sentence("Showcase")
	e := s.LocalE(l, logger)
	assert.True(t, e.Error() == "演示", "i18n not work")
}
