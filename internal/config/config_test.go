package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	for _, fileName := range []string{"resend.conf", "sync.conf"} {
		cfg, err := NewConfig("../../example/" + fileName)
		assert.NoError(t, err)
		assert.Equal(t, time.Duration(10), cfg.Interval)
		assert.Equal(t, "*{{.Author.Username}}*: {{.Text}}", cfg.Template)
		assert.NotEmpty(t, cfg.Src)
		assert.NotEmpty(t, cfg.Dst)
	}
}
