package cmd

import (
	"context"

	"gabe565.com/cloudflare-ddns/internal/output"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra"
)

func WithMarkdown() cobrax.Option {
	return func(cmd *cobra.Command) {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		cmd.SetContext(output.NewContext(ctx, output.FormatMarkdown))
	}
}
