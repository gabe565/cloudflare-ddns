package config

import (
	"time"

	"gabe565.com/utils/slogx"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
)

type Config struct {
	LogLevel  slogx.Level
	LogFormat slogx.Format

	Sources   Sources
	Domain    string
	Interval  time.Duration
	DNSUseTCP bool
	Proxied   bool
	TTL       float64

	CloudflareToken string
	CloudflareKey   string
	CloudflareEmail string
}

func New() *Config {
	return &Config{
		Sources: Sources{
			SourceCloudflareTLS,
			SourceOpenDNSTLS,
			SourceCloudflare,
			SourceOpenDNS,
			SourceIPInfo,
			SourceIPify,
		},
	}
}

func (c *Config) NewCloudflareClient() (*cloudflare.Client, error) {
	var opts []option.RequestOption
	switch {
	case c.CloudflareToken != "":
		opts = append(opts, option.WithAPIToken(c.CloudflareToken))
	case c.CloudflareEmail != "" && c.CloudflareKey != "":
		opts = append(opts,
			option.WithAPIEmail(c.CloudflareEmail),
			option.WithAPIKey(c.CloudflareKey),
		)
	}

	return cloudflare.NewClient(opts...), nil
}
