package config

import (
	"slices"

	"gabe565.com/utils/must"
	"gabe565.com/utils/slogx"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/spf13/cobra"
)

func (c *Config) RegisterCompletions(cmd *cobra.Command) {
	must.Must(cmd.RegisterFlagCompletionFunc(FlagLogLevel, func(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return slogx.LevelStrings(), cobra.ShellCompDirectiveNoFileComp
	}))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagLogFormat, func(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return slogx.FormatStrings(), cobra.ShellCompDirectiveNoFileComp
	}))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagSource, func(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return SourceStrings(), cobra.ShellCompDirectiveNoFileComp
	}))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagDomain, CompleteDomain))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagInterval, func(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return []string{"1m", "15m", "1h", "24h"}, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveKeepOrder
	}))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagTTL, func(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return []string{"0\tauto", "5m", "1h"}, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveKeepOrder
	}))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagCloudflareToken, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagCloudflareKey, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagCloudflareEmail, cobra.NoFileCompletions))
}

func CompleteDomain(cmd *cobra.Command, args []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
	conf, err := Load(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	client, err := conf.NewCloudflareClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var domains []string
	iter := client.Zones.ListAutoPaging(cmd.Context(), zones.ZoneListParams{})
	for iter.Next() {
		name := iter.Current().Name
		if !slices.Contains(args, name) {
			domains = append(domains, name)
		}
	}
	if err := iter.Err(); err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	return domains, cobra.ShellCompDirectiveNoFileComp
}
