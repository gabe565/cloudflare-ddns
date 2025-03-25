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
	FlagDomain    = "domain"
	FlagInterval  = "interval"
	FlagDNSUseTCP = "dns-tcp"
	FlagProxied   = "proxied"
	FlagTTL       = "ttl"

	FlagCloudflareToken     = "cloudflare-token"
	FlagCloudflareKey       = "cloudflare-key"
	FlagCloudflareEmail     = "cloudflare-email"
	FlagCloudflareAccountID = "cloudflare-account-id"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	fs := cmd.Flags()

	fs.Var(&c.LogLevel, FlagLogLevel, "Log level (one of "+strings.Join(slogx.LevelStrings(), ", ")+")")
	fs.Var(&c.LogFormat, FlagLogFormat, "Log format (one of "+strings.Join(slogx.FormatStrings(), ", ")+")")

	fs.StringSliceVar(&c.SourceStrs, FlagSource, c.SourceStrs, "Enabled IP sources (supports "+strings.Join(SourceStrings(), ", ")+")")
	fs.StringVar(&c.Domain, FlagDomain, c.Domain, "Domain to manage")
	fs.DurationVar(&c.Interval, FlagInterval, c.Interval, "Update interval")
	fs.BoolVar(&c.Proxied, FlagProxied, c.Proxied, "Enables Cloudflare proxy for the record")
	fs.Float64Var(&c.TTL, FlagTTL, c.TTL, "DNS record TTL (default auto)")
	fs.BoolVar(&c.DNSUseTCP, FlagDNSUseTCP, c.DNSUseTCP, "Force DNS to use TCP")

	fs.StringVar(&c.CloudflareToken, FlagCloudflareToken, c.CloudflareToken, "Cloudflare API token (recommended)")
	fs.StringVar(&c.CloudflareKey, FlagCloudflareKey, c.CloudflareKey, "Cloudflare API key")
	fs.StringVar(&c.CloudflareEmail, FlagCloudflareEmail, c.CloudflareEmail, "Cloudflare account email address")
	fs.StringVar(&c.CloudflareAccountID, FlagCloudflareAccountID, c.CloudflareAccountID, "Cloudflare account ID")
}
