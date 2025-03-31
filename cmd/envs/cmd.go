package envs

import (
	"cmp"
	"io"
	"slices"
	"strings"

	"gabe565.com/cloudflare-ddns/internal/config"
	"gabe565.com/cloudflare-ddns/internal/output"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const Name = "envs"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   Name,
		Short: "Environment variable reference",

		ValidArgsFunction: cobra.NoFileCompletions,
	}
	cmd.SetHelpFunc(helpFunc)
	return cmd
}

func helpFunc(cmd *cobra.Command, _ []string) {
	format, _ := output.FromContext(cmd.Context())

	var result strings.Builder
	if format == output.FormatMarkdown {
		result.WriteString("# Environment Variables\n\n")
	} else {
		result.WriteString("Environment Variables\n\n")
	}

	t := output.NewTable()

	bold := color.New(color.Bold).Sprint
	if format == output.FormatMarkdown {
		t.AppendHeader(table.Row{"Name", "Usage", "Default"})
	} else {
		t.AppendHeader(table.Row{bold("Name"), bold("Usage"), bold("Default")})
	}

	root := cmd.Root()
	excludeNames := []string{"completion", "help", "version"}
	var rows []table.Row
	italic := color.New(color.Italic).Sprint
	root.Flags().VisitAll(func(flag *pflag.Flag) {
		if slices.Contains(excludeNames, flag.Name) {
			return
		}

		var value string
		switch flag.Value.Type() {
		case "stringSlice":
			value = strings.Join(flag.Value.(pflag.SliceValue).GetSlice(), ",")
		default:
			value = flag.Value.String()
		}

		if format == output.FormatMarkdown {
			if value == "" {
				value = " "
			}
			rows = append(rows, table.Row{"`" + config.EnvName(flag.Name) + "`", flag.Usage, "`" + value + "`"})
		} else {
			rows = append(rows, table.Row{bold(config.EnvName(flag.Name)), flag.Usage, italic(value)})
		}
	})
	slices.SortFunc(rows, func(a, b table.Row) int {
		return cmp.Compare(a[0].(string), b[0].(string))
	})
	t.AppendRows(rows)

	if format == output.FormatMarkdown {
		result.WriteString(t.RenderMarkdown())
	} else {
		result.WriteString(t.Render())
	}
	result.WriteByte('\n')
	_, _ = io.WriteString(cmd.OutOrStdout(), result.String())
}
