package envs

import (
	"cmp"
	"io"
	"slices"
	"strings"

	"gabe565.com/cloudflare-ddns/internal/config"
	"gabe565.com/cloudflare-ddns/internal/output"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
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

	t := table.New().
		Headers("Name", "Usage", "Default")

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
		italic := pad.Italic(true)
		t.StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case col == 0, row == -1:
				return bold
			case col == 2:
				return italic
			default:
				return pad
			}
		})
	}

	root := cmd.Root()
	excludeNames := []string{"completion", "help", "version"}
	var rows [][]string
	root.Flags().VisitAll(func(flag *pflag.Flag) {
		if slices.Contains(excludeNames, flag.Name) {
			return
		}

		var value string
		switch fv := flag.Value.(type) {
		case pflag.SliceValue:
			value = strings.Join(fv.GetSlice(), ",")
		default:
			value = flag.Value.String()
		}

		if format == output.FormatMarkdown {
			if value == "" {
				value = " "
			}
			rows = append(rows, []string{"`" + config.EnvName(flag.Name) + "`", flag.Usage, "`" + value + "`"})
		} else {
			rows = append(rows, []string{config.EnvName(flag.Name), flag.Usage, value})
		}
	})
	slices.SortFunc(rows, func(a, b []string) int {
		return cmp.Compare(a[0], b[0])
	})
	t.Rows(rows...)

	result.WriteString(t.Render())
	result.WriteByte('\n')
	_, _ = io.WriteString(cmd.OutOrStdout(), result.String())
}
