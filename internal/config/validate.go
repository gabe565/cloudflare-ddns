package config

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidSource  = errors.New("invalid source")
	ErrNoProto        = errors.New("either v4 or v6 must be enabled")
	ErrCloudflareAuth = errors.New("missing Cloudflare auth")
	ErrNoDomain       = errors.New("at least one domain must be provided")
)

func (c *Config) Validate() error {
	for _, sourceStr := range c.SourceStrs {
		if _, err := SourceString(sourceStr); err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidSource, sourceStr)
		}
	}

	switch {
	case len(c.Domains) == 0:
		return ErrNoDomain
	case !c.UseV4 && !c.UseV6:
		return ErrNoProto
	case c.CloudflareToken == "" && c.CloudflareKey == "":
		return fmt.Errorf("%w: CF_API_KEY or CF_API_TOKEN is required", ErrCloudflareAuth)
	case c.CloudflareKey != "" && c.CloudflareEmail == "":
		return fmt.Errorf("%w: CF_API_EMAIL is required", ErrCloudflareAuth)
	}

	return nil
}
