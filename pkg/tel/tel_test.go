package tel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Fingerprint(t *testing.T) {
	hostname, ip, mac, ns := Fingerprint()
	assert.NotEmpty(t, hostname)
	assert.NotEmpty(t, ip)
	assert.NotEmpty(t, mac)
	assert.Empty(t, ns)
}

func Test_BuiltinAttributeStrings(t *testing.T) {
	attrs := BuiltinAttributeStrings()
	assert.True(t, len(attrs)%2 == 0)
}
