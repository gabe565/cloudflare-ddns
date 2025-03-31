package sources

import (
	"io"
	"strings"

	"gabe565.com/cloudflare-ddns/internal/config"
	"gabe565.com/cloudflare-ddns/internal/output"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

const Name = "sources"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   Name,
		Short: "Source reference",

		ValidArgsFunction: cobra.NoFileCompletions,
	}
	cmd.SetHelpFunc(helpFunc)
	return cmd
}

func helpFunc(cmd *cobra.Command, _ []string) {
	format, _ := output.FromContext(cmd.Context())

	italic := color.New(color.Italic).Sprint
	var result strings.Builder
	if format == output.FormatMarkdown {
		result.WriteString("# Sources\n\nThe `--source` flag lets you define which sources are used to get your public IP address.\n\n" +
			"## Available Sources\n\n")
	} else {
		result.WriteString("The " + italic("--source") + " flag lets you define which sources are used to get your public IP address.\n\n" +
			"Available Sources:\n")
	}

	t := output.NewTable()

	bold := color.New(color.Bold).Sprint
	if format == output.FormatMarkdown {
		t.AppendHeader(table.Row{"Name", "Description"})
	} else {
		t.AppendHeader(table.Row{bold("Name"), bold("Description")})
	}

	sources := config.SourceValues()

	for _, v := range sources {
		if format == output.FormatMarkdown {
			t.AppendRow(table.Row{"`" + v.String() + "`", v.Description(output.FormatMarkdown)})
		} else {
			t.AppendRow(table.Row{bold(v.String()), v.Description(output.FormatANSI)})
		}
	}

	if format == output.FormatMarkdown {
		result.WriteString(t.RenderMarkdown())
	} else {
		result.WriteString(t.Render())
	}
	result.WriteByte('\n')
	_, _ = io.WriteString(cmd.OutOrStdout(), result.String())
}
