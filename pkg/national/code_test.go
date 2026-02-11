package national

import (
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
)

const translations = `
Missing table id {{.id}}:
  zh_CN: 表ID缺失{{.id}}
  ja_JP: テーブルIDがありません{{.id}}
`

var code = Sentence("Missing table id {{.id}}")

func TestTranslation(t *testing.T) {
	LoadMessagesFromString(translations)

	id := "0012"
	zh := i18n.NewLocalizer(bundle, "zh_CN", "ja_JP", "en")
	err := code.LocalE(zh, nil, "id", id)
	assert.True(t, err.Error() == "表ID缺失0012")
	ja := i18n.NewLocalizer(bundle, "ja_JP", "en")
	err = code.LocalE(ja, nil, "id", id)
	assert.Equal(t, err.Error(), "テーブルIDがありません0012")
	en := i18n.NewLocalizer(bundle, "en_US", "en")
	err = code.LocalE(en, nil, "id", id)
	assert.Equal(t, err.Error(), "Missing table id 0012")
}
