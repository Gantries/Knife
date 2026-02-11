package cache

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Properties(t *testing.T) {
	p := Properties{}
	err := envconfig.Process("", &p)
	assert.Nil(t, err)
	assert.True(t, p.Pool.MaxConnectionIdleTime == time.Hour*1)
	assert.True(t, p.Pool.MaxConnectionLifeTime == time.Hour*10)
}
