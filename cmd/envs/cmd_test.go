package envs

import (
	"strings"
	"testing"

	"gabe565.com/cloudflare-ddns/internal/config"
	"gabe565.com/cloudflare-ddns/internal/output"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvs(t *testing.T) {
	cmd := New()
	conf := config.New()
	conf.RegisterFlags(cmd)

	t.Run("ansi", func(t *testing.T) {
		cmd.SetContext(output.NewContext(t.Context(), output.FormatANSI))
		var buf strings.Builder
		cmd.SetOut(&buf)
		require.NoError(t, cmd.Execute())
		assert.NotEmpty(t, buf.String())
	})

	t.Run("markdown", func(t *testing.T) {
		cmd.SetContext(output.NewContext(t.Context(), output.FormatMarkdown))
		var buf strings.Builder
		cmd.SetOut(&buf)
		require.NoError(t, cmd.Execute())
		assert.NotEmpty(t, buf.String())
	})
}
