package config

import (
	"log/slog"
	"slices"
	"strings"

	"gabe565.com/utils/must"
	"gabe565.com/utils/slogx"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/accounts"
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
	must.Must(cmd.RegisterFlagCompletionFunc(FlagDNSUseTCP, completeBool))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagProxied, completeBool))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagTTL, func(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return []string{"0\tauto", "5m", "1h"}, cobra.ShellCompDirectiveNoFileComp | cobra.ShellCompDirectiveKeepOrder
	}))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagTimeout, func(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
		return []string{"0\tno timeout", "30s", "1m"}, cobra.ShellCompDirectiveNoFileComp
	}))

	must.Must(cmd.RegisterFlagCompletionFunc(FlagCloudflareAccountID, completeAccount))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagCloudflareToken, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagCloudflareKey, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagCloudflareEmail, cobra.NoFileCompletions))
}

func completeBool(_ *cobra.Command, _ []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
	return []string{"true", "false"}, cobra.ShellCompDirectiveNoFileComp
}

func setupCompletion(cmd *cobra.Command, args []string) (*Config, *cloudflare.Client, error) {
	conf, err := Load(cmd, args)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return nil, nil, err
	}

	client, err := conf.NewCloudflareClient()
	if err != nil {
		slog.Error("Failed to create Cloudflare client", "error", err)
		return nil, nil, err
	}

	return conf, client, nil
}

func CompleteDomain(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
	conf, client, err := setupCompletion(cmd, args)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	trimmed, _, _ := strings.Cut(toComplete, ".")
	var domains []string
	iter := client.Zones.ListAutoPaging(cmd.Context(), conf.CloudflareZoneListParams())
	for iter.Next() {
		name := iter.Current().Name
		if toComplete != "" {
			if strings.HasSuffix(toComplete, name) {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			name = trimmed + "." + name
		}

		if !slices.Contains(conf.Domains, name) {
			domains = append(domains, name)
		}
	}
	if err := iter.Err(); err != nil {
		slog.Error("Failed to list zones", "error", err)
		return nil, cobra.ShellCompDirectiveError
	}

	return domains, cobra.ShellCompDirectiveNoFileComp
}

func completeAccount(cmd *cobra.Command, args []string, _ string) ([]cobra.Completion, cobra.ShellCompDirective) {
	_, client, err := setupCompletion(cmd, args)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var names []string
	iter := client.Accounts.ListAutoPaging(cmd.Context(), accounts.AccountListParams{})
	for iter.Next() {
		account := iter.Current()
		names = append(names, account.ID+"\t"+account.Name)
	}
	if err := iter.Err(); err != nil {
		slog.Error("Failed to list accounts", "error", err)
		return nil, cobra.ShellCompDirectiveError
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}
