package config

import (
	"time"

	"gabe565.com/cloudflare-ddns/internal/lookup"
	"gabe565.com/utils/slogx"
	"github.com/cloudflare/cloudflare-go/v5"
	"github.com/cloudflare/cloudflare-go/v5/option"
	"github.com/cloudflare/cloudflare-go/v5/zones"
)

type Config struct {
	LogLevel  slogx.Level
	LogFormat slogx.Format

	SourceStrs []string
	UseV4      bool
	UseV6      bool
	Domains    []string
	Interval   time.Duration
	DNSUseTCP  bool
	Proxied    bool
	TTL        float64
	Timeout    time.Duration
	DryRun     bool

	CloudflareToken     string
	CloudflareKey       string
	CloudflareEmail     string
	CloudflareAccountID string
}

func New() *Config {
	return &Config{
		UseV4:   true,
		Timeout: time.Minute,
		SourceStrs: []string{
			lookup.CloudflareTLS.String(),
			lookup.OpenDNSTLS.String(),
			lookup.ICanHazIP.String(),
			lookup.IPInfo.String(),
			lookup.IPify.String(),
			lookup.Cloudflare.String(),
			lookup.OpenDNS.String(),
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

func (c *Config) CloudflareZoneListParams() zones.ZoneListParams {
	var params zones.ZoneListParams
	if c.CloudflareAccountID != "" {
		params.Account = cloudflare.F(zones.ZoneListParamsAccount{
			ID: cloudflare.F(c.CloudflareAccountID),
		})
	}
	return params
}

func (c *Config) Sources() ([]lookup.Source, error) {
	s := make([]lookup.Source, 0, len(c.SourceStrs))
	for _, str := range c.SourceStrs {
		source, err := lookup.SourceString(str)
		if err != nil {
			return nil, err
		}
		s = append(s, source)
	}
	return s, nil
}
