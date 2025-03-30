package config

import (
	"strings"

	"gabe565.com/utils/slogx"
	"github.com/spf13/cobra"
)

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"

	FlagSource    = "source"
	FlagIPV4      = "ipv4"
	FlagIPV6      = "ipv6"
	FlagDomain    = "domain"
	FlagInterval  = "interval"
	FlagDNSUseTCP = "dns-tcp"
	FlagProxied   = "proxied"
	FlagTTL       = "ttl"
	FlagTimeout   = "timeout"
	FlagDryRun    = "dry-run"

	FlagCloudflareToken     = "cf-api-token" //nolint:gosec
	FlagCloudflareKey       = "cf-api-key"
	FlagCloudflareEmail     = "cf-api-email"
	FlagCloudflareAccountID = "cf-account-id"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	fs := cmd.Flags()

	fs.Var(&c.LogLevel, FlagLogLevel, "Log level (one of "+strings.Join(slogx.LevelStrings(), ", ")+")")
	fs.Var(&c.LogFormat, FlagLogFormat, "Log format (one of "+strings.Join(slogx.FormatStrings(), ", ")+")")

	fs.StringSliceVarP(&c.SourceStrs, FlagSource, "s", c.SourceStrs, "Enabled IP sources (supports "+strings.Join(SourceStrings(), ", ")+")")
	fs.BoolVarP(&c.UseV4, FlagIPV4, "4", c.UseV4, "Enables A records")
	fs.BoolVarP(&c.UseV6, FlagIPV6, "6", c.UseV6, "Enables AAAA records")
	fs.StringSliceVarP(&c.Domains, FlagDomain, "d", c.Domains, "Domains to manage")
	fs.DurationVarP(&c.Interval, FlagInterval, "i", c.Interval, "Update interval")
	fs.BoolVarP(&c.Proxied, FlagProxied, "p", c.Proxied, "Enables Cloudflare proxy for the record")
	fs.Float64VarP(&c.TTL, FlagTTL, "t", c.TTL, "DNS record TTL (default auto)")
	fs.BoolVar(&c.DNSUseTCP, FlagDNSUseTCP, c.DNSUseTCP, "Force DNS to use TCP")
	fs.DurationVar(&c.Timeout, FlagTimeout, c.Timeout, "Maximum length of time that an update may take")
	fs.BoolVarP(&c.DryRun, FlagDryRun, "n", c.DryRun, "Runs without changing any records")

	fs.StringVar(&c.CloudflareToken, FlagCloudflareToken, c.CloudflareToken, "Cloudflare API token (recommended)")
	fs.StringVar(&c.CloudflareKey, FlagCloudflareKey, c.CloudflareKey, "Cloudflare API key")
	fs.StringVar(&c.CloudflareEmail, FlagCloudflareEmail, c.CloudflareEmail, "Cloudflare account email address")
	fs.StringVar(&c.CloudflareAccountID, FlagCloudflareAccountID, c.CloudflareAccountID, "Cloudflare account ID")
}
