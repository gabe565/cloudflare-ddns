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

	SourceStrs []string
	Domains    []string
	Interval   time.Duration
	DNSUseTCP  bool
	Proxied    bool
	TTL        float64
	Timeout    time.Duration

	CloudflareToken     string
	CloudflareKey       string
	CloudflareEmail     string
	CloudflareAccountID string
}

func New() *Config {
	return &Config{
		Timeout: time.Minute,
		SourceStrs: []string{
			SourceCloudflareTLS.String(),
			SourceOpenDNSTLS.String(),
			SourceIPInfo.String(),
			SourceIPify.String(),
			SourceCloudflare.String(),
			SourceOpenDNS.String(),
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

func (c *Config) Sources() ([]Source, error) {
	s := make([]Source, 0, len(c.SourceStrs))
	for _, str := range c.SourceStrs {
		source, err := SourceString(str)
		if err != nil {
			return nil, err
		}
		s = append(s, source)
	}
	return s, nil
}
