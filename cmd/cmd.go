package cmd

import (
	"context"
	"errors"
	"fmt"
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
		Short: "Sync a Cloudflare DNS record with your public IP address",
		RunE:  run,
		Args:  cobra.MaximumNArgs(1),

		SilenceErrors: true,
	}

	conf := config.New()
	conf.RegisterFlags(cmd)

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

var (
	ErrCloudflareAuth = errors.New("missing Cloudflare auth")
	ErrDomainRequired = errors.New("domain is required")
)

func run(cmd *cobra.Command, args []string) error {
	conf, err := config.Load(cmd)
	if err != nil {
		return err
	}

	switch {
	case conf.CloudflareToken == "" && conf.CloudflareKey == "":
		return fmt.Errorf("%w: CF_API_KEY or CF_API_TOKEN is required", ErrCloudflareAuth)
	case conf.CloudflareKey != "" && conf.CloudflareEmail == "":
		return fmt.Errorf("%w: CF_API_EMAIL is required", ErrCloudflareAuth)
	}

	if len(args) != 0 {
		conf.Domain = args[0]
	}
	if conf.Domain == "" {
		return ErrDomainRequired
	}

	cmd.SilenceUsage = true

	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	if err := ddns.Update(ctx, conf, false); err != nil {
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
				if err := ddns.Update(ctx, conf, true); err != nil {
					slog.Error("Run failed", "error", err)
				}
			}
		}
	}

	return nil
}
