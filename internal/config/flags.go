package config

import (
	"strings"

	"gabe565.com/utils/slogx"
	"github.com/spf13/cobra"
)

const (
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"

	FlagSources   = "sources"
	FlagDomain    = "domain"
	FlagInterval  = "interval"
	FlagDNSUseTCP = "dns-tcp"
	FlagProxied   = "proxied"
	FlagTTL       = "ttl"

	FlagCloudflareToken = "cloudflare-token"
	FlagCloudflareKey   = "cloudflare-key"
	FlagCloudflareEmail = "cloudflare-email"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	fs := cmd.Flags()

	fs.Var(&c.LogLevel, FlagLogLevel, "Log level (one of "+strings.Join(slogx.LevelStrings(), ", ")+")")
	fs.Var(&c.LogFormat, FlagLogFormat, "Log format (one of "+strings.Join(slogx.FormatStrings(), ", ")+")")

	fs.Var(&c.Sources, FlagSources, "Enabled IP sources (supports "+strings.Join(SourceStrings(), ", ")+")")
	fs.StringVar(&c.Domain, FlagDomain, c.Domain, "Domain to manage")
	fs.DurationVar(&c.Interval, FlagInterval, c.Interval, "Update interval")
	fs.BoolVar(&c.Proxied, FlagProxied, c.Proxied, "Enables Cloudflare proxy for the record")
	fs.Float64Var(&c.TTL, FlagTTL, c.TTL, "DNS record TTL (default auto)")
	fs.BoolVar(&c.DNSUseTCP, FlagDNSUseTCP, c.DNSUseTCP, "Force DNS to use TCP")

	fs.StringVar(&c.CloudflareToken, FlagCloudflareToken, c.CloudflareToken, "Cloudflare token")
	fs.StringVar(&c.CloudflareKey, FlagCloudflareKey, c.CloudflareKey, "Cloudflare API key")
	fs.StringVar(&c.CloudflareEmail, FlagCloudflareEmail, c.CloudflareEmail, "Cloudflare email")
}
