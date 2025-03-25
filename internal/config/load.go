package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Load(cmd *cobra.Command) (*Config, error) {
	conf, ok := FromContext(cmd.Context())
	if !ok {
		panic("command missing config")
	}

	var errs []error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed {
			if val, ok := os.LookupEnv(EnvName(f.Name)); ok {
				if err := f.Value.Set(val); err != nil {
					errs = append(errs, err)
				}
			}
		}
	})
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	conf.InitLog(cmd.ErrOrStderr())

	return conf, nil
}

const EnvPrefix = "DDNS_"

func EnvName(name string) string {
	switch name {
	case FlagCloudflareToken:
		return "CF_API_TOKEN"
	case FlagCloudflareKey:
		return "CF_API_KEY"
	case FlagCloudflareEmail:
		return "CF_API_EMAIL"
	default:
		name = strings.ToUpper(name)
		name = strings.ReplaceAll(name, "-", "_")
		return EnvPrefix + name
	}
}
