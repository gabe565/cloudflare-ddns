package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gabe565.com/cloudflare-ddns/internal/config"
	"gabe565.com/cloudflare-ddns/internal/ddns"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra"
)

func New(opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloudflare-ddns",
		Short: "Sync a Cloudflare DNS record with your current public IP address",
		RunE:  run,

		ValidArgsFunction: config.CompleteDomain,
		SilenceErrors:     true,
		DisableAutoGenTag: true,
	}

	conf := config.New()
	conf.RegisterFlags(cmd)
	conf.RegisterCompletions(cmd)

	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	cmd.SetContext(config.NewContext(ctx, conf))

	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	conf, err := config.Load(cmd, args)
	if err != nil {
		return err
	}

	if err := conf.Validate(); err != nil {
		return err
	}

	cmd.SilenceUsage = true

	if conf.Interval != 0 {
		slog.Info("Cloudflare DDNS", "version", cobrax.GetVersion(cmd), "commit", cobrax.GetCommit(cmd))
	}

	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	if err := ddns.Update(ctx, conf); err != nil {
		return err
	}

	if conf.Interval != 0 {
		ticker := time.NewTicker(conf.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				if err := ddns.Update(ctx, conf); err != nil {
					slog.Error("Run failed", "error", err)
				}
			}
		}
	}

	return nil
}
