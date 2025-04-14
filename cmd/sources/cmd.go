package sources

import (
	"io"
	"strings"

	"gabe565.com/cloudflare-ddns/internal/lookup"
	"gabe565.com/cloudflare-ddns/internal/output"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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

	italic := lipgloss.NewStyle().Italic(true).Render
	var result strings.Builder
	if format == output.FormatMarkdown {
		result.WriteString(
			"# Sources\n\nThe `--source` flag lets you define which sources are used to get your public IP address.\n\n" +
				"## Available Sources\n\n",
		)
	} else {
		result.WriteString("The " + italic("--source") + " flag lets you define which sources are used to get your public IP address.\n\n" +
			"Available Sources:\n")
	}

	t := table.New().
		Headers("Name", "Description")

	pad := lipgloss.NewStyle().Padding(0, 1)
	if format == output.FormatMarkdown {
		t.Border(lipgloss.MarkdownBorder()).
			BorderTop(false).
			BorderBottom(false).
			StyleFunc(func(int, int) lipgloss.Style {
				return pad
			})
	} else {
		bold := pad.Bold(true)
		t.StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case col == 0, row == -1:
				return bold
			default:
				return pad
			}
		})
	}

	sources := lookup.SourceValues()

	for _, v := range sources {
		if format == output.FormatMarkdown {
			t.Row("`"+v.String()+"`", v.Description(output.FormatMarkdown))
		} else {
			t.Row(v.String(), v.Description(output.FormatANSI))
		}
	}

	result.WriteString(t.Render())
	result.WriteByte('\n')
	_, _ = io.WriteString(cmd.OutOrStdout(), result.String())
}
